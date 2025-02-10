package frame

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go.uber.org/fx"

	glog "github.com/zunkk/go-project-startup/pkg/log"
	"github.com/zunkk/go-project-startup/pkg/reqctx"
)

var log = glog.WithModule("sidecar")

func init() {
	RegisterComponents(NewSidecar)
}

type Component interface {
	ComponentName() string

	// Start must be non-blocking, use ComponentShutdown to stop on goroutine
	Start() error

	// Stop must be non-blocking
	Stop() error
}

type BuildConfig struct {
	Ctx       context.Context
	Wg        *sync.WaitGroup
	Version   string
	NodeIndex uint16
}

type Sidecar struct {
	// internal
	lc fx.Lifecycle

	sd                    fx.Shutdowner
	appReadyCallbacks     []func() error
	appReadyCallbackNames []string
	lock                  *sync.RWMutex
	wg                    *sync.WaitGroup
	version               string

	// common
	Ctx context.Context

	Logger        *slog.Logger
	UUIDGenerator *snowflake.Node
}

func NewSidecar(cfg *BuildConfig, lc fx.Lifecycle, sd fx.Shutdowner) (*Sidecar, error) {
	uuidGenerator, err := snowflake.NewNode(int64(cfg.NodeIndex))
	if err != nil {
		return nil, err
	}
	return &Sidecar{
		lc:            lc,
		sd:            sd,
		lock:          new(sync.RWMutex),
		wg:            cfg.Wg,
		version:       cfg.Version,
		Ctx:           cfg.Ctx,
		Logger:        log,
		UUIDGenerator: uuidGenerator,
	}, nil
}

func (c *Sidecar) RegisterLifecycleHook(component Component) {
	c.lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			start := time.Now()
			if err := component.Start(); err != nil {
				return errors.Wrapf(err, "componen[%s] start failed", component.ComponentName())
			}
			log.Info(fmt.Sprintf("Component[%s] started", component.ComponentName()), "time_cost", time.Since(start))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			start := time.Now()
			if err := component.Stop(); err != nil {
				return errors.Wrapf(err, "componen[%s] stop failed", component.ComponentName())
			}
			log.Info(fmt.Sprintf("Component[%s] stopped", component.ComponentName()), "time_cost", time.Since(start))
			return nil
		},
	})
}

func (c *Sidecar) RegisterAppReadyCallback(name string, callback func() error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.appReadyCallbackNames = append(c.appReadyCallbackNames, name)
	c.appReadyCallbacks = append(c.appReadyCallbacks, callback)
}

func (c *Sidecar) ExecuteAppReadyCallbacks() {
	lo.ForEach(c.appReadyCallbacks, func(callback func() error, i int) {
		callbackFn := callback
		idx := i
		c.SafeGo(func() {
			err := callbackFn()
			if err != nil {
				log.Warn("Failed to execute app ready callback", "err", err, "name", c.appReadyCallbackNames[idx])
				return
			}
			log.Info("Executed app ready callback", "name", c.appReadyCallbackNames[idx])
		})
	})
}

func (c *Sidecar) ComponentShutdown() {
	if err := c.sd.Shutdown(); err != nil {
		log.Error("App shutdown error", "err", err)
	}
}

func (c *Sidecar) SafeGo(fn func()) {
	go func() {
		defer func() {
			c.Recovery()
		}()
		fn()
	}()
}

func (c *Sidecar) SafeGoPersistentTask(fn func()) {
	c.wg.Add(1)
	go func() {
		defer func() {
			c.Recovery()
			c.wg.Done()
		}()
		fn()
	}()
}

func (c *Sidecar) Recovery() {
	if c.IsDevVersion() {
		return
	}

	if err := recover(); err != nil {
		log.Error(fmt.Sprintf("panic: %v", err))
		for i := 0; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			log.Error(fmt.Sprintf("%v %v", file, line))
		}
	}
}

func (c *Sidecar) RecoverExecute(executor func() error) (pErr error) {
	if !c.IsDevVersion() {
		defer func() {
			if r := recover(); r != nil {
				pErr = fmt.Errorf("%v:\n%s", r, panicTrace())
			}
		}()
	}

	return executor()
}

func (c *Sidecar) BackgroundContext() *reqctx.ReqCtx {
	reqID := c.UUIDGenerator.Generate()
	return reqctx.NewReqCtx(c.Ctx, log, int64(reqID), "system")
}

func (c *Sidecar) IsDevVersion() bool {
	return c.version == "dev"
}

func (c *Sidecar) IsTestVersion() bool {
	return c.version == "test"
}

func (c *Sidecar) IsProdVersion() bool {
	return c.version == "prod"
}

type ScheduledTaskHandler struct {
	taskName       string
	cancelCtx      context.Context
	cancelFunc     context.CancelFunc
	canceled       bool
	paused         bool
	waitCanceledCh chan struct{}
}

func (c *ScheduledTaskHandler) IsRunning() bool {
	return !c.canceled
}

func (c *ScheduledTaskHandler) IsPaused() bool {
	return c.paused
}

func (c *ScheduledTaskHandler) Pause() {
	c.paused = true
}

func (c *ScheduledTaskHandler) Resume() {
	c.paused = false
}

func (c *ScheduledTaskHandler) Cancel() {
	c.cancelFunc()
	if c.canceled {
		return
	}
	select {
	case <-c.waitCanceledCh:
	case <-time.After(10 * time.Second):
		log.Warn("Wait scheduled task canceled timeout", "task", c.taskName)
	}
}

func (c *Sidecar) runScheduledTask(taskName string, isPersistent bool, interval time.Duration, cancelCtx context.Context, cancelFunc context.CancelFunc, taskExecutorOnTick func(ctx context.Context) (cancel bool, err error)) *ScheduledTaskHandler {
	log.Info("Scheduled task started", "task", taskName, "interval", interval)
	handler := &ScheduledTaskHandler{
		taskName:       taskName,
		cancelCtx:      cancelCtx,
		cancelFunc:     cancelFunc,
		canceled:       false,
		paused:         false,
		waitCanceledCh: make(chan struct{}, 1),
	}

	runner := func() {
		tk := time.NewTicker(interval)
		defer tk.Stop()
		for {
			select {
			case <-tk.C:
				if err := RecoverExecute(func() error {
					if handler.paused {
						return nil
					}
					cancel, err := taskExecutorOnTick(cancelCtx)
					if cancel {
						handler.canceled = true
					}
					return err
				}); err != nil {
					if strings.Contains(err.Error(), "context canceled") {
						handler.canceled = true
					} else {
						log.Warn("Do scheduled task executor error", "task", taskName, "err", err)
					}
				}
			case <-cancelCtx.Done():
				handler.canceled = true
			}
			if handler.canceled {
				break
			}
		}
		log.Info("Scheduled task stopped", "task", taskName, "interval", interval)
		handler.cancelFunc()
		handler.waitCanceledCh <- struct{}{}
	}
	if isPersistent {
		c.SafeGoPersistentTask(runner)
	} else {
		c.SafeGo(runner)
	}
	return handler
}

// RunScheduledTask will poll the task executor when time tick reached
func (c *Sidecar) RunScheduledTask(taskName string, isPersistent bool, interval time.Duration, taskExecutorOnTick func(ctx context.Context) (err error)) *ScheduledTaskHandler {
	return c.RunScheduledTaskWithCtx(c.Ctx, taskName, isPersistent, interval, taskExecutorOnTick)
}

func (c *Sidecar) RunScheduledTaskWithCtx(ctx context.Context, taskName string, isPersistent bool, interval time.Duration, taskExecutorOnTick func(ctx context.Context) (err error)) *ScheduledTaskHandler {
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	return c.runScheduledTask(taskName, isPersistent, interval, cancelCtx, cancelFunc, func(ctx context.Context) (cancel bool, err error) {
		return false, taskExecutorOnTick(ctx)
	})
}

func (c *Sidecar) RunScheduledTaskWithCancel(taskName string, isPersistent bool, interval time.Duration, taskExecutorOnTick func(ctx context.Context) (cancel bool, err error)) *ScheduledTaskHandler {
	cancelCtx, cancelFunc := context.WithCancel(c.Ctx)
	return c.runScheduledTask(taskName, isPersistent, interval, cancelCtx, cancelFunc, taskExecutorOnTick)
}

func (c *Sidecar) RunScheduledTaskWithPrepare(taskName string, isPersistent bool, interval time.Duration, prepare func(ctx context.Context) (err error), taskExecutorOnTick func(ctx context.Context) (err error)) (*ScheduledTaskHandler, error) {
	cancelCtx, cancelFunc := context.WithCancel(c.Ctx)
	if err := prepare(cancelCtx); err != nil {
		cancelFunc()
		return nil, errors.Wrapf(err, "prepare core task %s failed", taskName)
	}
	return c.runScheduledTask(taskName, isPersistent, interval, cancelCtx, cancelFunc, func(ctx context.Context) (cancel bool, err error) {
		return false, taskExecutorOnTick(ctx)
	}), nil
}

type CoreTaskHandler struct {
	taskName       string
	cancelCtx      context.Context
	cancelFunc     context.CancelFunc
	canceled       bool
	waitCanceledCh chan struct{}
	cleaner        func()
	cleanerDoOnce  *sync.Once
}

func (c *CoreTaskHandler) IsRunning() bool {
	return !c.canceled
}

func (c *CoreTaskHandler) Cancel() {
	c.cancelFunc()
	if c.canceled {
		return
	}
	c.cleanerDoOnce.Do(func() {
		if c.cleaner != nil {
			c.cleaner()
		}
	})
	select {
	case <-c.waitCanceledCh:
	case <-time.After(10 * time.Second):
		log.Warn("Wait core task canceled timeout", "task", c.taskName)
	}
}

func (c *Sidecar) runCoreTask(taskName string, isPersistent bool, cancelCtx context.Context, cancelFunc context.CancelFunc, taskExecutor func(ctx context.Context) (cancel bool, err error), cleaner func()) *CoreTaskHandler {
	log.Info("Core task started", "task", taskName)
	handler := &CoreTaskHandler{
		taskName:       taskName,
		cancelCtx:      cancelCtx,
		cancelFunc:     cancelFunc,
		canceled:       false,
		waitCanceledCh: make(chan struct{}, 1),
		cleaner:        cleaner,
		cleanerDoOnce:  new(sync.Once),
	}
	runner := func() {
		for {
			select {
			case <-cancelCtx.Done():
				handler.canceled = true
			default:
			}
			if err := RecoverExecute(func() error {
				cancel, err := taskExecutor(cancelCtx)
				if cancel {
					handler.canceled = true
				}
				return err
			}); err != nil {
				if strings.Contains(err.Error(), "context canceled") {
					handler.canceled = true
				} else {
					log.Warn("Do core task executor error", "task", taskName, "err", err)
				}
			}
			if handler.canceled {
				break
			}
		}
		handler.cleanerDoOnce.Do(func() {
			if cleaner != nil {
				cleaner()
			}
		})
		log.Info("Core task stopped", "task", taskName)
		handler.cancelFunc()
		handler.waitCanceledCh <- struct{}{}
	}
	if isPersistent {
		c.SafeGoPersistentTask(runner)
	} else {
		c.SafeGo(runner)
	}
	return handler
}

// RunCoreTask will always poll the task executor
func (c *Sidecar) RunCoreTask(taskName string, isPersistent bool, taskExecutor func(ctx context.Context) (cancel bool, err error)) *CoreTaskHandler {
	cancelCtx, cancelFunc := context.WithCancel(c.Ctx)
	return c.runCoreTask(taskName, isPersistent, cancelCtx, cancelFunc, taskExecutor, nil)
}

func (c *Sidecar) RunCoreTaskWithPrepare(taskName string, isPersistent bool, prepare func(ctx context.Context) (cleaner func(), err error), taskExecutor func(ctx context.Context) (cancel bool, err error)) (*CoreTaskHandler, error) {
	cancelCtx, cancelFunc := context.WithCancel(c.Ctx)
	cleaner, err := prepare(cancelCtx)
	if err != nil {
		cancelFunc()
		return nil, errors.Wrapf(err, "prepare core task %s failed", taskName)
	}

	return c.runCoreTask(taskName, isPersistent, cancelCtx, cancelFunc, taskExecutor, cleaner), nil
}

func panicTrace() string {
	return string(debug.Stack())
}

func RecoverExecute(executor func() error) (pErr error) {
	defer func() {
		if r := recover(); r != nil {
			pErr = fmt.Errorf("%v:\n%s", r, panicTrace())
		}
	}()
	return executor()
}

package frame

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"runtime/debug"
	"sync"

	"github.com/bwmarrin/snowflake"
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

	sd                fx.Shutdowner
	appReadyCallbacks []func() error
	lock              *sync.RWMutex
	wg                *sync.WaitGroup
	version           string

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
			return component.Start()
		},
		OnStop: func(ctx context.Context) error {
			return component.Stop()
		},
	})
}

func (c *Sidecar) RegisterAppReadyCallback(callback func() error) {
	c.lock.Lock()
	c.appReadyCallbacks = append(c.appReadyCallbacks, callback)
	c.lock.Unlock()
}

func (c *Sidecar) ExecuteAppReadyCallbacks() {
	lo.ForEach(c.appReadyCallbacks, func(callback func() error, _ int) {
		callbackFn := callback
		c.SafeGo(func() {
			err := callbackFn()
			if err != nil {
				log.Warn("Failed to execute app ready callback", "err", err)
			}
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

func panicTrace() string {
	return string(debug.Stack())
}

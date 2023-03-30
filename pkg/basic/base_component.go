package basic

import (
	"context"
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"

	"github.com/zunkk/go-project-startup/pkg/reqctx"
)

func init() {
	RegisterComponents(NewBaseComponent)
}

type Component interface {
	// Start must be non-blocking, use ComponentShutdown to stop on goroutine
	Start() error
	// Stop must be non-blocking
	Stop() error
}

type BuildConfig struct {
	Ctx       context.Context
	Logger    *logrus.Logger
	Wg        *sync.WaitGroup
	Version   string
	NodeIndex uint16
}

type BaseComponent struct {
	// internal
	lc                fx.Lifecycle
	sd                fx.Shutdowner
	appReadyCallbacks []func() error
	lock              *sync.RWMutex
	wg                *sync.WaitGroup
	version           string

	// common
	Ctx           context.Context
	Logger        *logrus.Logger
	UUIDGenerator *snowflake.Node
}

func NewBaseComponent(cfg *BuildConfig, lc fx.Lifecycle, sd fx.Shutdowner) (*BaseComponent, error) {
	uuidGenerator, err := snowflake.NewNode(int64(cfg.NodeIndex))
	if err != nil {
		return nil, err
	}
	return &BaseComponent{
		lc:            lc,
		sd:            sd,
		lock:          new(sync.RWMutex),
		wg:            cfg.Wg,
		version:       cfg.Version,
		Ctx:           cfg.Ctx,
		Logger:        cfg.Logger,
		UUIDGenerator: uuidGenerator,
	}, nil
}

func (c *BaseComponent) RegisterLifecycleHook(component Component) {
	c.lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return component.Start()
		},
		OnStop: func(ctx context.Context) error {
			return component.Stop()
		},
	})
}

func (c *BaseComponent) RegisterAppReadyCallback(callback func() error) {
	c.lock.Lock()
	c.appReadyCallbacks = append(c.appReadyCallbacks, callback)
	c.lock.Unlock()
}

func (c *BaseComponent) ExecuteAppReadyCallbacks() {
	for _, callback := range c.appReadyCallbacks {
		callback := callback
		c.SafeGo(func() {
			err := callback()
			if err != nil {
				c.Logger.WithField("err", err).Warn("Failed to execute app ready callback")
			}
		})
	}
}

func (c *BaseComponent) ComponentShutdown() {
	if err := c.sd.Shutdown(); err != nil {
		c.Logger.Errorf("App shutdown error: %v", err)
	}
}

func (c *BaseComponent) SafeGo(fn func()) {
	go func() {
		defer func() {
			c.Recovery()
		}()
		fn()
	}()
}

func (c *BaseComponent) SafeGoPersistentTask(fn func()) {
	c.wg.Add(1)
	go func() {
		defer func() {
			c.Recovery()
			c.wg.Done()
		}()
		fn()
	}()
}

func (c *BaseComponent) Recovery() {
	if c.IsDevVersion() {
		return
	}

	if err := recover(); err != nil {
		c.Logger.Errorf("panic: %v", err)
		for i := 0; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			c.Logger.Errorf("%v %v", file, line)
		}
	}
}

func (c *BaseComponent) RecoverExecute(executor func() error) (pErr error) {
	if !c.IsDevVersion() {
		defer func() {
			if r := recover(); r != nil {
				pErr = fmt.Errorf("%v:\n%s", r, panicTrace())
			}
		}()
	}

	return executor()
}

func (c *BaseComponent) BackgroundContext() *reqctx.ReqCtx {
	reqID := c.UUIDGenerator.Generate()
	return reqctx.NewReqCtx(c.Ctx, c.Logger, int64(reqID), "system")
}

func (c *BaseComponent) IsDevVersion() bool {
	return c.version == "dev"
}

func panicTrace() string {
	return string(debug.Stack())
}

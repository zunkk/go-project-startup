package basic

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

const (
	lifecycleTimeout = 20 * time.Second
)

type App interface {
	Run() (exitCode int)
	Start(ctx context.Context) error
}

type appInternal struct {
	ctx    context.Context
	cancel context.CancelFunc
	logger *logrus.Logger
	wg     *sync.WaitGroup
	*fx.App
}

func (app *appInternal) Run() (exitCode int) {
	defer func() {
		app.cancel()
		ch := make(chan bool)
		go func() {
			app.wg.Wait()
			ch <- false
		}()

		select {
		case <-time.After(lifecycleTimeout):
			app.logger.Warn("Wait for components to terminate timeout")
		case <-ch:
		}
	}()

	startCtx, startCancel := context.WithTimeout(app.ctx, lifecycleTimeout)
	defer startCancel()
	if err := app.Start(startCtx); err != nil {
		app.logger.WithFields(logrus.Fields{"err": err}).Error("Start components failed")
		return 1
	}

	sig := <-app.Done()
	app.logger.Infof("Receive exit signal: %v", sig.String())

	stopCtx, stopCancel := context.WithTimeout(app.ctx, lifecycleTimeout)
	defer stopCancel()
	if err := app.Stop(stopCtx); err != nil {
		app.logger.WithFields(logrus.Fields{"err": err}).Error("Stop components failed")
		return 1
	}

	return 0
}

var lock = new(sync.Mutex)
var constructors []interface{}

func RegisterComponents(componentConstructors ...interface{}) {
	lock.Lock()
	constructors = append(constructors, componentConstructors...)
	lock.Unlock()
}

func BuildApp(bgCtx context.Context, logger *logrus.Logger, nodeIndex uint16, version string, supports []interface{}, fxInvokeFunc interface{}, targetPopulate ...interface{}) (App, error) {
	ctx, cancel := context.WithCancel(bgCtx)
	wg := new(sync.WaitGroup)
	supports = append(supports, &BuildConfig{
		Ctx:       ctx,
		Logger:    logger,
		Wg:        wg,
		Version:   version,
		NodeIndex: nodeIndex,
	})
	app := fx.New(
		fx.NopLogger,
		fx.Supply(supports...),
		fx.Provide(
			constructors...,
		),
		fx.Populate(targetPopulate...),
		fx.Invoke(fxInvokeFunc),
	)
	if app.Err() != nil {
		cancel()
		return nil, errors.Wrap(app.Err(), "app setup error")
	}
	return &appInternal{
		ctx:    ctx,
		cancel: cancel,
		logger: logger,
		wg:     wg,
		App:    app,
	}, nil
}

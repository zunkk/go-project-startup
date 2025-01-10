package frame

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/pkg/errors"
	"go.uber.org/fx"

	glog "github.com/zunkk/go-project-startup/pkg/log"
	"github.com/zunkk/go-project-startup/pkg/repo"
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
			log.Warn("Wait for components to terminate timeout")
		case <-ch:
		}
	}()

	startCtx, startCancel := context.WithTimeout(app.ctx, lifecycleTimeout)
	defer startCancel()
	if err := app.Start(startCtx); err != nil {
		log.Error("Start components failed", "err", err)
		return 1
	}

	log.Info(fmt.Sprintf("%s is ready", repo.AppName))
	fig := figure.NewFigure(repo.AppName, "slant", true)
	figWeight := 0
	for _, printRow := range fig.Slicify() {
		if len(printRow) > figWeight {
			figWeight = len(printRow)
		}
	}
	decorateLine := strings.Repeat("=", figWeight)
	log.Info(fmt.Sprintf("%s\n%s\n%s\n", decorateLine, fig.String(), decorateLine), glog.OnlyWriteMsgWithoutFormatterField, nil)

	sig := <-app.Done()
	log.Info(fmt.Sprintf("Receive exit signal: %v", sig.String()))

	stopCtx, stopCancel := context.WithTimeout(app.ctx, lifecycleTimeout)
	defer stopCancel()
	if err := app.Stop(stopCtx); err != nil {
		log.Error("Stop components failed", "err", err)
		return 1
	}

	return 0
}

var lock = new(sync.Mutex)
var constructors []any

func RegisterComponents(componentConstructors ...any) {
	lock.Lock()
	constructors = append(constructors, componentConstructors...)
	lock.Unlock()
}

func BuildApp(bgCtx context.Context, uuidNodeIndex uint16, version string, supports []any, fxInvokeFunc any, targetPopulate ...any) (App, error) {
	ctx, cancel := context.WithCancel(bgCtx)
	wg := new(sync.WaitGroup)
	supports = append(supports, &BuildConfig{
		Ctx:       ctx,
		Wg:        wg,
		Version:   version,
		NodeIndex: uuidNodeIndex,
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
		wg:     wg,
		App:    app,
	}, nil
}

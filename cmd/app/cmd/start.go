package cmd

import (
	"os"

	"github.com/zunkk/go-project-startup/api/rest"
	"github.com/zunkk/go-project-startup/internal/pkg/base"
	internalconfig "github.com/zunkk/go-project-startup/internal/pkg/config"
	"github.com/zunkk/go-project-startup/pkg/basic"
	"github.com/zunkk/go-project-startup/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type APP struct {
	baseComponent *base.Component
	stopFuncs     []func()
}

func NewApp(baseComponent *base.Component, server *rest.Server) *APP {
	app := &APP{
		baseComponent: baseComponent,
	}
	app.baseComponent.RegisterLifecycleHook(app)
	return app
}

// execute when all components started
func (app *APP) Start() error {
	if err := config.WritePid(app.baseComponent.Config.RootPath); err != nil {
		return err
	}
	app.baseComponent.Logger.Infof("%s is ready", config.AppName)
	app.baseComponent.ExecuteAppReadyCallbacks()
	return nil
}

func (app *APP) stopDebugService() {
	for _, stopFunc := range app.stopFuncs {
		stopFunc()
	}
}

func (app *APP) Stop() error {
	if err := config.RemovePID(app.baseComponent.Config.RootPath); err != nil {
		app.baseComponent.Logger.WithFields(logrus.Fields{"err": err}).Warn("Failed to remove pid file")
	}
	app.stopDebugService()
	return nil
}

func Start(ctx *cli.Context) error {
	cfg, err := config.Load(internalconfig.DefaultConfig)
	if err != nil {
		return err
	}
	logger, err := config.InitLogger(ctx.Context, cfg.RootPath, cfg.Log)
	if err != nil {
		return err
	}
	config.PrintSystemInfo(cfg.RootPath, logger.Infof)
	exe, err := os.Executable()
	if err == nil {
		logger.Infof("Binary path: %s", exe)
	}
	logger.Infof("PID: %d", os.Getpid())
	logger.Infof("Node index: %d", cfg.App.NodeIndex)

	basic.RegisterComponents(NewApp)
	app, err := basic.BuildApp(ctx.Context, logger, cfg.App.NodeIndex, config.Version, []interface{}{cfg}, func(app *APP) {})
	if err != nil {
		logger.WithField("err", err).Error("Build app failed")
		return nil
	}

	if exitCode := app.Run(); exitCode != 0 {
		logger.Infof("%s is stopped", config.AppName)
		os.Exit(exitCode)
		return nil
	}
	logger.Infof("%s is stopped", config.AppName)
	return nil
}

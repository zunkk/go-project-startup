package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/urfave/cli/v2"

	"github.com/zunkk/go-project-startup/api/rest"
	"github.com/zunkk/go-project-startup/internal/pkg/base"
	internalconfig "github.com/zunkk/go-project-startup/internal/pkg/config"
	"github.com/zunkk/go-project-startup/pkg/config"
	"github.com/zunkk/go-project-startup/pkg/frame"
	glog "github.com/zunkk/go-project-startup/pkg/log"
)

var log = glog.WithModule("app")

type APP struct {
	sidecar   *base.CustomSidecar
	stopFuncs []func()
}

func NewApp(sidecar *base.CustomSidecar, server *rest.Server) *APP {
	app := &APP{
		sidecar: sidecar,
	}
	app.sidecar.RegisterLifecycleHook(app)
	return app
}

// execute when all components started
func (app *APP) Start() error {
	log.Info(fmt.Sprintf("%s is ready", config.AppName))
	fig := figure.NewFigure(config.AppName, "slant", true)
	figWeight := 0
	for _, printRow := range fig.Slicify() {
		if len(printRow) > figWeight {
			figWeight = len(printRow)
		}
	}
	decorateLine := strings.Repeat("=", figWeight)
	log.Info(fmt.Sprintf("%s\n%s\n%s\n", decorateLine, fig.String(), decorateLine), glog.OnlyWriteMsgWithoutFormatterField, nil)
	app.sidecar.ExecuteAppReadyCallbacks()
	return nil
}

func (app *APP) stopDebugService() {
	for _, stopFunc := range app.stopFuncs {
		stopFunc()
	}
}

func (app *APP) Stop() error {
	app.stopDebugService()
	return nil
}

func Start(ctx *cli.Context) error {
	cfg, err := config.Load(internalconfig.DefaultConfig)
	if err != nil {
		return err
	}
	if err := config.InitLogger(ctx.Context, cfg.RepoPath, cfg.Log); err != nil {
		return err
	}
	config.PrintSystemInfo(cfg.RepoPath, func(format string, args ...any) {
		log.Info(fmt.Sprintf(format, args...))
	})
	exe, err := os.Executable()
	if err == nil {
		log.Info(fmt.Sprintf("Executable: %s", exe))
	}
	log.Info(fmt.Sprintf("PID: %d", os.Getpid()))
	log.Info(fmt.Sprintf("UUID node index: %d", cfg.App.UUIDNodeIndex))

	frame.RegisterComponents(NewApp)
	app, err := frame.BuildApp(ctx.Context, cfg.App.UUIDNodeIndex, config.Version, []any{cfg}, func(app *APP) {})
	if err != nil {
		log.Error("Build app failed", "err", err)
		return nil
	}

	if exitCode := app.Run(); exitCode != 0 {
		log.Info(fmt.Sprintf("%s is stopped", config.AppName))
		os.Exit(exitCode)
		return nil
	}
	log.Info(fmt.Sprintf("%s is stopped", config.AppName))
	return nil
}

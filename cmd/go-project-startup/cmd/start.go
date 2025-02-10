package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/zunkk/go-project-startup/api/rest"
	"github.com/zunkk/go-project-startup/internal/pkg/base"
	internalconfig "github.com/zunkk/go-project-startup/internal/pkg/config"
	"github.com/zunkk/go-sidecar/frame"
	glog "github.com/zunkk/go-sidecar/log"
	"github.com/zunkk/go-sidecar/repo"
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

func (app *APP) ComponentName() string {
	return "app"
}

// execute when all components started
func (app *APP) Start() error {
	if err := repo.WritePid(app.sidecar.Repo.RepoPath); err != nil {
		return errors.Wrap(err, "write pid failed")
	}
	app.sidecar.ExecuteAppReadyCallbacks()
	return nil
}

func (app *APP) Stop() error {
	for _, stopFunc := range app.stopFuncs {
		stopFunc()
	}
	if err := repo.RemovePID(app.sidecar.Repo.RepoPath); err != nil {
		return errors.Wrap(err, "remove pid failed")
	}
	return nil
}

func Start(ctx *cli.Context) error {
	rep, err := repo.Load(repo.RootPath, internalconfig.DefaultConfig)
	if err != nil {
		return err
	}
	if err := repo.InitLogger(ctx.Context, rep.RepoPath, rep.Cfg.Log); err != nil {
		return err
	}
	repo.PrintSystemInfo(rep.RepoPath, func(format string, args ...any) {
		log.Info(fmt.Sprintf(format, args...))
	})
	exe, err := os.Executable()
	if err == nil {
		log.Info(fmt.Sprintf("Executable: %s", exe))
	}
	log.Info(fmt.Sprintf("PID: %d", os.Getpid()))
	log.Info(fmt.Sprintf("UUID node index: %d", rep.Cfg.App.UUIDNodeIndex))

	frame.RegisterComponents(NewApp)
	app, err := frame.BuildApp(repo.AppName, ctx.Context, rep.Cfg.App.UUIDNodeIndex, repo.Version, []any{rep}, func(app *APP) {})
	if err != nil {
		log.Error("Build app failed", "err", err)
		return nil
	}

	if exitCode := app.Run(); exitCode != 0 {
		log.Info(fmt.Sprintf("%s is stopped", repo.AppName))
		os.Exit(exitCode)
		return nil
	}
	log.Info(fmt.Sprintf("%s is stopped", repo.AppName))
	return nil
}

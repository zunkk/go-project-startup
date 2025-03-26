package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/zunkk/go-project-startup/cmd/go-project-startup/cmd"
	clicmd "github.com/zunkk/go-project-startup/cmd/go-project-startup/cmd/cli"
	configcmd "github.com/zunkk/go-project-startup/cmd/go-project-startup/cmd/config"
	"github.com/zunkk/go-project-startup/internal/pkg/config"
	"github.com/zunkk/go-sidecar/repo"
)

func main() {
	repo.InitGlobalInfo(config.AppName, config.AppDesc, config.Version, config.BuildTime, config.CommitID)

	app := cli.NewApp()
	app.Name = repo.AppName
	app.Usage = repo.AppDesc
	app.HideVersion = true
	app.Description = "Run COMMAND --help for more information on a command"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "repo_path",
			Aliases:     []string{"rp"},
			Destination: &repo.RootPath,
			Usage:       "repo path",
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "env_file_path",
			Aliases:     []string{"efp"},
			Destination: &repo.EnvFilePath,
			Usage:       "env file path",
			Required:    false,
		},
	}
	app.Before = func(c *cli.Context) error {
		repo.LoadEnvFile()
		return nil
	}

	app.Commands = []*cli.Command{
		{
			Name:   "start",
			Usage:  "Start app",
			Action: cmd.Start,
		},
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "Show version",
			Action: func(c *cli.Context) error {
				repo.PrintSystemInfo("", func(format string, args ...any) {
					fmt.Printf(format+"\n", args...)
				})
				return nil
			},
		},
		configcmd.Command,
		clicmd.Command,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("app run error: %v\n", err)
	}
}

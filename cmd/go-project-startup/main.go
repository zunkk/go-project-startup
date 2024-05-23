package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/zunkk/go-project-startup/cmd/go-project-startup/cmd"
	configcmd "github.com/zunkk/go-project-startup/cmd/go-project-startup/cmd/config"
	"github.com/zunkk/go-project-startup/pkg/repo"
)

func main() {
	repo.LoadEnvFile()

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
			Destination: &repo.EnvPrefix,
			Usage:       "env file path",
			Required:    false,
		},
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
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("app run error: %v\n", err)
	}
}

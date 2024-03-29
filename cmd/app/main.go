package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/zunkk/go-project-startup/cmd/app/cmd"
	configcmd "github.com/zunkk/go-project-startup/cmd/app/cmd/config"
	"github.com/zunkk/go-project-startup/pkg/config"
)

func main() {
	app := cli.NewApp()
	app.Name = config.AppName
	app.Usage = config.AppDesc
	app.HideVersion = true
	app.Description = "Run COMMAND --help for more information on a command"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "repo_path",
			Aliases:     []string{"rp"},
			Destination: &config.RepoPath,
			Usage:       "root path",
			Required:    true,
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
				config.PrintSystemInfo("", func(format string, args ...any) {
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

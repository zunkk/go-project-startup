package config

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	internalconfig "github.com/zunkk/go-project-startup/internal/pkg/config"
	"github.com/zunkk/go-sidecar/repo"
)

// Command is The explorer commands
var Command = &cli.Command{
	Name:        "config",
	Usage:       "The config manage commands",
	Subcommands: subCommands,
}

var subCommands = []*cli.Command{
	{
		Name:   "check",
		Usage:  "Check if the config file is valid",
		Action: check,
	},
	{
		Name:   "show",
		Usage:  "Show the complete config processed by the environment variable",
		Action: show,
	},
	{
		Name:   "generate-default",
		Usage:  "Generate default config",
		Action: generateDefault,
	},
}

func check(ctx *cli.Context) error {
	cfg := internalconfig.DefaultConfig()
	if err := repo.ReadConfig(repo.RootPath, cfg); err != nil {
		fmt.Println("config file format error, please check:", err)
		os.Exit(1)
		return nil
	}
	return nil
}

func show(ctx *cli.Context) error {
	cfg := internalconfig.DefaultConfig()
	if err := repo.ReadConfig(repo.RootPath, cfg); err != nil {
		fmt.Println("read config file failed:", err)
		os.Exit(1)
		return nil
	}

	str, err := repo.MarshalConfig(cfg)
	if err != nil {
		fmt.Println("marshal config failed:", err)
		os.Exit(1)
		return nil
	}
	fmt.Println(str)
	return nil
}

func generateDefault(ctx *cli.Context) error {
	if repo.ExistConfigFile(repo.RootPath) {
		fmt.Println("config file already exists")
		os.Exit(1)
		return nil
	}

	cfg := internalconfig.DefaultConfig()
	if err := repo.WriteConfig(repo.RootPath, cfg); err != nil {
		fmt.Println("write config to file failed:", err)
		os.Exit(1)
		return nil
	}
	return nil
}

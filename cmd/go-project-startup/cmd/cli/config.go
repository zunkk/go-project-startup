package cli

import (
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/urfave/cli/v2"

	"github.com/zunkk/go-project-startup/internal/pkg/config"
)

var configCommand = &cli.Command{
	Name:  "config",
	Usage: "The config manage commands",
	Subcommands: []*cli.Command{
		{
			Name:   "show",
			Usage:  "Show config",
			Action: configShow,
		},
	},
}

func configShow(ctx *cli.Context) error {
	res, err := doRequest[config.Config](http.MethodGet, "/config/info", func(req *resty.Request) {})
	if err != nil {
		return err
	}
	return PrettyPrint(res)
}

package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/zunkk/go-project-startup/api/rest"
	"github.com/zunkk/go-sidecar/repo"
	"github.com/zunkk/go-sidecar/util"
)

const (
	baseIPCURL = "http://ipc.sock/api/v1"
)

var pingContent string

var Command = &cli.Command{
	Name:  "cli",
	Usage: "The ipc client tool",
	Subcommands: []*cli.Command{
		{
			Name:   "ping",
			Usage:  "Ping local server by ipc",
			Action: ping,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "content",
					Usage:       "Ping content",
					DefaultText: "now time",
					Required:    false,
					Destination: &pingContent,
				},
			},
		},
		configCommand,
	},
}

type emptyRes struct {
}

func doRequest[Res any](method string, url string, requestBuilder func(req *resty.Request)) (res *Res, err error) {
	ipcFile := filepath.Join(repo.RootPath, "ipc.sock")
	if !util.FileExist(ipcFile) {
		return nil, errors.New("bot is not running")
	}

	client := resty.New().SetTransport(&http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial("unix", ipcFile)
		},
	}).SetBaseURL(baseIPCURL)

	req := client.R().SetHeader("Content-Type", "application/json")
	if requestBuilder != nil {
		requestBuilder(req)
	}

	resp, err := req.Execute(method, url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, errors.New(resp.String())
	}
	return util.DecodeResponse[Res](resp.Body())
}

func PrettyPrint(d any) error {
	res, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(res))
	return nil
}

func ping(ctx *cli.Context) error {
	if pingContent == "" {
		pingContent = time.Now().Format(time.RFC3339)
	}
	pong, err := doRequest[rest.PingRes](http.MethodGet, "/ping", func(req *resty.Request) {
		req.SetQueryParam("ping", pingContent)
	})
	if err != nil {
		return err
	}
	fmt.Printf("send ping: %s, got pong: %s\n", pingContent, pong.Pong)
	return err
}

package profile

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/pkg/errors"
	glog "github.com/zunkk/go-project-startup/pkg/log"
)

var log = glog.WithModule("pprof")

type Pprof struct {
	ctx    context.Context
	enable bool
	port   uint64

	listener net.Listener
	server   *http.Server
}

func NewPprof(ctx context.Context, enable bool, port uint64) (*Pprof, error) {
	return &Pprof{
		ctx:    ctx,
		enable: enable,
		port:   port,
	}, nil
}

func (p *Pprof) Start() error {
	if p.enable {
		var err error
		p.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", p.port))
		if err != nil {
			return errors.Wrap(err, "Failed to start pprof server")
		}

		p.server = &http.Server{
			Addr: fmt.Sprintf(":%d", p.port),
		}
		err = p.server.Serve(p.listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Warn("Failed to start pprof server", "err", err, "port", p.port)
		}
	}

	return nil
}

func (p *Pprof) Stop() error {
	err := p.server.Close()
	if err != nil {
		log.Warn("Failed to stop pprof server", "err", err)
	}
	return nil
}

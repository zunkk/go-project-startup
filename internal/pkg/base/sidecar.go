package base

import (
	"github.com/zunkk/go-project-startup/internal/pkg/config"
	"github.com/zunkk/go-project-startup/pkg/frame"
)

func init() {
	frame.RegisterComponents(NewCustomSidecar)
}

type CustomSidecar struct {
	*frame.Sidecar
	Config *config.Config
}

func NewCustomSidecar(sidecar *frame.Sidecar, config *config.Config) (*CustomSidecar, error) {
	return &CustomSidecar{
		Sidecar: sidecar,
		Config:  config,
	}, nil
}

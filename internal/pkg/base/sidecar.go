package base

import (
	"github.com/zunkk/go-project-startup/internal/pkg/config"
	"github.com/zunkk/go-project-startup/pkg/frame"
	"github.com/zunkk/go-project-startup/pkg/repo"
)

func init() {
	frame.RegisterComponents(NewCustomSidecar)
}

type CustomSidecar struct {
	*frame.Sidecar
	Repo *repo.Repo[*config.Config]
}

func NewCustomSidecar(sidecar *frame.Sidecar, rep *repo.Repo[*config.Config]) (*CustomSidecar, error) {
	return &CustomSidecar{
		Sidecar: sidecar,
		Repo:    rep,
	}, nil
}

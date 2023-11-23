package base

import (
	"github.com/zunkk/go-project-startup/internal/pkg/config"
	"github.com/zunkk/go-project-startup/pkg/cache"
	"github.com/zunkk/go-project-startup/pkg/frame"
)

func init() {
	frame.RegisterComponents(NewCustomSidecar)
}

type CustomSidecar struct {
	*frame.Sidecar
	Config   *config.Config
	MemCache *cache.ExpiredMemCache
}

func NewCustomSidecar(sidecar *frame.Sidecar, config *config.Config) (*CustomSidecar, error) {
	memCache, err := cache.NewExpiredMemCache(config.Cache.ExpiredTime.ToDuration(), config.Cache.CleanupInterval.ToDuration())
	if err != nil {
		return nil, err
	}
	return &CustomSidecar{
		Sidecar:  sidecar,
		Config:   config,
		MemCache: memCache,
	}, nil
}

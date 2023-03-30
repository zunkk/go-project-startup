package base

import (
	"github.com/zunkk/go-project-startup/internal/pkg/config"
	"github.com/zunkk/go-project-startup/pkg/basic"
	"github.com/zunkk/go-project-startup/pkg/cache"
)

func init() {
	basic.RegisterComponents(NewBaseComponent)
}

type Component struct {
	*basic.BaseComponent
	Config   *config.Config
	MemCache *cache.MemCache
}

func NewBaseComponent(baseComponent *basic.BaseComponent, config *config.Config) (*Component, error) {
	memCache, err := cache.NewMemCache(config.Cache.ExpiredTime.ToDuration(), config.Cache.CleanupInterval.ToDuration())
	if err != nil {
		return nil, err
	}
	return &Component{
		BaseComponent: baseComponent,
		Config:        config,
		MemCache:      memCache,
	}, nil
}

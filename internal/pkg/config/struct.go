package config

import (
	"github.com/zunkk/go-project-startup/pkg/db"
	"github.com/zunkk/go-project-startup/pkg/repo"
)

type App struct {
	UUIDNodeIndex uint16 `mapstructure:"uuid_node_index" toml:"uuid_node_index"`
}

type Cache struct {
	ExpiredTime repo.Duration `mapstructure:"expired_time" toml:"expired_time"`
	Capacity    int           `mapstructure:"capacity" toml:"capacity"`
}

type DB struct {
	Type        db.Type `mapstructure:"type" toml:"type"`
	repo.DBInfo `mapstructure:",squash" toml:""`
}

type Config struct {
	App   App       `mapstructure:"app" toml:"app"`
	DB    DB        `mapstructure:"db" toml:"db"`
	HTTP  repo.HTTP `mapstructure:"http" toml:"http"`
	Cache Cache     `mapstructure:"cache" toml:"cache"`
	Log   repo.Log  `mapstructure:"log" toml:"log"`
}

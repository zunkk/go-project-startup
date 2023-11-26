package config

import (
	"github.com/zunkk/go-project-startup/pkg/config"
	"github.com/zunkk/go-project-startup/pkg/db"
)

type App struct {
	UUIDNodeIndex uint16 `mapstructure:"uuid_node_index" toml:"uuid_node_index"`
}

type Cache struct {
	ExpiredTime     config.Duration `mapstructure:"expired_time" toml:"expired_time"`
	CleanupInterval config.Duration `mapstructure:"cleanup_interval" toml:"cleanup_interval"`
}

type DB struct {
	Type          db.Type `mapstructure:"type" toml:"type"`
	config.DBInfo `mapstructure:",squash" toml:""`
}

type Config struct {
	RepoPath string      `mapstructure:"-" toml:"-"`
	App      App         `mapstructure:"app" toml:"app"`
	DB       DB          `mapstructure:"db" toml:"db"`
	HTTP     config.HTTP `mapstructure:"http" toml:"http"`
	Cache    Cache       `mapstructure:"cache" toml:"cache"`
	Log      config.Log  `mapstructure:"log" toml:"log"`
}

func (c *Config) GetRepoPath() string {
	return c.RepoPath
}

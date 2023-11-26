package config

import (
	"time"

	"github.com/zunkk/go-project-startup/pkg/config"
	"github.com/zunkk/go-project-startup/pkg/db"
	glog "github.com/zunkk/go-project-startup/pkg/log"
)

func DefaultConfig(repoPath string) *Config {
	return &Config{
		RepoPath: repoPath,
		App: App{
			UUIDNodeIndex: 0,
		},
		DB: DB{
			Type: db.DBTypePostgres,
			DBInfo: config.DBInfo{
				Host:     "127.0.0.1",
				Port:     5432,
				User:     "zunkk",
				Password: "zunkk",
				Schema:   "public",
				DBName:   "test",
				SSLMode:  "disable",
			},
		},
		HTTP: config.HTTP{
			Port:                  8080,
			MultipartMemory:       1024 * 1024 * 1024,
			ReadTimeout:           config.Duration(200 * time.Second),
			WriteTimeout:          config.Duration(200 * time.Second),
			TLSEnable:             false,
			TLSCertFilePath:       "",
			TLSKeyFilePath:        "",
			JWTTokenValidDuration: config.Duration(30 * time.Minute),
			JWTTokenHMACKey:       config.AppName + "_awsd_2024",
		},
		Cache: Cache{
			ExpiredTime:     config.Duration(24 * time.Hour),
			CleanupInterval: config.Duration(48 * time.Hour),
		},
		Log: config.Log{
			Level:            glog.LevelInfo,
			Filename:         config.AppName,
			MaxAge:           config.Duration(7 * 24 * time.Hour),
			MaxSizeStr:       "64mb",
			MaxSize:          64 << 20,
			RotationTime:     24 * config.Duration(time.Hour),
			EnableColor:      true,
			EnableCaller:     false,
			DisableTimestamp: false,
			ModuleLevelMap: map[string]glog.Level{
				"app":     glog.LevelInfo,
				"api":     glog.LevelDebug,
				"sidecar": glog.LevelDebug,
			},
		},
	}
}

package config

import (
	"time"

	"github.com/zunkk/go-sidecar/db"
	glog "github.com/zunkk/go-sidecar/log"
	"github.com/zunkk/go-sidecar/repo"
)

func DefaultConfig() *Config {
	return &Config{
		App: App{
			UUIDNodeIndex: 0,
		},
		DB: DB{
			Type: db.DBTypeSqlite,
			DBInfo: repo.DBInfo{
				Host:     "127.0.0.1",
				Port:     5432,
				User:     "zunkk",
				Password: "zunkk",
				Schema:   "public",
				DBName:   "test",
				SSLMode:  "disable",
			},
		},
		HTTP: repo.HTTP{
			Enable:                false,
			Port:                  8080,
			MultipartMemory:       1024 * 1024 * 1024,
			ReadTimeout:           repo.Duration(200 * time.Second),
			WriteTimeout:          repo.Duration(200 * time.Second),
			TLSEnable:             false,
			TLSCertFilePath:       "",
			TLSKeyFilePath:        "",
			JWTTokenValidDuration: repo.Duration(30 * time.Minute),
			JWTTokenHMACKey:       repo.AppName + "_awsd_2024",
		},
		Cache: Cache{
			ExpiredTime: repo.Duration(24 * time.Hour),
			Capacity:    10000,
		},
		Log: repo.Log{
			Level:            glog.LevelInfo,
			Filename:         repo.AppName,
			MaxAge:           repo.Duration(7 * 24 * time.Hour),
			MaxSizeStr:       "64mb",
			MaxSize:          64 << 20,
			RotationTime:     24 * repo.Duration(time.Hour),
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

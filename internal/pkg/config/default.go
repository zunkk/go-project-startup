package config

import (
	"time"

	"github.com/zunkk/go-project-startup/pkg/config"
)

func DefaultConfig(rootPath string) *Config {
	return &Config{
		RootPath: rootPath,
		App: App{
			NodeIndex: 0,
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
			Level:        "debug",
			Filename:     "app",
			MaxAge:       config.Duration(7 * 24 * time.Hour),
			MaxSize:      10 << 20,
			MaxSizeStr:   "10mb",
			RotationTime: config.Duration(time.Hour),
		},
	}
}

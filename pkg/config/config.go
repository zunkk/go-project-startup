package config

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/zunkk/go-project-startup/pkg/log"
	"github.com/zunkk/go-project-startup/pkg/util"
)

type CustomConfig interface {
	GetRootPath() string
}

func Load[T CustomConfig](defaultConfigFunc func(rootPath string) T) (t T, err error) {
	cfg, err := func() (T, error) {
		cfg := defaultConfigFunc(RootPath)
		existConfig := ExistConfigFile(cfg)
		if existConfig {
			if err := ReadConfig(cfg); err != nil {
				return t, err
			}
		}

		return cfg, nil
	}()
	if err != nil {
		return t, errors.Wrap(err, "failed to load config")
	}
	return cfg, nil
}

func ReadConfig[T CustomConfig](config T) error {
	viper.SetConfigFile(filepath.Join(config.GetRootPath(), cfgFileName))
	viper.SetConfigType("toml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix(AppName)
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	if err := viper.Unmarshal(config, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(";"),
		)),
	); err != nil {
		return err
	}

	return nil
}

func InitLogger(ctx context.Context, rootPath string, config Log) (*logrus.Logger, error) {
	logger, err := log.New(
		ctx,
		config.Level,
		filepath.Join(rootPath, logsDirName),
		config.Filename,
		config.MaxSize,
		config.MaxAge.ToDuration(),
		config.RotationTime.ToDuration(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init logger")
	}
	return logger, nil
}

func PrintSystemInfo(rootPath string, writer func(format string, args ...interface{})) {
	writer("%s version: %s", AppName, Version)
	writer("System version: %s", runtime.GOOS+"/"+runtime.GOARCH)
	writer("Golang version: %s", runtime.Version())
	writer("App build time: %s", BuildTime)
	writer("Git commit id: %s", CommitID)
	if rootPath != "" {
		writer("Config path: %s", rootPath)
	}
}

func WritePid(rootPath string) error {
	pid := os.Getpid()
	pidStr := strconv.Itoa(pid)
	if err := os.WriteFile(filepath.Join(rootPath, pidFileName), []byte(pidStr), 0755); err != nil {
		return errors.Wrap(err, "failed to write pid file")
	}
	return nil
}

func RemovePID(rootPath string) error {
	return os.Remove(filepath.Join(rootPath, pidFileName))
}

func WriteDebugInfo(rootPath string, debugInfo interface{}) error {
	p := filepath.Join(rootPath, debugFileName)
	_ = os.Remove(p)

	raw, err := json.Marshal(debugInfo)
	if err != nil {
		return err
	}
	if err := os.WriteFile(p, raw, 0755); err != nil {
		return errors.Wrap(err, "failed to write debug info file")
	}
	return nil
}

func ExistConfigFile[T CustomConfig](config T) bool {
	return util.FileExist(filepath.Join(config.GetRootPath(), cfgFileName))
}

func WriteConfig[T CustomConfig](config T) error {
	raw, err := MarshalConfig(config)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(config.GetRootPath(), cfgFileName), []byte(raw), 0755); err != nil {
		return err
	}
	return nil
}

func MarshalConfig[T CustomConfig](config T) (string, error) {
	buf := bytes.NewBuffer([]byte{})
	e := toml.NewEncoder(buf)
	e.SetIndentTables(true)
	e.SetArraysMultiline(true)
	err := e.Encode(config)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

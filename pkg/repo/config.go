package repo

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"

	"github.com/zunkk/go-project-startup/pkg/util"
)

type CustomConfig any

func ReadConfig[T CustomConfig](repoPath string, config T) error {
	vp := viper.New()
	vp.AutomaticEnv()
	vp.SetEnvPrefix(EnvPrefix)
	vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if ExistConfigFile(repoPath) {
		vp.SetConfigFile(filepath.Join(repoPath, cfgFileName))
		vp.SetConfigType("toml")
		err := vp.ReadInConfig()
		if err != nil {
			return err
		}
	}

	if err := vp.Unmarshal(config, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			StringToTimeDurationHookFunc(),
			StringToLevelHookFunc(),
			mapstructure.StringToSliceHookFunc(";"),
		)),
	); err != nil {
		return err
	}

	return nil
}

func ExistConfigFile(repoPath string) bool {
	return util.FileExist(filepath.Join(repoPath, cfgFileName))
}

func WriteConfig[T CustomConfig](repoPath string, config T) error {
	raw, err := MarshalConfig(config)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(repoPath, cfgFileName), []byte(raw), 0755); err != nil {
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

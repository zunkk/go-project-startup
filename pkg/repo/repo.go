package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"

	glog "github.com/zunkk/go-project-startup/pkg/log"
	"github.com/zunkk/go-project-startup/pkg/util"
)

type Repo[T CustomConfig] struct {
	RepoPath string
	Cfg      T
}

func Load[T CustomConfig](repoPath string, defaultConfigFunc func() T) (rep *Repo[T], err error) {
	cfg := defaultConfigFunc()
	if err := ReadConfig(repoPath, cfg); err != nil {
		return nil, errors.Wrap(err, "failed to load config")
	}
	return &Repo[T]{
		RepoPath: RootPath,
		Cfg:      cfg,
	}, nil
}

func InitLogger(ctx context.Context, repoPath string, config Log) error {
	err := glog.Init(
		ctx,
		config.Level,
		filepath.Join(repoPath, logsDirName),
		config.Filename,
		config.MaxSize,
		config.MaxAge.ToDuration(),
		config.RotationTime.ToDuration(),
		config.EnableColor,
		config.EnableCaller,
		config.DisableTimestamp,
		config.ModuleLevelMap,
	)
	if err != nil {
		return errors.Wrap(err, "failed to init logger")
	}
	return nil
}

func PrintSystemInfo(repoPath string, writer func(format string, args ...any)) {
	writer("%s version: %s", AppName, Version)
	writer("System version: %s", runtime.GOOS+"/"+runtime.GOARCH)
	writer("Golang version: %s", runtime.Version())
	writer("App build time: %s", BuildTime)
	writer("Git commit id: %s", CommitID)
	if repoPath != "" {
		writer("Repo path: %s", repoPath)
	}
}

func WriteDebugInfo(repoPath string, debugInfo any) error {
	p := filepath.Join(repoPath, debugFileName)
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

func LoadEnvFile() {
	if util.FileExist(EnvFilePath) {
		if err := godotenv.Load(EnvFilePath); err != nil {
			fmt.Printf("load env file %s failed: %s\n", EnvFilePath, err)
			return
		}
	}
}

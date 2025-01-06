package repo

import "strings"

var (
	RootPath = ""

	EnvFilePath = ".env"

	JWTTokenHeaderKey = "token"
)

var (
	AppName = ""

	AppDesc = ""

	Version = ""

	BuildTime = ""

	CommitID = ""

	EnvPrefix = strings.ToUpper(AppName)
)

const (
	cfgFileName = "config.toml"

	debugFileName = "debug-info.json"

	pidFileName = "process.pid"

	logsDirName = "logs"

	IPCFileName = "ipc.sock"
)

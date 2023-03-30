package config

var (
	RootPath          = ""
	JWTTokenHeaderKey = "token"
)

var (
	AppName = ""

	AppDesc = ""

	Version = ""

	BuildTime = ""

	CommitID = ""
)

const (
	cfgFileName = "config.toml"

	pidFileName = "process.pid"

	debugFileName = "debug-info.json"

	logsDirName = "logs"
)

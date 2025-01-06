package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	sloglogrus "github.com/samber/slog-logrus/v2"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Level int

var (
	LevelPanic = Level(24)
	LevelFatal = Level(16)
	LevelError = Level(slog.LevelError)
	LevelWarn  = Level(slog.LevelWarn)
	LevelInfo  = Level(slog.LevelInfo)
	LevelDebug = Level(slog.LevelDebug)
	LevelTrace = Level(-8)
)

func (l Level) MarshalText() (text []byte, err error) {
	levelStr := l.String()
	if levelStr == "" {
		return nil, fmt.Errorf("not a valid Level: %d", l)
	}
	return []byte(levelStr), nil
}

func (l *Level) UnmarshalText(b []byte) error {
	x, err := func() (Level, error) {
		lStr := string(b)
		switch strings.ToLower(lStr) {
		case "panic":
			return LevelPanic, nil
		case "fatal":
			return LevelFatal, nil
		case "error":
			return LevelError, nil
		case "warn", "warning":
			return LevelWarn, nil
		case "info":
			return LevelInfo, nil
		case "debug":
			return LevelDebug, nil
		case "trace":
			return LevelTrace, nil
		}

		var l Level
		return l, fmt.Errorf("not a valid Level string: %s", lStr)
	}()
	if err != nil {
		return err
	}
	*l = x
	return nil
}

func (l *Level) String() string {
	switch *l {
	case LevelPanic:
		return "panic"
	case LevelFatal:
		return "fatal"
	case LevelError:
		return "error"
	case LevelWarn:
		return "warn"
	case LevelInfo:
		return "info"
	case LevelDebug:
		return "debug"
	case LevelTrace:
		return "trace"
	}
	return ""
}

var globalModuleLevelMap = map[string]slog.Level{}

var globalLogrusLogger = logrus.New()

var globalSlog = WithModule("")

func init() {
	globalLogrusLogger.SetFormatter(&Formatter{
		FirstFieldsOrder: []string{"uri", "err_code", "err_msg", "err"},
		LastFieldsOrder:  []string{"method", "ip", "http_code", "req_id", "time_cost", "caller"},
		TimestampFormat:  "01/02 15:04:05.000",
		EnableColor:      true,
		EnableCaller:     false,
		DisableTimestamp: false,
	})
}

func Disable() {
	globalLogrusLogger.SetOutput(io.Discard)
}

func Default() *slog.Logger {
	return globalSlog
}

func WithModule(module string) *slog.Logger {
	return slog.New(&ModuleLevel{
		module: module,
		Handler: sloglogrus.Option{
			Level:  slog.Level(-1000),
			Logger: globalLogrusLogger,
		}.NewLogrusHandler(),
	})
}

func Init(ctx context.Context, level Level, filePath string, fileName string, maxSize int64, maxAge time.Duration, rotationTime time.Duration, enableColor bool, enableCaller bool, disableTimestamp bool, moduleLevelMap map[string]Level) error {
	logrusLogger, err := newLogrusLogger(ctx, level, filePath, fileName, maxSize, maxAge, rotationTime, enableColor, enableCaller, disableTimestamp)
	if err != nil {
		return err
	}
	globalLogrusLogger.Formatter = logrusLogger.Formatter
	globalLogrusLogger.ReportCaller = logrusLogger.ReportCaller
	globalLogrusLogger.Level = logrusLogger.Level
	globalLogrusLogger.Out = logrusLogger.Out
	globalLogrusLogger.Hooks = logrusLogger.Hooks
	for m, l := range moduleLevelMap {
		globalModuleLevelMap[m] = slog.Level(l)
	}
	slog.SetDefault(globalSlog)
	return nil
}

func SetModuleLevel(module string, level Level) {
	globalModuleLevelMap[module] = slog.Level(level)
}

func newLogrusLogger(ctx context.Context, level Level, filePath string, fileName string, maxSize int64, maxAge time.Duration, rotationTime time.Duration, enableColor bool, enableCaller bool, disableTimestamp bool) (*logrus.Logger, error) {
	logger := logrus.New()
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		return nil, err
	}
	logger.SetFormatter(&Formatter{
		FirstFieldsOrder: []string{"uri", "err_code", "err_msg"},
		LastFieldsOrder:  []string{"method", "ip", "http_code", "req_id", "time_cost", "caller"},
		TimestampFormat:  "01/02 15:04:05.000",
		EnableColor:      enableColor,
		EnableCaller:     enableCaller,
		DisableTimestamp: disableTimestamp,
	})
	logger.SetReportCaller(true)
	lvl, err := logrus.ParseLevel(level.String())
	if err != nil {
		lvl = logrus.ErrorLevel
	}
	logger.SetLevel(lvl)
	logger.SetOutput(os.Stdout)

	h, err := newRotateHook(ctx, &lumberjack.Logger{
		Filename:  filepath.Join(filePath, fileName) + ".log",
		MaxSize:   int(maxSize),
		MaxAge:    int(maxAge),
		LocalTime: true,
		Compress:  false,
	}, rotationTime)
	if err != nil {
		return nil, err
	}
	logger.AddHook(h)

	if err := redirectPanic(filepath.Join(filePath, "error.log")); err != nil {
		return nil, errors.Wrap(err, "failed to redirect panic")
	}

	return logger, nil
}

func New(ctx context.Context, level Level, filePath string, fileName string, maxSize int64, maxAge time.Duration, rotationTime time.Duration, enableColor bool, enableCaller bool, disableTimestamp bool) (*slog.Logger, error) {
	logrusLogger, err := newLogrusLogger(ctx, level, filePath, fileName, maxSize, maxAge, rotationTime, enableColor, enableCaller, disableTimestamp)
	if err != nil {
		return nil, err
	}
	return slog.New(sloglogrus.Option{
		Level:  slog.LevelDebug,
		Logger: logrusLogger,
	}.NewLogrusHandler()), nil
}

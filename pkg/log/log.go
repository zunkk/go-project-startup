package log

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func New(ctx context.Context, level string, filePath string, fileName string, maxSize int64, maxAge time.Duration, rotationTime time.Duration) (*logrus.Logger, error) {
	log := logrus.New()
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		return nil, err
	}
	log.SetFormatter(&Formatter{
		FirstFieldsOrder: []string{"uri", "err_code", "err_msg"},
		LastFieldsOrder:  []string{"method", "ip", "http_code", "req_id", "time_cost", "caller"},
		TimestampFormat:  "2006/01/02 15:04:05.000",
	})
	log.SetReportCaller(true)
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.ErrorLevel
	}
	log.SetLevel(lvl)
	log.SetOutput(os.Stdout)

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
	log.AddHook(h)
	return log, nil
}

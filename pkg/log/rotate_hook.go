package log

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type RotateHook struct {
	logger *lumberjack.Logger
}

func newRotateHook(ctx context.Context, logger *lumberjack.Logger, rotationTime time.Duration) (*RotateHook, error) {
	if logger == nil {
		return nil, errors.New("Logger cannot be nil")
	}
	go func() {
		_ = logger.Rotate()

		tk := time.NewTicker(rotationTime)
		select {
		case <-tk.C:
			_ = logger.Rotate()
		case <-ctx.Done():
			return
		}
		for {
			nowTime := time.Now()
			nowTimeStr := nowTime.Format("2006-01-02")
			t2, _ := time.ParseInLocation("2006-01-02", nowTimeStr, time.Local)
			next := t2.AddDate(0, 0, 1)
			after := next.UnixNano() - nowTime.UnixNano() - 1
			<-time.After(time.Duration(after) * time.Nanosecond)
		}
	}()

	return &RotateHook{
		logger: logger,
	}, nil
}

func (hook *RotateHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	_, err = hook.logger.Write([]byte(line))
	if err != nil {
		return err
	}

	return nil
}

func (hook *RotateHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.TraceLevel,
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

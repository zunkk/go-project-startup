package log

import (
	"fmt"
	"os"
	"time"
)

var stdErrFile *os.File

func redirectPanic(errLogFilePath string) error {
	var err error
	stdErrFile, err = os.OpenFile(errLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(stdErrFile, "\n\n-------------------- process start time: %s --------------------\n", time.Now().Format("2006-01-02 15:04:05")); err != nil {
		return err
	}
	redirectStderr(stdErrFile)
	return nil
}

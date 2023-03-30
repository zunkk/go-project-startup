package util

import (
	"time"

	"github.com/pkg/errors"
)

var jsonTimeStrFormats = []string{"2006", "2006-01", "2006-01-02"}

func StringToDate(data string) (time.Time, error) {
	for _, format := range jsonTimeStrFormats {
		if date, err := time.ParseInLocation(format, data, time.Local); err == nil {
			return date, nil
		}
	}
	return time.Time{}, errors.Errorf("date time format error, expected: %v, get: %v", jsonTimeStrFormats, string(data))
}

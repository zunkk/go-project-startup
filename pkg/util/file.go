package util

import (
	"os"
)

func FileExist(path string) bool {
	fi, err := os.Lstat(path)
	if fi != nil || (err != nil && !os.IsNotExist(err)) {
		return true
	}
	return false
}

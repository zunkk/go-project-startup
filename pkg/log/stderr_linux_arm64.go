// Log the panic under unix to the log file

//go:build linux && arm64
// +build linux,arm64

package log

import (
	"log"
	"os"
	"syscall"
)

// redirectStderr to the file passed in
func redirectStderr(f *os.File) {
	err := syscall.Dup3(int(f.Fd()), int(os.Stderr.Fd()), 0)
	if err != nil {
		log.Fatalf("Failed to redirect stderr to file: %v", err)
	}
}

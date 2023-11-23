package coreapi

import (
	"github.com/zunkk/go-project-startup/pkg/frame"
	"github.com/zunkk/go-project-startup/pkg/mutex"
)

func init() {
	frame.RegisterComponents(NewCoreAPI, mutex.NewKeyMutex)
}

type CoreAPI struct {
}

func NewCoreAPI() (*CoreAPI, error) {
	return &CoreAPI{}, nil
}

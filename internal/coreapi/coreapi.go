package coreapi

import (
	"github.com/zunkk/go-project-startup/pkg/basic"
	"github.com/zunkk/go-project-startup/pkg/mutex"
)

func init() {
	basic.RegisterComponents(NewCoreAPI, mutex.NewKeyMutex)
}

type CoreAPI struct {
}

func NewCoreAPI() (*CoreAPI, error) {
	return &CoreAPI{}, nil
}

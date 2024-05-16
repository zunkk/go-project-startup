package coreapi

import (
	"github.com/zunkk/go-project-startup/pkg/frame"
	"github.com/zunkk/go-project-startup/pkg/mutex"
)

func init() {
	frame.RegisterComponents(NewCoreAPI, mutex.NewKeyMutex)
}

type CoreAPI struct {
	// UserSrv *service.User
}

func NewCoreAPI() (*CoreAPI, error) {
	return &CoreAPI{}, nil
}

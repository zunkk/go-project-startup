package coreapi

import (
	"github.com/zunkk/go-project-startup/internal/core/service"
	"github.com/zunkk/go-sidecar/frame"
	"github.com/zunkk/go-sidecar/mutex"
)

func init() {
	frame.RegisterComponents(NewCoreAPI, mutex.NewKeyMutex)
}

type CoreAPI struct {
	UserService *service.UserService
}

func NewCoreAPI(userSrv *service.UserService) (*CoreAPI, error) {
	return &CoreAPI{
		UserService: userSrv,
	}, nil
}

package coreapi

import (
	"github.com/zunkk/go-project-startup/internal/core/service"
	"github.com/zunkk/go-project-startup/pkg/frame"
	"github.com/zunkk/go-project-startup/pkg/mutex"
)

func init() {
	frame.RegisterComponents(NewCoreAPI, mutex.NewKeyMutex)
}

type CoreAPI struct {
	UserSrv *service.User
}

func NewCoreAPI(userSrv *service.User) (*CoreAPI, error) {
	return &CoreAPI{
		UserSrv: userSrv,
	}, nil
}

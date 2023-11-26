package service

import (
	"context"

	"github.com/zunkk/go-project-startup/internal/core/dao"
	"github.com/zunkk/go-project-startup/internal/core/model"
	"github.com/zunkk/go-project-startup/internal/pkg/base"
)

type User struct {
	sidecar *base.CustomSidecar
	userDao *dao.User
}

func NewUser(sidecar *base.CustomSidecar, userDao *dao.User) (*User, error) {
	return &User{
		sidecar: sidecar,
		userDao: userDao,
	}, nil
}

func (d *User) QueryByID(ctx context.Context, id int64) (*model.User, error) {
	return d.userDao.QueryByID(ctx, id)
}

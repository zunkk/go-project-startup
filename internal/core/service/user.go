package service

import (
	"context"

	"github.com/stephenafamo/bob"

	"github.com/zunkk/go-project-startup/internal/core/dao"
	"github.com/zunkk/go-project-startup/internal/core/model"
	"github.com/zunkk/go-project-startup/internal/pkg/base"
)

type UserService struct {
	sidecar *base.CustomSidecar
	db      *bob.DB
}

func NewUserService(sidecar *base.CustomSidecar, sqlConnector *dao.SQLConnector) (*UserService, error) {
	return &UserService{
		sidecar: sidecar,
		db:      sqlConnector.DB,
	}, nil
}

func (d *UserService) QueryByID(ctx context.Context, id int64) (*model.User, error) {
	return model.FindUser(ctx, d.db, id)
}

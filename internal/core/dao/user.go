package dao

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/zunkk/go-project-startup/internal/core/model"
	"github.com/zunkk/go-project-startup/internal/pkg/base"
)

type User struct {
	sidecar *base.CustomSidecar
	db      *sqlx.DB
}

func NewUser(sidecar *base.CustomSidecar, sqlConnector *SQLConnector) (*User, error) {
	return &User{
		sidecar: sidecar,
		db:      sqlConnector.DB,
	}, nil
}

func (d *User) QueryByID(ctx context.Context, id int) (*model.User, error) {
	return model.FindUser(ctx, d.db, id)
}

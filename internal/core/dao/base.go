package dao

import (
	"github.com/jmoiron/sqlx"

	"github.com/zunkk/go-project-startup/internal/pkg/base"
	"github.com/zunkk/go-project-startup/pkg/db/sql"
	"github.com/zunkk/go-project-startup/pkg/frame"
)

func init() {
	frame.RegisterComponents(NewSQLConnector)
}

type SQLConnector struct {
	DB *sqlx.DB
}

func NewSQLConnector(sidecar *base.CustomSidecar) (*SQLConnector, error) {
	sqlDB, err := sql.Open(sidecar.Repo.Cfg.DB.Type, sidecar.Repo.RepoPath, sidecar.Repo.Cfg.DB.DBInfo)
	if err != nil {
		return nil, err
	}
	return &SQLConnector{DB: sqlDB}, nil
}

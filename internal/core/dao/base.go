package dao

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/zunkk/go-project-startup/internal/pkg/base"
	"github.com/zunkk/go-project-startup/pkg/db/sql"
	"github.com/zunkk/go-project-startup/pkg/frame"
	glog "github.com/zunkk/go-project-startup/pkg/log"
)

var log = glog.WithModule("db")

func init() {
	frame.RegisterComponents(NewSQLConnector)
}

type DBAction func(dbTX boil.ContextExecutor) error

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

func (t *SQLConnector) SubmitDBChangesByTransaction(dbActions ...DBAction) error {
	dbTX, err := t.DB.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to begin db transaction")
	}
	// The rollback will be ignored if the tx has been committed later in the function.
	defer func() {
		if err != nil {
			if rollbackErr := dbTX.Rollback(); rollbackErr != nil {
				log.Warn("Failed to rollback", "err", rollbackErr)
			}
		}
	}()
	for _, dbAction := range dbActions {
		if err := dbAction(dbTX); err != nil {
			return err
		}
	}
	if err := dbTX.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit db transaction")
	}
	return nil
}

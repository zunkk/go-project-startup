package dao

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/zunkk/go-project-startup/build"
	"github.com/zunkk/go-project-startup/internal/pkg/base"
	"github.com/zunkk/go-sidecar/db/sql"
	"github.com/zunkk/go-sidecar/frame"
	glog "github.com/zunkk/go-sidecar/log"
)

var log = glog.WithModule("db")

func init() {
	frame.RegisterComponents(NewSQLConnector)
}

type DBAction func(dbTX boil.ContextExecutor) error

type SQLConnector struct {
	sidecar *base.CustomSidecar
	DB      *sqlx.DB
}

func NewSQLConnector(sidecar *base.CustomSidecar) (*SQLConnector, error) {
	sqlDB, err := sql.Open(sidecar.Repo.Cfg.DB.Type, sidecar.Repo.RepoPath, sidecar.Repo.Cfg.DB.DBInfo)
	if err != nil {
		return nil, err
	}
	sqlConnector := &SQLConnector{
		sidecar: sidecar,
		DB:      sqlDB,
	}
	sidecar.RegisterLifecycleHook(sqlConnector)
	return sqlConnector, nil
}

func NewSQLConnectorWithDB(sidecar *base.CustomSidecar, db *sqlx.DB) (*SQLConnector, error) {
	sqlConnector := &SQLConnector{
		sidecar: sidecar,
		DB:      db,
	}
	sidecar.RegisterLifecycleHook(sqlConnector)
	return sqlConnector, nil
}

func (c *SQLConnector) ComponentName() string {
	return "sql-connector"
}

func (c *SQLConnector) Start() error {
	return build.TryCreateDDLTables(c.sidecar.Ctx, c.DB)
}

func (c *SQLConnector) Stop() error {
	return nil
}

func (c *SQLConnector) SubmitDBChangesByTransaction(dbActions ...DBAction) error {
	dbTX, err := c.DB.Begin()
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

package build

import (
	"context"
	_ "embed"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	//go:embed ddl.sql
	DDL string
)

func TryCreateDDLTables(ctx context.Context, db *sqlx.DB) error {
	_, err := db.ExecContext(ctx, DDL)
	if err != nil {
		return errors.Wrap(err, "failed to create ddl tables")
	}
	return nil
}

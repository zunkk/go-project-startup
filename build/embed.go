package build

import (
	"context"
	_ "embed"

	"github.com/pkg/errors"
	"github.com/stephenafamo/bob"
)

//go:embed ddl.sql
var DDL string

func TryCreateDDLTables(ctx context.Context, db *bob.DB) error {
	_, err := db.ExecContext(ctx, DDL)
	if err != nil {
		return errors.Wrap(err, "failed to create ddl tables")
	}
	return nil
}

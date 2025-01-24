package memory

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/marcboeker/go-duckdb"
)

func OpenSQLDB() (*sqlx.DB, error) {
	return sqlx.Connect("duckdb", "")
}

package sql

import (
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/zunkk/go-project-startup/pkg/config"
	"github.com/zunkk/go-project-startup/pkg/db"
)

var dbType2DriverName = map[db.Type]string{
	db.DBTypePostgres: "pgx",
	db.DBTypeMysql:    "mysql",
	db.DBTypeSqlite:   "sqlite3",
}

var dbType2DSNGenerator = map[db.Type]func(repoPath string, info config.DBInfo) string{
	db.DBTypePostgres: func(repoPath string, info config.DBInfo) string {
		return fmt.Sprintf("user=%s password=%s host=%s port=%d database=%s sslmode=%s", info.User, info.Password, info.Host, info.Port, info.DBName, info.SSLMode)
	},
	db.DBTypeMysql: func(repoPath string, info config.DBInfo) string {
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", info.User, info.Password, info.Host, info.Port, info.DBName)
	},
	db.DBTypeSqlite: func(repoPath string, info config.DBInfo) string {
		dbDir := filepath.Join(repoPath, "db")
		if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
			panic(err)
		}
		dbFilePath := filepath.Join(dbDir, "sqlite.db")
		return fmt.Sprintf("file:%s?cache=shared&mode=rwc", dbFilePath)
	},
}

func Open(dbType db.Type, repoPath string, info config.DBInfo) (*sqlx.DB, error) {
	if _, ok := dbType2DSNGenerator[dbType]; !ok {
		return nil, fmt.Errorf("unsupported db type: %s", dbType)
	}
	return sqlx.Connect(dbType2DriverName[dbType], dbType2DSNGenerator[dbType](repoPath, info))
}

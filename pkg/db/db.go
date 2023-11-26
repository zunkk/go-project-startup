package db

type Type string

const (
	DBTypePostgres Type = "postgres"
	DBTypeMysql    Type = "mysql"
	DBTypeSqlite   Type = "sqlite"
	DBTypeMongo    Type = "mongo"
)

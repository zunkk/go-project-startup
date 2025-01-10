package db

type Type string

const (
	DBTypePostgres  Type = "postgres"
	DBTypeMysql     Type = "mysql"
	DBTypeSqlite    Type = "sqlite"
	DBTypeSqlMemory Type = "sql_memory"
	DBTypeMongo     Type = "mongo"
)

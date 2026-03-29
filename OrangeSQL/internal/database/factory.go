package database

import "fmt"

// SupportedDrivers は対応するドライバ名の一覧。
var SupportedDrivers = []string{"sqlite", "postgres", "mysql", "sqlserver"}

// New はドライバ名と接続パラメータから Database 実装を生成するファクトリ。
func New(driver string, params ConnectionParams) (Database, error) {
	switch driver {
	case "sqlite":
		return NewSQLite(params.Path)
	case "postgres":
		return NewPostgres(params)
	case "mysql":
		return NewMySQL(params)
	case "sqlserver":
		return NewSQLServer(params)
	default:
		return nil, fmt.Errorf("unsupported driver: %s", driver)
	}
}

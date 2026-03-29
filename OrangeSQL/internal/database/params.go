package database

import "fmt"

// ConnectionParams は各 DB ドライバの接続パラメータ。
type ConnectionParams struct {
	// SQLite
	Path string `json:"path,omitempty"`

	// TCP 系共通（PostgreSQL / MySQL / SQL Server）
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	DBName   string `json:"dbName,omitempty"`

	// SSL/TLS
	SSLMode string `json:"sslMode,omitempty"`
}

// DefaultPort はドライバのデフォルトポートを返す。
func DefaultPort(driver string) int {
	switch driver {
	case "postgres":
		return 5432
	case "mysql":
		return 3306
	case "sqlserver":
		return 1433
	default:
		return 0
	}
}

// effectivePort は Port が 0 の場合にデフォルトポートを返す。
func (p ConnectionParams) effectivePort(driver string) int {
	if p.Port > 0 {
		return p.Port
	}
	return DefaultPort(driver)
}

// PostgresDSN は PostgreSQL の接続文字列を生成する。
func (p ConnectionParams) PostgresDSN() string {
	port := p.effectivePort("postgres")
	sslMode := p.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host, port, p.User, p.Password, p.DBName, sslMode,
	)
}

// MySQLDSN は MySQL の接続文字列を生成する。
func (p ConnectionParams) MySQLDSN() string {
	port := p.effectivePort("mysql")
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		p.User, p.Password, p.Host, port, p.DBName,
	)
}

// SQLServerDSN は SQL Server の接続文字列を生成する。
func (p ConnectionParams) SQLServerDSN() string {
	port := p.effectivePort("sqlserver")
	return fmt.Sprintf(
		"sqlserver://%s:%s@%s:%d?database=%s",
		p.User, p.Password, p.Host, port, p.DBName,
	)
}

package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySQL ドライバ登録
)

// MySQL は MySQL の Database 実装。
type MySQL struct {
	db   *sql.DB
	name string
}

// NewMySQL は MySQL に接続する。
func NewMySQL(params ConnectionParams) (*MySQL, error) {
	dsn := params.MySQLDSN()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("mysql open: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("mysql ping: %w", err)
	}
	port := params.effectivePort("mysql")
	name := fmt.Sprintf("%s:%d/%s", params.Host, port, params.DBName)
	return &MySQL{db: db, name: name}, nil
}

func (m *MySQL) Name() string  { return m.name }
func (m *MySQL) Close() error  { return m.db.Close() }

func (m *MySQL) Query(q string) (QueryResult, error) {
	rows, err := m.db.Query(q)
	if err != nil {
		return QueryResult{}, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func (m *MySQL) Exec(q string) (ExecResult, error) {
	return execStatement(m.db, q)
}

func (m *MySQL) ExecBatch(statements []string) (int64, error) {
	return execBatch(m.db, statements)
}

func (m *MySQL) Tables() ([]TableInfo, error) {
	rows, err := m.db.Query(`
		SELECT table_name,
		       CASE table_type WHEN 'BASE TABLE' THEN 'table' ELSE 'view' END
		FROM information_schema.tables
		WHERE table_schema = DATABASE()
		ORDER BY table_name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []TableInfo
	for rows.Next() {
		var t TableInfo
		if err := rows.Scan(&t.Name, &t.Type); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	return tables, rows.Err()
}

func (m *MySQL) Columns(tableName string) ([]ColumnInfo, error) {
	rows, err := m.db.Query(`
		SELECT
			COLUMN_NAME,
			COLUMN_TYPE,
			CASE WHEN IS_NULLABLE = 'NO' THEN true ELSE false END,
			CASE WHEN COLUMN_KEY = 'PRI' THEN true ELSE false END
		FROM information_schema.columns
		WHERE table_schema = DATABASE() AND table_name = ?
		ORDER BY ORDINAL_POSITION
	`, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var col ColumnInfo
		if err := rows.Scan(&col.Name, &col.Type, &col.NotNull, &col.PK); err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(columns) == 0 {
		return nil, fmt.Errorf("table not found: %s", tableName)
	}
	return columns, nil
}

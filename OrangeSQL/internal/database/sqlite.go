package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // SQLite ドライバ登録
)

// SQLite は modernc.org/sqlite を使った Database 実装。
type SQLite struct {
	db   *sql.DB
	name string
}

// NewSQLite は指定パスの SQLite ファイルに接続する。
// ファイルが存在しない場合は SQLite が自動作成する。
func NewSQLite(path string) (*SQLite, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("sqlite open: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("sqlite ping: %w", err)
	}
	return &SQLite{db: db, name: path}, nil
}

func (s *SQLite) Name() string {
	return s.name
}

func (s *SQLite) Close() error {
	return s.db.Close()
}

func (s *SQLite) Query(q string) (QueryResult, error) {
	rows, err := s.db.Query(q)
	if err != nil {
		return QueryResult{}, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func (s *SQLite) Exec(q string) (ExecResult, error) {
	return execStatement(s.db, q)
}

func (s *SQLite) Tables() ([]TableInfo, error) {
	rows, err := s.db.Query(
		"SELECT name, type FROM sqlite_master WHERE type IN ('table', 'view') ORDER BY name",
	)
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

func (s *SQLite) Columns(tableName string) ([]ColumnInfo, error) {
	rows, err := s.db.Query(fmt.Sprintf("PRAGMA table_info(%q)", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var cid int
		var name, colType string
		var notNull, pk int
		var dfltValue interface{}
		if err := rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk); err != nil {
			return nil, err
		}
		columns = append(columns, ColumnInfo{
			Name:    name,
			Type:    colType,
			PK:      pk > 0,
			NotNull: notNull > 0,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(columns) == 0 {
		return nil, fmt.Errorf("table not found: %s", tableName)
	}
	return columns, nil
}

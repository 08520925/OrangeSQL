package database

import (
	"database/sql"
	"fmt"

	_ "github.com/microsoft/go-mssqldb" // SQL Server ドライバ登録
)

// SQLServer は SQL Server の Database 実装。
type SQLServer struct {
	db   *sql.DB
	name string
}

// NewSQLServer は SQL Server に接続する。
func NewSQLServer(params ConnectionParams) (*SQLServer, error) {
	dsn := params.SQLServerDSN()
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlserver open: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("sqlserver ping: %w", err)
	}
	port := params.effectivePort("sqlserver")
	name := fmt.Sprintf("%s:%d/%s", params.Host, port, params.DBName)
	return &SQLServer{db: db, name: name}, nil
}

func (s *SQLServer) Name() string  { return s.name }
func (s *SQLServer) Close() error  { return s.db.Close() }

func (s *SQLServer) Query(q string) (QueryResult, error) {
	rows, err := s.db.Query(q)
	if err != nil {
		return QueryResult{}, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func (s *SQLServer) Exec(q string) (ExecResult, error) {
	return execStatement(s.db, q)
}

func (s *SQLServer) ExecBatch(statements []string) (int64, error) {
	return execBatch(s.db, statements)
}

func (s *SQLServer) Tables() ([]TableInfo, error) {
	rows, err := s.db.Query(`
		SELECT TABLE_NAME,
		       CASE TABLE_TYPE WHEN 'BASE TABLE' THEN 'table' ELSE 'view' END
		FROM INFORMATION_SCHEMA.TABLES
		ORDER BY TABLE_NAME
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

func (s *SQLServer) Columns(tableName string) ([]ColumnInfo, error) {
	rows, err := s.db.Query(`
		SELECT
			c.COLUMN_NAME,
			c.DATA_TYPE,
			CASE WHEN c.IS_NULLABLE = 'NO' THEN 1 ELSE 0 END,
			CASE WHEN kcu.COLUMN_NAME IS NOT NULL THEN 1 ELSE 0 END,
			COALESCE(ep.value, '')
		FROM INFORMATION_SCHEMA.COLUMNS c
		LEFT JOIN INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc
			ON tc.TABLE_NAME = c.TABLE_NAME
			AND tc.TABLE_SCHEMA = c.TABLE_SCHEMA
			AND tc.CONSTRAINT_TYPE = 'PRIMARY KEY'
		LEFT JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
			ON kcu.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
			AND kcu.TABLE_SCHEMA = tc.TABLE_SCHEMA
			AND kcu.COLUMN_NAME = c.COLUMN_NAME
		LEFT JOIN sys.extended_properties ep
			ON ep.major_id = OBJECT_ID(c.TABLE_SCHEMA + '.' + c.TABLE_NAME)
			AND ep.minor_id = c.ORDINAL_POSITION
			AND ep.name = 'MS_Description'
		WHERE c.TABLE_NAME = @p1
		ORDER BY c.ORDINAL_POSITION
	`, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var col ColumnInfo
		if err := rows.Scan(&col.Name, &col.Type, &col.NotNull, &col.PK, &col.Comment); err != nil {
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

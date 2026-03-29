package database

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL ドライバ登録
)

// Postgres は PostgreSQL の Database 実装。
type Postgres struct {
	db   *sql.DB
	name string
}

// NewPostgres は PostgreSQL に接続する。
func NewPostgres(params ConnectionParams) (*Postgres, error) {
	dsn := params.PostgresDSN()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres open: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("postgres ping: %w", err)
	}
	port := params.effectivePort("postgres")
	name := fmt.Sprintf("%s:%d/%s", params.Host, port, params.DBName)
	return &Postgres{db: db, name: name}, nil
}

func (p *Postgres) Name() string  { return p.name }
func (p *Postgres) Close() error  { return p.db.Close() }

func (p *Postgres) Query(q string) (QueryResult, error) {
	rows, err := p.db.Query(q)
	if err != nil {
		return QueryResult{}, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func (p *Postgres) Exec(q string) (ExecResult, error) {
	return execStatement(p.db, q)
}

func (p *Postgres) Tables() ([]TableInfo, error) {
	rows, err := p.db.Query(`
		SELECT table_name,
		       CASE table_type WHEN 'BASE TABLE' THEN 'table' ELSE 'view' END
		FROM information_schema.tables
		WHERE table_schema = 'public'
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

func (p *Postgres) Columns(tableName string) ([]ColumnInfo, error) {
	rows, err := p.db.Query(`
		SELECT
			c.column_name,
			c.data_type,
			CASE WHEN c.is_nullable = 'NO' THEN true ELSE false END,
			CASE WHEN kcu.column_name IS NOT NULL THEN true ELSE false END
		FROM information_schema.columns c
		LEFT JOIN information_schema.table_constraints tc
			ON tc.table_schema = c.table_schema
			AND tc.table_name = c.table_name
			AND tc.constraint_type = 'PRIMARY KEY'
		LEFT JOIN information_schema.key_column_usage kcu
			ON kcu.constraint_name = tc.constraint_name
			AND kcu.table_schema = tc.table_schema
			AND kcu.column_name = c.column_name
		WHERE c.table_schema = 'public' AND c.table_name = $1
		ORDER BY c.ordinal_position
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

package sqlite

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"

	"OrangeSQL/internal/domain"
)

type Repository struct {
	db *sql.DB
}

func New(path string) (*Repository, error) {
	db, err := sql.Open("sqlite", path) // ★ "sqlite3" ではなく "sqlite"
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) Query(q string) (domain.QueryResult, error) {
	rows, err := r.db.Query(q)
	if err != nil {
		return domain.QueryResult{}, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return domain.QueryResult{}, err
	}

	result := domain.QueryResult{Columns: cols, Rows: make([][]string, 0)}

	raw := make([]interface{}, len(cols))
	ptrs := make([]interface{}, len(cols))
	for i := range raw {
		ptrs[i] = &raw[i]
	}

	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return domain.QueryResult{}, err
		}
		line := make([]string, len(cols))
		for i, v := range raw {
			switch t := v.(type) {
			case []byte:
				line[i] = string(t)
			default:
				line[i] = fmt.Sprint(t)
			}
		}
		result.Rows = append(result.Rows, line)
	}

	if err := rows.Err(); err != nil {
		return domain.QueryResult{}, err
	}

	return result, nil
}

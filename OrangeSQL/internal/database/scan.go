package database

import (
	"database/sql"
	"fmt"
)

// scanRows は *sql.Rows から QueryResult を構築する共通ヘルパー。
// 全ドライバの Query() で再利用する。
func scanRows(rows *sql.Rows) (QueryResult, error) {
	cols, err := rows.Columns()
	if err != nil {
		return QueryResult{}, err
	}

	result := QueryResult{
		Columns: cols,
		Rows:    make([][]*string, 0),
	}

	raw := make([]interface{}, len(cols))
	ptrs := make([]interface{}, len(cols))
	for i := range raw {
		ptrs[i] = &raw[i]
	}

	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return QueryResult{}, err
		}
		row := make([]*string, len(cols))
		for i, v := range raw {
			if v == nil {
				row[i] = nil
				continue
			}
			var s string
			switch t := v.(type) {
			case []byte:
				s = string(t)
			default:
				s = fmt.Sprint(t)
			}
			row[i] = &s
		}
		result.Rows = append(result.Rows, row)
	}

	if err := rows.Err(); err != nil {
		return QueryResult{}, err
	}

	return result, nil
}

// execStatement は *sql.DB で DML/DDL を実行する共通ヘルパー。
func execStatement(db *sql.DB, q string) (ExecResult, error) {
	res, err := db.Exec(q)
	if err != nil {
		return ExecResult{}, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return ExecResult{}, err
	}
	return ExecResult{AffectedRows: affected}, nil
}

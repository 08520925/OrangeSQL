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

// execBatch は複数 SQL をトランザクション内で一括実行する共通ヘルパー。
// いずれかが失敗した場合はロールバックしてエラーを返す。
func execBatch(db *sql.DB, statements []string) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}

	var total int64
	for i, stmt := range statements {
		res, err := tx.Exec(stmt)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("statement %d: %w", i+1, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("statement %d rows affected: %w", i+1, err)
		}
		total += affected
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}
	return total, nil
}

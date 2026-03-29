package query

import (
	"OrangeSQL/internal/database"
)

const maxRows = 10000

// BuildResponse は database.QueryResult をAPIレスポンスに変換する。
func BuildResponse(result database.QueryResult, elapsedMs int64, sql string, db database.Database) Response {
	truncated := false
	rows := result.Rows
	if len(rows) > maxRows {
		rows = rows[:maxRows]
		truncated = true
	}

	// [][]*string → [][]any に変換（nil を JSON null にするため）
	jsonRows := make([][]any, len(rows))
	for i, row := range rows {
		jsonRow := make([]any, len(row))
		for j, v := range row {
			if v == nil {
				jsonRow[j] = nil
			} else {
				jsonRow[j] = *v
			}
		}
		jsonRows[i] = jsonRow
	}

	editableTable, pkColumns := DetectEditableTable(sql, result.Columns, db)

	return Response{
		Columns:       result.Columns,
		Rows:          jsonRows,
		RowCount:      len(jsonRows),
		ExecutionTime: elapsedMs,
		Truncated:     truncated,
		EditableTable: editableTable,
		PKColumns:     pkColumns,
	}
}

package schema

import (
	"OrangeSQL/internal/database"
)

// BuildResponse は database.TableInfo のスライスをAPIレスポンスに変換する。
func BuildResponse(tables []database.TableInfo) Response {
	entries := make([]TableEntry, len(tables))
	for i, t := range tables {
		entries[i] = TableEntry{Name: t.Name, Type: t.Type}
	}
	return Response{Tables: entries}
}

// BuildColumnsResponse は database.ColumnInfo のスライスをAPIレスポンスに変換する。
func BuildColumnsResponse(cols []database.ColumnInfo) ColumnsResponse {
	entries := make([]ColumnEntry, len(cols))
	for i, c := range cols {
		entries[i] = ColumnEntry{Name: c.Name, Type: c.Type, PK: c.PK, NotNull: c.NotNull}
	}
	return ColumnsResponse{Columns: entries}
}

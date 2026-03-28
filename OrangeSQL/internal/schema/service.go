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

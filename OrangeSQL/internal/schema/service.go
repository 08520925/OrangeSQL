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

// BuildCompletionsResponse は全テーブル + カラムを一括取得してレスポンスに変換する。
func BuildCompletionsResponse(db database.Database) (CompletionsResponse, error) {
	tables, err := db.Tables()
	if err != nil {
		return CompletionsResponse{}, err
	}

	result := make([]CompletionTable, 0, len(tables))
	for _, t := range tables {
		cols, err := db.Columns(t.Name)
		if err != nil {
			// カラム取得失敗はスキップ（権限不足等）
			result = append(result, CompletionTable{Name: t.Name, Type: t.Type, Columns: nil})
			continue
		}
		entries := make([]ColumnEntry, len(cols))
		for i, c := range cols {
			entries[i] = ColumnEntry{Name: c.Name, Type: c.Type, PK: c.PK, NotNull: c.NotNull, Comment: c.Comment}
		}
		result = append(result, CompletionTable{Name: t.Name, Type: t.Type, Columns: entries})
	}
	return CompletionsResponse{Tables: result}, nil
}

// BuildColumnsResponse は database.ColumnInfo のスライスをAPIレスポンスに変換する。
func BuildColumnsResponse(cols []database.ColumnInfo) ColumnsResponse {
	entries := make([]ColumnEntry, len(cols))
	for i, c := range cols {
		entries[i] = ColumnEntry{Name: c.Name, Type: c.Type, PK: c.PK, NotNull: c.NotNull, Comment: c.Comment}
	}
	return ColumnsResponse{Columns: entries}
}

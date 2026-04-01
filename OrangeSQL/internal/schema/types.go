package schema

// Response は GET /api/schema/tables の JSON レスポンス。
type Response struct {
	Tables []TableEntry `json:"tables"`
}

// TableEntry はテーブルまたはビューの情報。
type TableEntry struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// ColumnsResponse は GET /api/schema/columns の JSON レスポンス。
type ColumnsResponse struct {
	Columns []ColumnEntry `json:"columns"`
}

// ColumnEntry はカラムの情報。
type ColumnEntry struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	PK      bool   `json:"pk"`
	NotNull bool   `json:"notNull"`
	Comment string `json:"comment"`
}

// CompletionsResponse は GET /api/schema/completions の JSON レスポンス。
type CompletionsResponse struct {
	Tables []CompletionTable `json:"tables"`
}

// CompletionTable はテーブル名 + カラム一覧。
type CompletionTable struct {
	Name    string        `json:"name"`
	Type    string        `json:"type"`
	Columns []ColumnEntry `json:"columns"`
}

// ErrorResponse はエラー時の JSON レスポンス。
type ErrorResponse struct {
	Error string `json:"error"`
}

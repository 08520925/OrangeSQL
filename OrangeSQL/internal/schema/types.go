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

// ErrorResponse はエラー時の JSON レスポンス。
type ErrorResponse struct {
	Error string `json:"error"`
}

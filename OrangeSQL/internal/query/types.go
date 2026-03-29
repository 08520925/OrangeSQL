package query

// Response は POST /api/query の JSON レスポンス。
type Response struct {
	Columns       []string `json:"columns"`
	Rows          [][]any  `json:"rows"`
	RowCount      int      `json:"rowCount"`
	ExecutionTime int64    `json:"executionTimeMs"`
	Truncated     bool     `json:"truncated,omitempty"`
	EditableTable string   `json:"editableTable,omitempty"`
	PKColumns     []string `json:"pkColumns,omitempty"`
}

// ErrorResponse はエラー時の JSON レスポンス。
type ErrorResponse struct {
	Error string `json:"error"`
}

// Request は POST /api/query の JSON リクエスト。
type Request struct {
	SQL string `json:"sql"`
}

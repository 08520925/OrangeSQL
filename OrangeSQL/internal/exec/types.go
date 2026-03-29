package exec

// Response は POST /api/exec の JSON レスポンス。
type Response struct {
	AffectedRows int64 `json:"affectedRows"`
	ExecutionTime int64 `json:"executionTimeMs"`
}

// ErrorResponse はエラー時の JSON レスポンス。
type ErrorResponse struct {
	Error string `json:"error"`
}

// Request は POST /api/exec の JSON リクエスト。
type Request struct {
	SQL string `json:"sql"`
}

// BatchRequest は POST /api/exec/batch の JSON リクエスト。
type BatchRequest struct {
	Statements []string `json:"statements"`
}

// BatchResponse は POST /api/exec/batch の JSON レスポンス。
type BatchResponse struct {
	TotalAffectedRows int64 `json:"totalAffectedRows"`
	ExecutionTime     int64 `json:"executionTimeMs"`
}

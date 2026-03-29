package exec

import (
	"encoding/json"
	"net/http"
	"time"

	"OrangeSQL/internal/database"
)

// BatchHandler は POST /api/exec/batch のHTTPハンドラを返す。
// 複数 SQL をトランザクション内で一括実行する。
func BatchHandler(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req BatchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid JSON"})
			return
		}

		if len(req.Statements) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "statements field is required"})
			return
		}

		start := time.Now()
		total, err := db.ExecBatch(req.Statements)
		elapsed := time.Since(start).Milliseconds()

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(BatchResponse{
			TotalAffectedRows: total,
			ExecutionTime:     elapsed,
		})
	}
}

package query

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"OrangeSQL/internal/database"
)

// Handler は POST /api/query のHTTPハンドラを返す。
func Handler(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid JSON"})
			return
		}

		if strings.TrimSpace(req.SQL) == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "sql field is required"})
			return
		}

		start := time.Now()
		result, err := db.Query(req.SQL)
		elapsed := time.Since(start).Milliseconds()

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		resp := BuildResponse(result, elapsed)
		json.NewEncoder(w).Encode(resp)
	}
}

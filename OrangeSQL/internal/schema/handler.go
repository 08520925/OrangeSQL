package schema

import (
	"encoding/json"
	"net/http"

	"OrangeSQL/internal/database"
)

// Handler は GET /api/schema/tables のHTTPハンドラを返す。
func Handler(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		tables, err := db.Tables()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		resp := BuildResponse(tables)
		json.NewEncoder(w).Encode(resp)
	}
}

package schema

import (
	"encoding/json"
	"net/http"

	"OrangeSQL/internal/database"
)

// TablesHandler は GET /api/schema/tables のHTTPハンドラを返す。
func TablesHandler(db database.Database) http.HandlerFunc {
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

// CompletionsHandler は GET /api/schema/completions のHTTPハンドラを返す。
// 全テーブル + カラムを一括返却する。
func CompletionsHandler(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		resp, err := BuildCompletionsResponse(db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(resp)
	}
}

// ColumnsHandler は GET /api/schema/columns のHTTPハンドラを返す。
func ColumnsHandler(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		tableName := r.URL.Query().Get("table")
		if tableName == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "table parameter is required"})
			return
		}

		cols, err := db.Columns(tableName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		resp := BuildColumnsResponse(cols)
		json.NewEncoder(w).Encode(resp)
	}
}

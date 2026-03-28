package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"OrangeSQL/internal/database"
	"OrangeSQL/internal/exec"
	"OrangeSQL/internal/query"
	"OrangeSQL/internal/schema"
)

// NewRouter は API ルーティングを設定した http.Handler を返す。
func NewRouter(db database.Database) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/query", query.Handler(db))
	mux.HandleFunc("POST /api/exec", exec.Handler(db))
	mux.HandleFunc("GET /api/schema/tables", schema.Handler(db))

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		status := "ok"
		if _, err := db.Query("SELECT 1"); err != nil {
			status = "error"
		}
		json.NewEncoder(w).Encode(map[string]string{
			"status":   status,
			"database": db.Name(),
		})
	})

	return withCORS(withLogging(mux))
}

// withCORS は開発時の CORS ヘッダーを付与するミドルウェア。
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:5173" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// withLogging はリクエストログを出力するミドルウェア。
func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

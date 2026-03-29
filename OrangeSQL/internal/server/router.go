package server

import (
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	"OrangeSQL/internal/exec"
	"OrangeSQL/internal/profile"
	"OrangeSQL/internal/query"
	"OrangeSQL/internal/schema"
)

// NewRouter は API ルーティングを設定した http.Handler を返す。
// 全ハンドラは cm.DB() 経由で現在の DB 接続にアクセスする。
// staticFS が nil でなければ、API 以外のリクエストに埋め込みファイルを返す。
func NewRouter(cm *profile.ConnectionManager, staticFS fs.FS) http.Handler {
	mux := http.NewServeMux()

	// DB を使うハンドラはラムダで cm.DB() を渡す
	mux.HandleFunc("POST /api/query", func(w http.ResponseWriter, r *http.Request) {
		query.Handler(cm.DB())(w, r)
	})
	mux.HandleFunc("POST /api/exec", func(w http.ResponseWriter, r *http.Request) {
		exec.Handler(cm.DB())(w, r)
	})
	mux.HandleFunc("POST /api/exec/batch", func(w http.ResponseWriter, r *http.Request) {
		exec.BatchHandler(cm.DB())(w, r)
	})
	mux.HandleFunc("GET /api/schema/tables", func(w http.ResponseWriter, r *http.Request) {
		schema.TablesHandler(cm.DB())(w, r)
	})
	mux.HandleFunc("GET /api/schema/columns", func(w http.ResponseWriter, r *http.Request) {
		schema.ColumnsHandler(cm.DB())(w, r)
	})
	mux.HandleFunc("GET /api/schema/completions", func(w http.ResponseWriter, r *http.Request) {
		schema.CompletionsHandler(cm.DB())(w, r)
	})

	// プロファイル管理 API
	mux.HandleFunc("GET /api/profiles", profile.ListHandler(cm))
	mux.HandleFunc("POST /api/profiles", profile.CreateHandler(cm))
	mux.HandleFunc("PUT /api/profiles/{id}", profile.UpdateHandler(cm))
	mux.HandleFunc("DELETE /api/profiles/{id}", profile.DeleteHandler(cm))
	mux.HandleFunc("POST /api/profiles/{id}/connect", profile.ConnectHandler(cm))

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		db := cm.DB()
		status := "ok"
		if _, err := db.Query("SELECT 1"); err != nil {
			status = "error"
		}
		json.NewEncoder(w).Encode(map[string]string{
			"status":   status,
			"database": db.Name(),
		})
	})

	// 静的ファイル配信（埋め込みフロントエンド）
	if staticFS != nil {
		fileServer := http.FileServerFS(staticFS)
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// API リクエストはここに来ない（上で登録済み）
			// 存在しないパスは index.html にフォールバック（SPA 対応）
			path := r.URL.Path
			if path != "/" && !strings.HasPrefix(path, "/api/") {
				// ファイルが存在するか確認
				f, err := staticFS.Open(strings.TrimPrefix(path, "/"))
				if err != nil {
					// ファイルが見つからなければ index.html を返す
					r.URL.Path = "/"
				} else {
					f.Close()
				}
			}
			fileServer.ServeHTTP(w, r)
		})
	}

	return withCORS(withLogging(mux))
}

// withCORS は開発時の CORS ヘッダーを付与するミドルウェア。
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:5173" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
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

# Phase 7: シングルバイナリ配布 — 仕様書

## 概要

`go:embed` でフロントエンドのビルド成果物（`frontend/dist`）を Go バイナリに埋め込み、
exe 1つで配布・起動できるようにする。

---

## 1. 仕組み

### 1.1 ビルドフロー

```
pnpm build → frontend/dist/ に成果物生成
go build   → dist/ を embed して exe に埋め込み
```

### 1.2 実行時の動作

1. exe を起動
2. Go サーバーが `:5522` で起動
3. API（`/api/*`）はバックエンドが処理
4. それ以外のリクエスト（`/`, `/assets/*` 等）は埋め込みファイルを返す
5. ブラウザが自動で開く（`--no-browser` フラグで抑制可能）

### 1.3 開発モードとの両立

- 開発時: `pnpm dev`（Vite dev server、ポート5173）+ `go run`（ポート5522）
- 本番時: exe 単体で起動（ポート5522 で API + 静的ファイル両方）
- 開発時は Vite proxy が `/api` を 5522 に転送するため、変更不要

---

## 2. バックエンド変更

### 2.1 embed ファイル

```go
// internal/server/embed.go
package server

import "embed"

//go:embed all:dist
var distFS embed.FS
```

`frontend/dist` を `internal/server/dist` にコピー（ビルドスクリプトで自動化）するか、
または `go:embed` のパスを調整する。

→ **シンプルに `go:embed` のパスを main.go 側に置き、NewRouter に渡す方式を採用。**

```go
// main.go
//go:embed frontend/dist
var frontendFS embed.FS
```

### 2.2 静的ファイルサーバー

`NewRouter` に `embed.FS` を渡し、API 以外のリクエストに埋め込みファイルを返す。

```go
func NewRouter(cm *profile.ConnectionManager, staticFS fs.FS) http.Handler {
    // ... API ルート登録 ...

    // 静的ファイル: API 以外は埋め込みファイルを返す
    // SPA のため、存在しないパスは index.html にフォールバック
    mux.Handle("/", http.FileServerFS(staticFS))
}
```

### 2.3 ブラウザ自動オープン

起動時にデフォルトブラウザで `http://localhost:5522` を開く。

```go
// Windows
exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
```

`--no-browser` フラグで抑制。

---

## 3. ビルドスクリプト

`build.sh`（bash）を作成:

```bash
#!/bin/bash
cd frontend && pnpm build && cd ..
go build -o orangesql.exe .
```

---

## 4. Issue 一覧（Phase 7）

| # | タイトル | スコープ | 依存 |
|---|---------|---------|------|
| 44 | go:embed で frontend/dist を埋め込み + 静的ファイル配信 | main.go, router.go | なし |
| 45 | ブラウザ自動オープン | main.go | #44 |
| 46 | ビルドスクリプト + テスト | build.sh | #44 |

---

## 5. Phase 7 完了条件

- [ ] `pnpm build && go build` で exe が生成される
- [ ] exe 単体で起動 → ブラウザが開く → OrangeSQL が動作する
- [ ] `--no-browser` で自動オープンを抑制できる
- [ ] 開発モード（pnpm dev + go run）も引き続き動作する
- [ ] `go test ./internal/...` で全テストがパスする

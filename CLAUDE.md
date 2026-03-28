# OrangeSQL

軽量デスクトップSQLエディタ。Go Webサーバー + Vue 3 フロントエンド。

## 技術スタック

- Backend: Go 1.26.1, `net/http`, `modernc.org/sqlite` (CGO=0)
- Frontend: Vue 3.5, TypeScript 6.0, Vite 8.0, CodeMirror 6, pnpm
- API: REST (JSON), ポート 5522
- テスト: `go test ./internal/...`（httptest + インメモリSQLite、モック不使用）

## アーキテクチャ

バーティカルスライス構成。機能単位で handler/service/types をコロケーション。

- Backend: `internal/{query,exec,schema,database,server}/`
- Frontend: `src/features/{sql-editor,results-panel,schema-sidebar}/` + `src/shared/`
- 詳細: `docs/spec/phase1-foundation.md`

## コーディングルール

- 1ファイル300行以下
- ファイル名で中身がわかる命名（`utils.ts` 等の曖昧な名前を避ける）
- Go: 標準ライブラリ優先、エラーは全て JSON レスポンスで返す
- TypeScript: 最厳格設定（strict + noUncheckedIndexedAccess + exactOptionalPropertyTypes 等）
- TypeScript: `any` 禁止、`as` キャスト禁止、`unknown` + 型ガードを使う
- Vue SFC は composable を積極的に抽出
- CodeMirror 6 は Vue ラッパー禁止、素で `new EditorView()` を使う

## API エンドポイント

- `POST /api/query` — SELECT 系
- `POST /api/exec` — INSERT/UPDATE/DELETE 系
- `GET /api/schema/tables` — テーブル一覧
- `GET /api/health` — ヘルスチェック

## コマンド

```bash
# バックエンド ビルド+起動（Windows では go run は使わない）
cd OrangeSQL && go build -o orangesql.exe . && ./orangesql.exe

# テスト
cd OrangeSQL && go test ./internal/...

# フロントエンド開発
cd OrangeSQL/frontend && pnpm dev

# フロントエンドビルド
cd OrangeSQL/frontend && pnpm build

# E2E テスト（バックエンド+フロントエンド起動中に実行）
cd OrangeSQL/frontend && pnpm test:e2e

# パッケージ追加
cd OrangeSQL/frontend && pnpm add <package>
```

## 回答言語

日本語で回答すること。

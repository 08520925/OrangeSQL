# OrangeSQL

軽量デスクトップ SQL エディタ。Go Web サーバー + Vue 3 フロントエンド。
exe 1つで配布可能（go:embed でフロントエンドを埋め込み）。

## 技術スタック

- Backend: Go 1.26, `net/http`, CGO=0
- DB ドライバ: `modernc.org/sqlite`, `pgx/v5`, `go-sql-driver/mysql`, `go-mssqldb`
- Frontend: Vue 3.5, TypeScript 6.0, Vite 8.0, CodeMirror 6, pnpm
- API: REST (JSON), ポート 5522
- テスト: Go (`httptest` + インメモリ SQLite)、E2E (Playwright)
- CI/CD: GitHub Actions（タグ push で 4 OS 向けバイナリを自動リリース）

## プロジェクト構成

```
OrangeSQL/
├── main.go                      # エントリポイント（embed + サーバー起動）
├── build.sh                     # ビルドスクリプト
├── docker-compose.test.yml      # 統合テスト用 DB コンテナ
├── internal/                    # バックエンド（→ OrangeSQL/internal/CLAUDE.md）
│   ├── database/                # DB 抽象化（Interface + 4 ドライバ）
│   ├── query/                   # POST /api/query + 編集可能テーブル判定
│   ├── exec/                    # POST /api/exec, /api/exec/batch
│   ├── schema/                  # GET /api/schema/*
│   ├── profile/                 # 接続プロファイル管理
│   ├── server/                  # ルーター + ミドルウェア
│   └── testutil/                # テスト用ヘルパー
├── frontend/                    # フロントエンド（→ OrangeSQL/frontend/CLAUDE.md）
│   ├── src/features/            # 機能別コンポーネント
│   ├── src/shared/              # 共通 API クライアント + 型定義
│   └── e2e/                     # Playwright E2E テスト
└── docs/
    ├── roadmap.md               # ロードマップ
    └── spec/                    # Phase 別仕様書
```

## コマンド

```bash
# プロダクションビルド（exe 生成）
cd OrangeSQL && bash build.sh

# バックエンド起動（開発用）
cd OrangeSQL && go build -o orangesql.exe . && ./orangesql.exe --no-browser

# フロントエンド開発サーバー（ポート 5173、API を 5522 にプロキシ）
cd OrangeSQL/frontend && pnpm dev

# バックエンドテスト
cd OrangeSQL && go test ./internal/...

# E2E テスト（バックエンド + フロントエンド起動中に実行）
cd OrangeSQL/frontend && pnpm test:e2e

# Docker 統合テスト（実 DB 接続）
cd OrangeSQL && docker compose -f docker-compose.test.yml up -d
go test ./internal/database/ -tags=integration -v
docker compose -f docker-compose.test.yml down

# フロントエンドビルド（単体）
cd OrangeSQL/frontend && pnpm build

# パッケージ追加
cd OrangeSQL/frontend && pnpm add <package>
```

## API エンドポイント

| メソッド | パス | 用途 |
|---------|------|------|
| POST | `/api/query` | SELECT 系（editableTable, pkColumns 付き） |
| POST | `/api/exec` | INSERT/UPDATE/DELETE/DDL |
| POST | `/api/exec/batch` | 複数 SQL をトランザクション一括実行 |
| GET | `/api/schema/tables` | テーブル一覧 |
| GET | `/api/schema/columns?table=X` | カラム情報 |
| GET | `/api/schema/completions` | 全テーブル+カラム一括（オートコンプリート用） |
| GET/POST/PUT/DELETE | `/api/profiles[/{id}]` | 接続プロファイル CRUD |
| POST | `/api/profiles/{id}/connect` | 接続切り替え |
| GET | `/api/health` | ヘルスチェック |

## コーディングルール

- 1ファイル 300 行以下
- ファイル名で中身がわかる命名（`utils.ts` 等の曖昧な名前を避ける）
- バーティカルスライス構成: 機能単位で handler/service/types をコロケーション
- 仕様駆動: 仕様書 → GitHub Issue → 実装 → テスト の順で進める
- 日本語で回答すること

## 開発フロー

1. `docs/spec/` に仕様書を作成
2. GitHub Issue を作成（`phase<N>` ラベル付き）
3. 実装 + テスト
4. Issue を close、ロードマップ更新
5. コミット + プッシュ
6. リリース時: `git tag v<X.Y.Z> && git push origin v<X.Y.Z>` で GitHub Actions が自動ビルド

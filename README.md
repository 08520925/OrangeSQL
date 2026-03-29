# OrangeSQL

軽量デスクトップ SQL エディタ。exe 1つで起動し、ブラウザ上で SQL を書いて実行できます。

## 特徴

- **マルチDB対応** — SQLite / PostgreSQL / MySQL / SQL Server
- **SQL オートコンプリート** — テーブル名・カラム名・キーワード補完
- **テーブルデータ編集** — 結果グリッド上でインライン編集、行追加・削除（トランザクション保証）
- **複数タブ** — タブごとに独立した SQL とクエリ結果
- **接続プロファイル管理** — 複数 DB を登録して切り替え
- **結果エクスポート** — CSV / JSON ダウンロード、CSV / Markdown クリップボードコピー
- **シングルバイナリ** — Go / Node.js のインストール不要、exe 1つで動作

## ダウンロード

[Releases](https://github.com/08520925/OrangeSQL/releases) から OS に合ったバイナリをダウンロードしてください。

| ファイル | 対象 |
|---------|------|
| `orangesql-windows-amd64.exe` | Windows (64bit) |
| `orangesql-darwin-amd64` | macOS (Intel) |
| `orangesql-darwin-arm64` | macOS (Apple Silicon) |
| `orangesql-linux-amd64` | Linux (64bit) |

## 使い方

ダウンロードした exe を実行するだけです。ブラウザが自動で開きます。

```bash
# Windows: ダブルクリック or
./orangesql-windows-amd64.exe

# macOS / Linux
chmod +x orangesql-*
./orangesql-linux-amd64
```

### オプション

```
--db <path>       初回起動時の SQLite ファイルパス（デフォルト: ./data.db）
--no-browser      ブラウザの自動オープンを無効化
```

### 接続先 DB の設定

起動後、ヘッダーのドロップダウンから「+ 新規接続」でプロファイルを作成します。

- **SQLite**: ファイルパスを指定
- **PostgreSQL / MySQL / SQL Server**: ホスト・ポート・ユーザー・パスワード・データベース名を入力

## 技術スタック

- **Backend**: Go, `net/http`, `modernc.org/sqlite`, `pgx`, `go-sql-driver/mysql`, `go-mssqldb`
- **Frontend**: Vue 3, TypeScript, Vite, CodeMirror 6
- **配布**: `go:embed` でフロントエンドを exe に埋め込み

## 開発

### 必要環境

- Go 1.26+
- Node.js 22+
- pnpm

### セットアップ

```bash
cd OrangeSQL/frontend
pnpm install
```

### 開発モード

```bash
# バックエンド（ポート 5522）
cd OrangeSQL && go build -o orangesql.exe . && ./orangesql.exe --no-browser

# フロントエンド（ポート 5173、API は 5522 にプロキシ）
cd OrangeSQL/frontend && pnpm dev
```

ブラウザで http://localhost:5173 を開いてください。

### プロダクションビルド

```bash
cd OrangeSQL && bash build.sh
```

### テスト

```bash
# バックエンド
cd OrangeSQL && go test ./internal/...

# E2E（バックエンド + フロントエンド起動中に実行）
cd OrangeSQL/frontend && pnpm test:e2e

# Docker 統合テスト（PostgreSQL / MySQL / SQL Server）
cd OrangeSQL && docker compose -f docker-compose.test.yml up -d
go test ./internal/database/ -tags=integration -v
docker compose -f docker-compose.test.yml down
```

## ライセンス

MIT

# Phase 1: 基盤 — 仕様書

## 概要

OrangeSQL の最小動作基盤を構築する。
Go Webサーバー + Vue 3 フロントエンドで、SQLiteに接続してSQLを実行し、結果をテーブル表示する。

---

## 1. 画面構成

### 1.1 レイアウト

```
┌─────────────────────────────────────────────────┐
│  ヘッダーバー（アプリ名・接続状態）                   │
├────────────┬────────────────────────────────────┤
│            │  SQLエディタ（Monaco Editor）         │
│  スキーマ   │                                     │
│  サイドバー  │                                     │
│            ├────────────────────────────────────┤
│ (テーブル   │  結果パネル                           │
│  一覧)     │  ┌─────┬─────┬─────┐               │
│            │  │ col1│ col2│ col3│               │
│            │  ├─────┼─────┼─────┤               │
│            │  │  …  │  …  │  …  │               │
│            │  └─────┴─────┴─────┘               │
│            │  ステータスバー（実行時間・行数）        │
└────────────┴────────────────────────────────────┘
```

### 1.2 各エリアの役割

| エリア | 役割 | Phase 1 スコープ |
|--------|------|-----------------|
| ヘッダーバー | アプリ名表示、接続中のDB表示 | アプリ名 + SQLiteファイルパス表示 |
| スキーマサイドバー | テーブル一覧・カラム情報 | テーブル名一覧のみ（クリックで `SELECT * FROM table LIMIT 100` をエディタに挿入） |
| SQLエディタ | SQL入力・編集 | Monaco Editor、SQL構文ハイライト、Ctrl+Enter で実行 |
| 結果パネル | クエリ結果のテーブル表示 | カラムヘッダ固定、スクロール対応 |
| ステータスバー | 実行結果のサマリ | 実行時間、取得行数、エラーメッセージ |

### 1.3 操作フロー

1. アプリ起動 → ブラウザで `http://localhost:5522` を開く
2. SQLite ファイルに自動接続（デフォルト: `./data.db`、コマンドライン引数 `-db` で変更可能）
3. エディタに SQL を入力
4. Ctrl+Enter または実行ボタンで実行
5. SELECT → 結果パネルにテーブル表示
6. INSERT/UPDATE/DELETE/DDL → ステータスバーに affected rows 表示
7. エラー → ステータスバーにエラーメッセージ表示

### 1.4 SQLite ファイルパスの指定

```bash
# デフォルト（カレントディレクトリに data.db を作成/接続）
go run .

# 指定したファイルに接続
go run . -db /path/to/mydb.sqlite
```

- コマンドライン引数 `-db` でファイルパスを指定できる
- 省略時はカレントディレクトリの `data.db` を使用
- ファイルが存在しない場合は **SQLite が自動作成する**（空のDBとして起動）
- ヘッダーバーに接続中のファイルパスを表示する

### 1.5 query / exec の振り分けロジック

フロントエンド側で SQL の先頭キーワードを見て API を振り分ける。

| 先頭キーワード | API |
|---------------|-----|
| `SELECT`, `SHOW`, `EXPLAIN`, `PRAGMA`, `WITH` | `POST /api/query` |
| それ以外（`INSERT`, `UPDATE`, `DELETE`, `CREATE`, `DROP`, `ALTER` 等） | `POST /api/exec` |

- 判定はケースインセンシティブ、先頭の空白・改行は無視
- 複数SQL文（セミコロン区切り）は **先頭の1文だけ実行する**（Phase 1 の制約）

### 1.6 大量結果の制限

- サーバー側で結果セットの行数に上限を設ける: **最大 10,000 行**
- 上限を超えた場合、10,000 行で切り詰めてレスポンスに `truncated: true` を含める
- フロントは `truncated` が true のとき「結果が切り詰められました」と表示する

### 1.7 ローディング状態

- SQL実行中は実行ボタンを無効化し、スピナーを表示する
- 実行中に再度 Ctrl+Enter を押しても無視する

### 1.8 サイドバーの更新

- 初回表示時にテーブル一覧を取得
- exec API の実行成功後、自動でテーブル一覧を再取得する（CREATE/DROP に対応）
- 手動リフレッシュボタンも設ける

### 1.9 レイアウト分割

- サイドバー幅: **250px 固定**（Phase 1）。リサイズは Phase 2 以降
- エディタ / 結果パネルの上下分割: **50:50 固定**（Phase 1）。ドラッグリサイズは Phase 2 以降
- ヘッダーバー高さ: **40px 固定**
- ステータスバー高さ: **28px 固定**

### 1.10 テーマ

- Phase 1 ではダークテーマ固定（Monaco: `vs-dark`）
- ライトテーマ対応は Phase 2 以降

---

## 2. API 仕様

### 2.1 基本方針

- REST API（JSON）
- ベースURL: `http://localhost:5522/api`
- 自分専用ツールのため認証なし
- Go 標準 `net/http` を使用
- Content-Type: `application/json`（リクエスト・レスポンスとも）

### 2.2 HTTPステータスコード

| ステータス | 用途 |
|-----------|------|
| 200 | 正常完了 |
| 400 | SQLエラー、リクエスト不正（JSON パース失敗、`sql` フィールド欠落・空文字列） |
| 405 | 許可されていない HTTP メソッド |
| 500 | サーバー内部エラー（DB接続断等） |

### 2.3 リクエストバリデーション（query / exec 共通）

以下の場合は HTTP 400 を返す:
- リクエストボディが空、または不正な JSON
- `sql` フィールドが存在しない、または空文字列

```json
{
  "error": "sql field is required"
}
```

### 2.4 エンドポイント一覧

#### `POST /api/query`

SELECT 系の SQL を実行し、結果セットを返す。

**Request:**
```json
{
  "sql": "SELECT * FROM users LIMIT 10"
}
```

**Response（成功 / HTTP 200）:**
```json
{
  "columns": ["id", "name", "email"],
  "rows": [
    ["1", "Alice", "alice@example.com"],
    ["2", "Bob", null]
  ],
  "rowCount": 2,
  "executionTimeMs": 12
}
```

**Response（切り詰め / HTTP 200）:**

10,000 行を超えた場合、10,000 行で切り詰めて `truncated: true` を付与する。
```json
{
  "columns": ["id", "name"],
  "rows": [
    ["1", "Alice"],
    ["2", "Bob"],
    "..."
  ],
  "rowCount": 10000,
  "executionTimeMs": 230,
  "truncated": true
}
```
※ `rows` の実際の要素数は 10,000。上記は省略表記。

**Response（エラー / HTTP 400）:**
```json
{
  "error": "near \"SELEC\": syntax error"
}
```

#### `POST /api/exec`

INSERT / UPDATE / DELETE / DDL（CREATE TABLE, DROP TABLE, ALTER TABLE 等）を実行する。
DDL の場合、`affectedRows` は 0 を返す。

**Request:**
```json
{
  "sql": "INSERT INTO users (name, email) VALUES ('Alice', 'alice@example.com')"
}
```

**Response（成功 / HTTP 200）:**
```json
{
  "affectedRows": 1,
  "executionTimeMs": 5
}
```

**Response（エラー / HTTP 400）:**
```json
{
  "error": "no such table: user"
}
```

#### `GET /api/schema/tables`

テーブル一覧を取得する。

**Response（成功 / HTTP 200）:**
```json
{
  "tables": [
    { "name": "users", "type": "table" },
    { "name": "orders", "type": "table" },
    { "name": "user_view", "type": "view" }
  ]
}
```

**Response（エラー / HTTP 500）:**
```json
{
  "error": "database connection lost"
}
```

#### `GET /api/health`

サーバーの生存確認。

**Response（正常 / HTTP 200）:**
```json
{
  "status": "ok",
  "database": "data.db"
}
```

**Response（異常 / HTTP 500）:**
```json
{
  "status": "error",
  "database": "data.db",
  "error": "database is locked"
}
```

---

## 3. データフォーマット

### 3.1 クエリ結果の型

```typescript
// フロントエンド側の型定義

// POST /api/query のレスポンス
type QueryResult = {
  columns: string[];
  rows: (string | null)[][];  // NULLはnull
  rowCount: number;
  executionTimeMs: number;
  truncated?: boolean;         // 10,000行超で切り詰めた場合 true
};

// POST /api/exec のレスポンス
type ExecResult = {
  affectedRows: number;
  executionTimeMs: number;
};

// エラー時（両APIで共通）
type ApiError = {
  error: string;
};
```

### 3.2 rows の値は全て string、NULL は null

- DB の型に関わらず、結果は全て文字列に変換して返す
- NULL は JSON の `null` として返す
- フロント側では `null` を `"NULL"` と表示し、グレーアウト+イタリックで視覚的に区別する
- BLOB（バイナリ）は `"[BLOB (N bytes)]"` という文字列に変換して返す
- 型表示が必要になったら Phase 2 以降で拡張

### 3.3 Database interface と共通型（Go 側）

循環依存を避けるため、Database interface と それが返す型は全て `database` パッケージ内に定義する。
各スライス（query / exec / schema）はこの interface と型を import して使う。

```go
// internal/database/interface.go

type Database interface {
    // SELECT 系。columns と rows を返す。
    Query(sql string) (QueryResult, error)
    // INSERT/UPDATE/DELETE/DDL 系。affected rows を返す。
    Exec(sql string) (ExecResult, error)
    // テーブル・ビュー一覧を返す。
    Tables() ([]TableInfo, error)
    // DB 名（ファイル名等）を返す。
    Name() string
    // 接続を閉じる。
    Close() error
}

type QueryResult struct {
    Columns  []string
    Rows     [][]*string  // nil = SQL NULL
    RowCount int
}

type ExecResult struct {
    AffectedRows int64
}

type TableInfo struct {
    Name string
    Type string  // "table" or "view"
}
```

各スライスの `types.go` は HTTP レスポンス用の JSON 構造体を定義し、`database.QueryResult` → JSON レスポンスへの変換を担当する。

---

## 4. 技術スタック（Phase 1）

| レイヤー | 技術 | バージョン | 備考 |
|---------|------|-----------|------|
| Go | Go `net/http` | 1.26.1 | ルーティングは Go 1.22+ のパターン |
| JSON 処理 | Go `encoding/json` | (Go標準) | |
| DB ドライバ | `modernc.org/sqlite` | 最新 | CGO=0 |
| フロント | Vue 3 | 3.5.31 | |
| 型システム | TypeScript | 6.0.2 | 最厳格設定（下記 4.1 参照） |
| パッケージマネージャ | pnpm | 最新 | npm/yarn は使わない |
| ビルドツール | Vite | 8.0.3 | |
| エディタ | Monaco Editor | 0.55.1 | |
| API 通信 | `fetch` API | (ブラウザ標準) | 外部ライブラリ不要 |

### 4.1 TypeScript 厳格設定

最も厳しい型チェックを適用する。`any` の使用は原則禁止。

```jsonc
// tsconfig.json の compilerOptions（抜粋）
{
  "strict": true,                           // 基本の strict 一式を有効化
  "noUncheckedIndexedAccess": true,         // 配列・辞書アクセスに undefined を付与
  "noUnusedLocals": true,                   // 未使用ローカル変数をエラー
  "noUnusedParameters": true,               // 未使用パラメータをエラー
  "exactOptionalPropertyTypes": true,       // optional と undefined を厳密に区別
  "noImplicitReturns": true,                // 全パスで return を強制
  "noFallthroughCasesInSwitch": true,       // switch の fall-through を禁止
  "noPropertyAccessFromIndexSignature": true, // インデックスシグネチャは [] でアクセス強制
  "verbatimModuleSyntax": true,             // import/export の型修飾を強制
  "skipLibCheck": false                     // node_modules の型もチェック
}
```

#### 型ルール

- `any` の使用は禁止。`unknown` を使い、型ガードで絞り込む
- API レスポンスは必ず型定義を通してから使う（`as` によるキャストは禁止）
- `null` と `undefined` を明確に区別する（API の NULL は `null`、未設定は `undefined`）
- Vue の `ref()` には必ず型引数を付ける（`ref<string>('')`）

---

## 5. ディレクトリ構成（Phase 1 完了後の想定）

AIフレンドリーなバーティカルスライス構成を採用する。
関連するファイルを機能単位で同じディレクトリに置き、AI が少ないファイル読み込みで変更を完結できるようにする。

### 設計原則

- 1ファイル300行以下（AIのコンテキスト効率を最大化）
- ファイル名で中身がわかる命名（`utils.ts` 等の曖昧な名前を避ける）
- 機能単位のコロケーション（関連するhandler/service/typesを同じディレクトリに）
- `shared/` は複数機能から使われるものだけ

```
OrangeSQL/
├── main.go                        # HTTPサーバー起動 + 静的ファイル配信
├── internal/
│   ├── query/                     # SQL実行（SELECT系）
│   │   ├── handler.go             #   POST /api/query のHTTPハンドラ
│   │   ├── handler_test.go        #   インテグレーションテスト
│   │   ├── service.go             #   クエリ実行ロジック・時間計測
│   │   └── types.go               #   QueryResult 等の型定義
│   ├── exec/                      # SQL実行（DML系）
│   │   ├── handler.go             #   POST /api/exec のHTTPハンドラ
│   │   ├── handler_test.go        #   インテグレーションテスト
│   │   ├── service.go             #   DML実行ロジック・時間計測
│   │   └── types.go               #   ExecResult 等の型定義
│   ├── schema/                    # スキーマ情報取得
│   │   ├── handler.go             #   GET /api/schema/tables のHTTPハンドラ
│   │   ├── handler_test.go        #   インテグレーションテスト
│   │   ├── service.go             #   テーブル一覧取得ロジック
│   │   └── types.go               #   TableInfo 等の型定義
│   ├── database/                  # DB抽象・接続管理（横断的基盤）
│   │   ├── interface.go           #   Database interface 定義
│   │   ├── sqlite.go              #   SQLite 実装
│   │   └── sqlite_test.go         #   DB接続・基本操作テスト
│   ├── server/                    # HTTPサーバー設定
│   │   └── router.go              #   ルーティング定義・CORS・ミドルウェア
│   └── testutil/                  # テスト共通ヘルパー
│       └── testdb.go              #   インメモリSQLiteセットアップ
├── frontend/
│   ├── package.json               # 依存定義
│   ├── pnpm-lock.yaml             # pnpm lockfile
│   └── src/
│       ├── features/
│       │   ├── sql-editor/        # SQLエディタ機能
│       │   │   ├── SqlEditor.vue  #   Monaco Editor ラップコンポーネント
│       │   │   ├── use-sql-editor.ts  #   エディタ操作のcomposable
│       │   │   └── types.ts       #   エディタ関連の型
│       │   ├── results-panel/     # 結果表示機能
│       │   │   ├── ResultsPanel.vue   #   テーブル表示コンポーネント
│       │   │   ├── use-results.ts     #   結果データ管理のcomposable
│       │   │   └── types.ts       #   結果関連の型
│       │   └── schema-sidebar/    # スキーマブラウザ機能
│       │       ├── SchemaSidebar.vue  #   テーブル一覧コンポーネント
│       │       ├── use-schema.ts      #   スキーマ取得のcomposable
│       │       └── types.ts       #   スキーマ関連の型
│       ├── shared/
│       │   ├── api.ts             # fetch wrapper（全API共通）
│       │   └── types.ts           # ApiError 等の共通型
│       ├── App.vue                # ルートコンポーネント（レイアウト）
│       └── main.ts                # エントリーポイント
├── docs/
│   └── spec/                      # 仕様書
├── go.mod
└── go.sum
```

---

## 6. 変更対象（現状 → Phase 1）

### 削除するもの

| ファイル/ディレクトリ | 理由 |
|---------------------|------|
| `app.go` | Wails Bind専用 → handler に移行 |
| `wails.json` | Wails設定ファイル、不要に |
| `internal/app/` | レイヤード構成 → バーティカルスライスに移行 |
| `internal/domain/` | 同上、各スライスの types.go に分散 |
| `internal/infra/` | database/ に統合 |
| `frontend/src/wailsjs/` | Wails IPC → fetch API に移行 |
| `frontend/src/app/` | FSD構成 → features構成に移行 |
| `frontend/src/pages/` | 同上、App.vue にレイアウト統合 |
| `frontend/src/widgets/` | 同上、features/ に移行 |
| `frontend/src/entities/` | 同上 |

### 新規作成するもの

| ファイル | 内容 |
|---------|------|
| `internal/query/` | SELECT系の handler + service + types |
| `internal/exec/` | DML系の handler + service + types |
| `internal/schema/` | スキーマ取得の handler + service + types |
| `internal/database/` | Database interface + SQLite 実装 |
| `internal/server/router.go` | ルーティング・CORS設定 |
| `frontend/src/features/` | 各機能の Vue + composable + types |
| `frontend/src/shared/api.ts` | fetch wrapper |

### 改修するもの

| ファイル | 変更内容 |
|---------|---------|
| `main.go` | Wails初期化 → `net/http` サーバー起動 |
| `frontend/vite.config.ts` | Wailsプラグイン削除、APIプロキシ設定追加 |
| `frontend/src/App.vue` | レイアウト構成（ヘッダー+サイドバー+エディタ+結果） |
| `frontend/src/main.ts` | FSDエントリー → シンプルなVueマウント |

---

## 7. テスト戦略（Phase 1）

### 方針

- **バックエンド: インテグレーションテストを主軸にする**
- フロントエンド: Phase 1 ではテストなし（Phase 2 以降で導入）
- モックは使わない（インメモリSQLiteに実際にSQLを実行する）

### 理由

- service層が薄い（SQLをDBに渡すだけ）ため、ユニットテスト単独では価値が低い
- `httptest` + インメモリSQLite で handler → service → DB を一気通貫で検証できる
- 自分専用ツールなのでテストカバレッジ目標は設けず、主要パスの動作確認を優先

### テスト構成

| テストファイル | テスト内容 |
|--------------|-----------|
| `query/handler_test.go` | SELECT 成功、SQLエラー、空結果、NULL含む結果 |
| `exec/handler_test.go` | INSERT/UPDATE/DELETE 成功、SQLエラー |
| `schema/handler_test.go` | テーブル一覧取得、空DB |
| `database/sqlite_test.go` | SQLite接続、Query/Exec基本動作 |

### テスト実行コマンド

```bash
go test ./internal/...
```

### テスト用ヘルパー（`internal/testutil/testdb.go`）

```go
// テスト用のインメモリSQLiteを作成し、テストデータを投入する
func NewTestDB(t *testing.T) database.Database {
    db := database.NewSQLite(":memory:")
    // users テーブル作成 + サンプルデータ投入
    return db
}
```

---

## 8. 非機能要件（Phase 1）

| 項目 | 方針 |
|------|------|
| 認証 | なし（localhost のみバインド、自分専用） |
| CORS | 開発時のみ許可（下記 8.1 参照） |
| エラーハンドリング | Go 側でパニックしない、全てJSONエラーレスポンスで返す |
| ログ | Go 標準 `log` パッケージで最低限のリクエストログ |
| パフォーマンス | 特別な対策なし（SQLite のローカルアクセスは十分速い） |

### 8.1 CORS 設定

開発時は Vite devサーバー（:5173）から Go API（:5522）へクロスオリジンリクエストが発生する。

| 項目 | 値 |
|------|-----|
| `Access-Control-Allow-Origin` | `http://localhost:5173`（開発時のみ） |
| `Access-Control-Allow-Methods` | `GET, POST, OPTIONS` |
| `Access-Control-Allow-Headers` | `Content-Type` |

- 本番（embed配信）では同一オリジンなので CORS 不要
- Go 側のミドルウェアで、リクエストの Origin を見て開発時のみヘッダーを付与する
- `OPTIONS` プリフライトリクエストには HTTP 204 を返す

### 8.2 静的ファイル配信

| モード | フロント | API |
|--------|---------|-----|
| 開発時 | `pnpm dev`（Vite devサーバー :5173） | Go サーバー :5522。Vite 側で `/api` を :5522 にプロキシ |
| 本番 | `pnpm build` → `frontend/dist/` 生成 | Go が `embed.FS` で `frontend/dist/` を埋め込み、`/` で配信 |

開発時はフロント・バックエンドを別プロセスで起動する。
本番は `go build` で単一バイナリにフロントを埋め込み、1プロセスで完結する。

---

## 9. 開発ワークフロー

### GitHub Issue 駆動開発

Phase 1 の全作業を GitHub Issue に分解し、1 Issue = 1 PR で進める。

```
Issue作成（スマホ可）
  → @claude にメンション or ローカルで claude -p
    → Claude が CLAUDE.md + 仕様書を読んで実装
      → PR 自動作成
        → テスト通過確認 → マージ
```

### Issue の書き方ルール

- 1 Issue = 1つの小さなタスク（1機能 or 1エンドポイント）
- タイトルは具体的に（例: `POST /api/query エンドポイント実装`）
- 本文に受け入れ条件を明記
- 関連する仕様書を参照（`docs/spec/phase1-foundation.md` のセクション番号）

### Issue 一覧（Phase 1）

| # | タイトル | スコープ | 依存 |
|---|---------|---------|------|
| 1 | プロジェクト基盤: Wails除去 + Go HTTPサーバー構築 | main.go, server/router.go, database/ | なし |
| 2 | POST /api/query エンドポイント実装 | internal/query/ 一式 + テスト | #1 |
| 3 | POST /api/exec エンドポイント実装 | internal/exec/ 一式 + テスト | #1 |
| 4 | GET /api/schema/tables エンドポイント実装 | internal/schema/ 一式 + テスト | #1 |
| 5 | GET /api/health エンドポイント実装 | server/router.go に追加 | #1 |
| 6 | フロントエンド基盤: Wails除去 + fetch API + Vite設定 | shared/, App.vue, vite.config.ts | #1 |
| 7 | SQLエディタ機能（Monaco Editor） | features/sql-editor/ | #6 |
| 8 | 結果パネル機能 | features/results-panel/ | #6, #2 |
| 9 | スキーマサイドバー機能 | features/schema-sidebar/ | #6, #4 |
| 10 | 統合: Ctrl+Enter でSQL実行 → 結果表示 | App.vue で各機能を接続 | #7, #8, #9 |

---

## 10. Phase 1 完了条件

### サーバー
- [ ] `go run .` でHTTPサーバーがポート 5522 で起動する
- [ ] 不正なリクエスト（空JSON、sqlフィールドなし）に HTTP 400 を返す
- [ ] `go test ./internal/...` で全テストがパスする

### フロントエンド
- [ ] ブラウザで `http://localhost:5522` を開くとSQLエディタが表示される（ダークテーマ）
- [ ] SQLを入力してCtrl+Enterで実行できる
- [ ] SQL先頭キーワードで query/exec APIが自動振り分けされる
- [ ] SELECT結果がテーブル形式で表示される（カラムヘッダ固定、スクロール対応）
- [ ] NULL は「NULL」とグレーアウト+イタリックで表示される
- [ ] INSERT/UPDATE/DELETE/DDL の実行結果（affected rows）がステータスバーに表示される
- [ ] SQLエラー時にエラーメッセージがステータスバーに表示される
- [ ] 実行時間・取得行数がステータスバーに表示される
- [ ] 実行中はボタン無効化+スピナー表示
- [ ] サイドバーにテーブル一覧が表示される
- [ ] テーブル名クリックで `SELECT * FROM {table} LIMIT 100` がエディタに挿入される
- [ ] exec 実行後にサイドバーが自動更新される
- [ ] `pnpm build` が成功する

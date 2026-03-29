# Phase 4: マルチDB対応 — 仕様書

## 概要

PostgreSQL・MySQL・SQL Server への接続をサポートし、SQLite 以外のデータベースでも OrangeSQL を使えるようにする。
既存の `database.Database` インターフェースを活かし、各 DB ドライバの実装を追加する。

---

## 1. 対応データベース

| DB | ドライバパッケージ | driver 名 |
|----|-------------------|-----------|
| SQLite | `modernc.org/sqlite`（既存） | `sqlite` |
| PostgreSQL | `github.com/jackc/pgx/v5/stdlib` | `postgres` |
| MySQL | `github.com/go-sql-driver/mysql` | `mysql` |
| SQL Server | `github.com/microsoft/go-mssqldb` | `sqlserver` |

すべて `database/sql` 標準インターフェースを使用する。CGO=0 を維持するため、CGO 不要のドライバを選択する。

---

## 2. バックエンド変更

### 2.1 Database ファクトリ

既存の `database.NewSQLite()` を直接呼ぶ箇所をファクトリに置き換える。

```go
// internal/database/factory.go
func New(driver string, params ConnectionParams) (Database, error)
```

```go
// internal/database/types.go に追加
type ConnectionParams struct {
    // SQLite
    Path string `json:"path,omitempty"`

    // TCP 系共通（PostgreSQL / MySQL / SQL Server）
    Host     string `json:"host,omitempty"`
    Port     int    `json:"port,omitempty"`
    User     string `json:"user,omitempty"`
    Password string `json:"password,omitempty"`
    DBName   string `json:"dbName,omitempty"`

    // SSL/TLS
    SSLMode string `json:"sslMode,omitempty"` // disable, require 等
}
```

ファクトリは `driver` を見て対応する実装を生成する。不正な `driver` 値はエラーを返す。

### 2.2 各 DB ドライバの実装

`database.Database` インターフェースを各 DB 向けに実装する。

```
internal/database/
├── interface.go          # Database インターフェース（既存・変更なし）
├── types.go              # ConnectionParams 追加
├── factory.go            # New() ファクトリ
├── sqlite.go             # SQLite 実装（既存・変更なし）
├── postgres.go           # PostgreSQL 実装
├── mysql.go              # MySQL 実装
└── sqlserver.go          # SQL Server 実装
```

各実装が対応するメソッド:

#### Query / Exec

`database/sql` の標準 API を使用するため、既存の SQLite 実装とほぼ同じコード。
`scanRows()` ヘルパーを共通化し、各ドライバで再利用する。

```go
// internal/database/scan.go
func scanRows(rows *sql.Rows) (QueryResult, error)
```

#### Tables()

| DB | クエリ |
|----|--------|
| SQLite | `SELECT name, type FROM sqlite_master WHERE type IN ('table','view')` |
| PostgreSQL | `SELECT table_name, table_type FROM information_schema.tables WHERE table_schema = 'public'` |
| MySQL | `SELECT table_name, table_type FROM information_schema.tables WHERE table_schema = DATABASE()` |
| SQL Server | `SELECT TABLE_NAME, TABLE_TYPE FROM INFORMATION_SCHEMA.TABLES` |

#### Columns(tableName)

| DB | クエリ |
|----|--------|
| SQLite | `PRAGMA table_info(...)` |
| PostgreSQL | `information_schema.columns` + `pg_constraint`（PK 判定） |
| MySQL | `information_schema.columns`（`COLUMN_KEY = 'PRI'` で PK 判定） |
| SQL Server | `INFORMATION_SCHEMA.COLUMNS` + `INFORMATION_SCHEMA.KEY_COLUMN_USAGE`（PK 判定） |

#### Name()

- SQLite: ファイルパス（既存）
- TCP 系: `host:port/dbname` 形式の文字列

### 2.3 Profile 型の拡張

```go
// internal/profile/types.go
type Profile struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    Driver    string `json:"driver"`
    CreatedAt string `json:"createdAt"`

    // SQLite
    Path string `json:"path,omitempty"`

    // TCP 系共通
    Host     string `json:"host,omitempty"`
    Port     int    `json:"port,omitempty"`
    User     string `json:"user,omitempty"`
    Password string `json:"password,omitempty"`
    DBName   string `json:"dbName,omitempty"`
    SSLMode  string `json:"sslMode,omitempty"`
}
```

- `Path` は SQLite のみ使用
- TCP 系は `Host`, `Port`, `User`, `Password`, `DBName` を使用
- `Password` は profiles.json にプレーンテキストで保存する（ローカルツールなのでシンプルに）

Profile → ConnectionParams 変換メソッドを追加:

```go
func (p Profile) ToConnectionParams() database.ConnectionParams
```

### 2.4 ConnectionManager の修正

```go
// internal/profile/service.go
func (m *ConnectionManager) SwitchTo(p Profile) error {
    // 変更前: database.NewSQLite(p.Path)
    // 変更後: database.New(p.Driver, p.ToConnectionParams())
    newDB, err := database.New(p.Driver, p.ToConnectionParams())
    if err != nil {
        return fmt.Errorf("failed to connect: %w", err)
    }
    // ... 既存のスイッチロジック
}
```

### 2.5 main.go の修正

起動時の DB 接続も同様にファクトリを使用する:

```go
// 変更前: database.NewSQLite(connectProfile.Path)
// 変更後: database.New(connectProfile.Driver, connectProfile.ToConnectionParams())
```

### 2.6 デフォルトポート

| DB | デフォルトポート |
|----|----------------|
| PostgreSQL | 5432 |
| MySQL | 3306 |
| SQL Server | 1433 |

フロントエンドでドライバ切り替え時にデフォルト値をセットする。バックエンドでは `Port == 0` の場合にデフォルトポートを適用する。

### 2.7 バリデーション

プロファイル作成・更新時にドライバごとの必須フィールドを検証する:

| DB | 必須フィールド |
|----|---------------|
| SQLite | `path` |
| PostgreSQL / MySQL / SQL Server | `host`, `dbName` |

`user`, `password` は任意（DB 側の設定による）。

---

## 3. フロントエンド変更

### 3.1 接続ダイアログの拡張

ドライバ選択で表示するフィールドを切り替える:

```
┌─── 新規接続 ──────────────────┐
│                               │
│  接続名:  [開発DB           ] │
│                               │
│  ドライバ: [PostgreSQL ▼]     │
│                               │
│  ホスト:    [localhost       ] │
│  ポート:    [5432            ] │
│  ユーザー:  [postgres        ] │
│  パスワード:[••••••          ] │
│  データベース:[mydb          ] │
│  SSL:       [disable ▼]      │
│                               │
│      [キャンセル]  [接続]     │
└───────────────────────────────┘
```

SQLite 選択時:

```
┌─── 新規接続 ──────────────────┐
│                               │
│  接続名:  [開発DB           ] │
│                               │
│  ドライバ: [SQLite ▼]        │
│                               │
│  ファイルパス:                 │
│  [C:/work/dev.db           ]  │
│                               │
│      [キャンセル]  [接続]     │
└───────────────────────────────┘
```

### 3.2 ドライバ選択肢

```typescript
const DRIVERS = [
  { value: 'sqlite', label: 'SQLite' },
  { value: 'postgres', label: 'PostgreSQL' },
  { value: 'mysql', label: 'MySQL' },
  { value: 'sqlserver', label: 'SQL Server' },
] as const;
```

### 3.3 型定義の拡張

```typescript
// features/connection/types.ts
type ConnectionProfile = {
  id: string;
  name: string;
  driver: 'sqlite' | 'postgres' | 'mysql' | 'sqlserver';
  createdAt: string;

  // SQLite
  path?: string;

  // TCP 系共通
  host?: string;
  port?: number;
  user?: string;
  password?: string;
  dbName?: string;
  sslMode?: string;
};
```

### 3.4 接続管理ダイアログの変更

ドライバ名をプロファイル一覧に表示する:

```
┌─── 接続管理 ─────────────────────────┐
│                                       │
│  開発DB      PostgreSQL  [編集][削除] │
│  ローカル    SQLite      [編集][削除] │
│  テスト      MySQL       [編集][削除] │
│                                       │
│                    [閉じる]           │
└───────────────────────────────────────┘
```

### 3.5 API クライアントの変更

`createProfile` / `updateProfile` のリクエストボディに新しいフィールドを含める。
レスポンス型も新しいフィールドに対応させる。

---

## 4. テスト

### 4.1 バックエンド単体テスト

| テストファイル | テスト内容 |
|--------------|-----------|
| `database/factory_test.go` | ファクトリが driver ごとに正しい実装を返す |
| `database/postgres_test.go` | PostgreSQL の Tables/Columns クエリが正しい SQL を生成する |
| `database/mysql_test.go` | MySQL の Tables/Columns クエリが正しい SQL を生成する |
| `database/sqlserver_test.go` | SQL Server の Tables/Columns クエリが正しい SQL を生成する |
| `database/scan_test.go` | 共通 scanRows の動作確認 |
| `profile/handler_test.go` | 新フィールドの CRUD・バリデーション |

**テスト方針:**
- 実際の DB サーバーが不要なテストを優先する
- ファクトリのルーティングテスト: ドライバ名 → 正しい型が返ること
- スキーマクエリのテスト: `database/sql` を直接叩かず、生成される SQL 文字列が正しいことを確認する方法を検討
- SQLite は引き続きインメモリ DB でのテストが可能
- PostgreSQL / MySQL / SQL Server のインテグレーションテストは `//go:build integration` タグで分離し、CI やローカルに DB がある環境でのみ実行する

### 4.2 E2E テスト

| テスト | 内容 |
|--------|------|
| ドライバ切り替え UI | ドライバ選択でフォームフィールドが切り替わる |
| PostgreSQL プロファイル作成 | TCP 系フィールドで保存できる |
| バリデーション | 必須フィールド未入力時にエラーが出る |

E2E テストでは SQLite のみ実接続テストを行う（外部 DB 不要にするため）。
TCP 系はフロントエンドの UI 切り替えとバリデーションのみテストする。

---

## 5. Issue 一覧（Phase 4）

| # | タイトル | スコープ | 依存 |
|---|---------|---------|------|
| 23 | scanRows 共通化 + ConnectionParams 型追加 | database/scan.go, database/types.go | なし |
| 24 | Database ファクトリ関数 | database/factory.go | #23 |
| 25 | PostgreSQL ドライバ実装 | database/postgres.go | #23 |
| 26 | MySQL ドライバ実装 | database/mysql.go | #23 |
| 27 | SQL Server ドライバ実装 | database/sqlserver.go | #23 |
| 28 | Profile 型拡張 + バリデーション | profile/types.go | なし |
| 29 | ConnectionManager ファクトリ対応 + main.go 修正 | profile/service.go, main.go | #24, #28 |
| 30 | フロント: 接続ダイアログのマルチDB対応 | features/connection/ | #28 |
| 31 | テスト追加（Phase 4 機能） | database/*_test.go, e2e/ | #29, #30 |

---

## 6. 変更不要な箇所

以下は既に `database.Database` インターフェース経由で抽象化されており、変更不要:

- `internal/query/handler.go` — Query ハンドラ
- `internal/query/service.go` — Query レスポンスビルダー
- `internal/exec/handler.go` — Exec ハンドラ
- `internal/exec/service.go` — Exec レスポンスビルダー
- `internal/schema/handler.go` — Schema ハンドラ
- `internal/schema/service.go` — Schema レスポンスビルダー
- `internal/server/router.go` — ルーター

---

## 7. Phase 4 完了条件

### サーバー
- [ ] `database.New(driver, params)` ファクトリが 4 種のドライバに対応する
- [ ] PostgreSQL で Tables / Columns / Query / Exec が動作する
- [ ] MySQL で Tables / Columns / Query / Exec が動作する
- [ ] SQL Server で Tables / Columns / Query / Exec が動作する
- [ ] Profile に TCP 系接続パラメータが保存される
- [ ] ConnectionManager がドライバに応じた DB 接続を生成する
- [ ] プロファイル作成時にドライバ別の必須フィールドがバリデーションされる
- [ ] `go test ./internal/...` で全テストがパスする

### フロントエンド
- [ ] ドライバ選択でフォームフィールドが動的に切り替わる
- [ ] PostgreSQL / MySQL / SQL Server のプロファイルが作成できる
- [ ] 接続管理ダイアログにドライバ名が表示される
- [ ] 必須フィールド未入力時にバリデーションエラーが出る
- [ ] `pnpm build` が成功する
- [ ] `pnpm test:e2e` で全テストがパスする

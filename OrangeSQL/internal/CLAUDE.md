# バックエンド（internal/）

## アーキテクチャ

バーティカルスライス構成。各機能パッケージに handler（HTTP）、service（ロジック）、types（型定義）をコロケーション。

```
database/      DB 抽象化層
├── interface.go    Database インターフェース（Query, Exec, ExecBatch, Tables, Columns）
├── scan.go         共通ヘルパー（scanRows, execStatement, execBatch）
├── params.go       ConnectionParams + DSN 生成
├── factory.go      New(driver, params) ファクトリ
├── sqlite.go       SQLite 実装
├── postgres.go     PostgreSQL 実装
├── mysql.go        MySQL 実装
├── sqlserver.go    SQL Server 実装
└── integration_test.go  Docker 統合テスト（//go:build integration）

query/         POST /api/query
├── handler.go      HTTP ハンドラ
├── service.go      BuildResponse（結果変換 + 編集可能テーブル判定）
├── detect.go       DetectEditableTable（単一テーブル SELECT 判定 + PK 検出）
└── types.go        Response（editableTable, pkColumns 含む）

exec/          POST /api/exec, /api/exec/batch
├── handler.go        単一 SQL 実行
├── batch_handler.go  バッチ実行（トランザクション）
└── types.go          Request, Response, BatchRequest, BatchResponse

schema/        GET /api/schema/*
├── handler.go      Tables / Columns / Completions ハンドラ
├── service.go      レスポンスビルダー
└── types.go        型定義

profile/       接続プロファイル管理
├── handler.go      CRUD + Connect ハンドラ
├── service.go      ConnectionManager（スレッドセーフな DB 切り替え）
├── storage.go      profiles.json 読み書き
└── types.go        Profile, Validate(), ToConnectionParams()

server/        ルーター + ミドルウェア
└── router.go       NewRouter(cm, staticFS) — API + 静的ファイル配信

testutil/      テスト用
└── testdb.go       NewTestDB() — インメモリ SQLite + サンプルデータ
```

## ルール

- 標準ライブラリ優先（外部依存は DB ドライバのみ）
- エラーは全て JSON レスポンス `{"error": "..."}` で返す
- DB アクセスは必ず `database.Database` インターフェース経由
- 新しいドライバを追加する場合: interface.go のメソッドを実装 + factory.go に登録
- テストは `httptest` + `testutil.NewTestDB()`（インメモリ SQLite）
- Docker 統合テストは `//go:build integration` タグで分離
- `sync.RWMutex` で DB 参照を保護（ConnectionManager）

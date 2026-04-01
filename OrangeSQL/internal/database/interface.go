package database

// Database はDB種別を抽象化するインターフェース。
// 各スライス（query/exec/schema）はこのインターフェースを受け取る。
type Database interface {
	// Query は SELECT 系 SQL を実行し、結果セットを返す。
	Query(sql string) (QueryResult, error)
	// Exec は INSERT/UPDATE/DELETE/DDL を実行し、影響行数を返す。
	Exec(sql string) (ExecResult, error)
	// ExecBatch は複数 SQL をトランザクション内で一括実行する。
	// いずれかが失敗した場合はロールバックしてエラーを返す。
	ExecBatch(statements []string) (totalAffected int64, err error)
	// Tables はテーブル・ビュー一覧を返す。
	Tables() ([]TableInfo, error)
	// Columns は指定テーブルのカラム情報を返す。
	Columns(tableName string) ([]ColumnInfo, error)
	// Name は接続中のDB名（ファイルパス等）を返す。
	Name() string
	// Close は接続を閉じる。
	Close() error
}

// QueryResult は SELECT 系の実行結果。
type QueryResult struct {
	Columns []string
	Rows    [][]*string // nil = SQL NULL
}

// ExecResult は DML/DDL の実行結果。
type ExecResult struct {
	AffectedRows int64
}

// TableInfo はテーブルまたはビューの情報。
type TableInfo struct {
	Name string // テーブル名
	Type string // "table" or "view"
}

// ColumnInfo はテーブルのカラム情報。
type ColumnInfo struct {
	Name    string
	Type    string
	PK      bool
	NotNull bool
	Comment string // カラムコメント（論理名）。非対応DBでは空文字。
}

package testutil

import (
	"testing"

	"OrangeSQL/internal/database"
)

// NewTestDB はテスト用のインメモリ SQLite を作成し、サンプルデータを投入する。
func NewTestDB(t *testing.T) database.Database {
	t.Helper()
	db, err := database.NewSQLite(":memory:")
	if err != nil {
		t.Fatalf("failed to create test db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	setup := []string{
		`CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT
		)`,
		`INSERT INTO users (name, email) VALUES ('Alice', 'alice@example.com')`,
		`INSERT INTO users (name, email) VALUES ('Bob', NULL)`,
	}
	for _, s := range setup {
		if _, err := db.Exec(s); err != nil {
			t.Fatalf("failed to setup test db: %v", err)
		}
	}
	return db
}

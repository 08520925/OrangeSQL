//go:build integration

package database_test

import (
	"testing"

	"OrangeSQL/internal/database"
)

// テスト用の共通ヘルパー: テーブル作成 → CRUD → スキーマ確認
func runDriverTest(t *testing.T, db database.Database) {
	t.Helper()

	// 0. 前回の残りを掃除
	db.Exec("DROP TABLE IF EXISTS integration_test")

	// 1. テーブル作成
	_, err := db.Exec(`CREATE TABLE integration_test (
		id INTEGER PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(200)
	)`)
	if err != nil {
		t.Fatalf("CREATE TABLE failed: %v", err)
	}
	t.Cleanup(func() {
		db.Exec("DROP TABLE IF EXISTS integration_test")
	})

	// 2. INSERT
	res, err := db.Exec("INSERT INTO integration_test (id, name, email) VALUES (1, 'Alice', 'alice@test.com')")
	if err != nil {
		t.Fatalf("INSERT failed: %v", err)
	}
	if res.AffectedRows != 1 {
		t.Fatalf("expected 1 affected row, got %d", res.AffectedRows)
	}

	res, err = db.Exec("INSERT INTO integration_test (id, name) VALUES (2, 'Bob')")
	if err != nil {
		t.Fatalf("INSERT (NULL email) failed: %v", err)
	}

	// 3. SELECT
	qr, err := db.Query("SELECT id, name, email FROM integration_test ORDER BY id")
	if err != nil {
		t.Fatalf("SELECT failed: %v", err)
	}
	if len(qr.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(qr.Columns))
	}
	if len(qr.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(qr.Rows))
	}

	// Alice の email は non-nil
	if qr.Rows[0][2] == nil {
		t.Fatal("expected Alice's email to be non-nil")
	}
	if *qr.Rows[0][2] != "alice@test.com" {
		t.Fatalf("expected 'alice@test.com', got '%s'", *qr.Rows[0][2])
	}

	// Bob の email は NULL
	if qr.Rows[1][2] != nil {
		t.Fatalf("expected Bob's email to be nil, got '%s'", *qr.Rows[1][2])
	}

	// 4. UPDATE
	res, err = db.Exec("UPDATE integration_test SET email = 'bob@test.com' WHERE id = 2")
	if err != nil {
		t.Fatalf("UPDATE failed: %v", err)
	}
	if res.AffectedRows != 1 {
		t.Fatalf("expected 1 affected row, got %d", res.AffectedRows)
	}

	// 5. DELETE
	res, err = db.Exec("DELETE FROM integration_test WHERE id = 1")
	if err != nil {
		t.Fatalf("DELETE failed: %v", err)
	}
	if res.AffectedRows != 1 {
		t.Fatalf("expected 1 affected row, got %d", res.AffectedRows)
	}

	// 6. Tables()
	tables, err := db.Tables()
	if err != nil {
		t.Fatalf("Tables() failed: %v", err)
	}
	found := false
	for _, tbl := range tables {
		if tbl.Name == "integration_test" && tbl.Type == "table" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected 'integration_test' in tables, got %v", tables)
	}

	// 7. Columns()
	cols, err := db.Columns("integration_test")
	if err != nil {
		t.Fatalf("Columns() failed: %v", err)
	}
	if len(cols) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(cols))
	}
	// id カラム: PK
	if cols[0].Name != "id" || !cols[0].PK {
		t.Fatalf("expected id column with PK, got %+v", cols[0])
	}
	// name カラム: NOT NULL
	if cols[1].Name != "name" || !cols[1].NotNull {
		t.Fatalf("expected name column with NOT NULL, got %+v", cols[1])
	}
	// email カラム: nullable
	if cols[2].Name != "email" || cols[2].NotNull {
		t.Fatalf("expected email column nullable, got %+v", cols[2])
	}

	// 8. Name()
	name := db.Name()
	if name == "" {
		t.Fatal("expected non-empty name")
	}
	t.Logf("DB Name: %s", name)

	// 9. エラーケース
	_, err = db.Query("SELECT * FROM nonexistent_table_xyz")
	if err == nil {
		t.Fatal("expected error for nonexistent table")
	}

	_, err = db.Columns("nonexistent_table_xyz")
	if err == nil {
		t.Fatal("expected error for nonexistent table columns")
	}
}

func TestPostgres_Integration(t *testing.T) {
	params := database.ConnectionParams{
		Host:     "localhost",
		Port:     15432,
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
		SSLMode:  "disable",
	}
	db, err := database.New("postgres", params)
	if err != nil {
		t.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	t.Log("Connected to PostgreSQL")
	runDriverTest(t, db)
}

func TestMySQL_Integration(t *testing.T) {
	params := database.ConnectionParams{
		Host:     "localhost",
		Port:     13306,
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
	}
	db, err := database.New("mysql", params)
	if err != nil {
		t.Fatalf("failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	t.Log("Connected to MySQL")
	runDriverTest(t, db)
}

func TestSQLServer_Integration(t *testing.T) {
	params := database.ConnectionParams{
		Host:     "localhost",
		Port:     11433,
		User:     "sa",
		Password: "TestPass123!",
		DBName:   "master",
	}
	db, err := database.New("sqlserver", params)
	if err != nil {
		t.Fatalf("failed to connect to SQL Server: %v", err)
	}
	defer db.Close()

	t.Log("Connected to SQL Server")
	runDriverTest(t, db)
}

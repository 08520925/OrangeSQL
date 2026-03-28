package database_test

import (
	"testing"

	"OrangeSQL/internal/testutil"
)

func TestSQLite_Query(t *testing.T) {
	db := testutil.NewTestDB(t)

	result, err := db.Query("SELECT id, name, email FROM users ORDER BY id")
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}

	if len(result.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(result.Columns))
	}
	if result.Columns[0] != "id" || result.Columns[1] != "name" || result.Columns[2] != "email" {
		t.Fatalf("unexpected columns: %v", result.Columns)
	}
	if len(result.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(result.Rows))
	}

	// Alice の行: email は non-nil
	if result.Rows[0][2] == nil {
		t.Fatal("expected Alice's email to be non-nil")
	}
	if *result.Rows[0][2] != "alice@example.com" {
		t.Fatalf("expected alice@example.com, got %s", *result.Rows[0][2])
	}

	// Bob の行: email は NULL (nil)
	if result.Rows[1][2] != nil {
		t.Fatalf("expected Bob's email to be nil, got %v", *result.Rows[1][2])
	}
}

func TestSQLite_Query_Empty(t *testing.T) {
	db := testutil.NewTestDB(t)

	result, err := db.Query("SELECT * FROM users WHERE id = -1")
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if len(result.Rows) != 0 {
		t.Fatalf("expected 0 rows, got %d", len(result.Rows))
	}
}

func TestSQLite_Query_Error(t *testing.T) {
	db := testutil.NewTestDB(t)

	_, err := db.Query("SELEC * FROM users")
	if err == nil {
		t.Fatal("expected error for invalid SQL")
	}
}

func TestSQLite_Exec(t *testing.T) {
	db := testutil.NewTestDB(t)

	result, err := db.Exec("INSERT INTO users (name, email) VALUES ('Charlie', 'c@example.com')")
	if err != nil {
		t.Fatalf("exec failed: %v", err)
	}
	if result.AffectedRows != 1 {
		t.Fatalf("expected 1 affected row, got %d", result.AffectedRows)
	}
}

func TestSQLite_Exec_DDL(t *testing.T) {
	db := testutil.NewTestDB(t)

	_, err := db.Exec("CREATE TABLE test_table (id INTEGER)")
	if err != nil {
		t.Fatalf("exec DDL failed: %v", err)
	}
	// DDL の AffectedRows は DB 実装依存のため値は検証しない
}

func TestSQLite_Exec_Error(t *testing.T) {
	db := testutil.NewTestDB(t)

	_, err := db.Exec("INSERT INTO nonexistent (name) VALUES ('test')")
	if err == nil {
		t.Fatal("expected error for nonexistent table")
	}
}

func TestSQLite_Tables(t *testing.T) {
	db := testutil.NewTestDB(t)

	tables, err := db.Tables()
	if err != nil {
		t.Fatalf("tables failed: %v", err)
	}
	// AUTOINCREMENT 付きテーブルがあると sqlite_sequence も含まれる
	found := false
	for _, tbl := range tables {
		if tbl.Name == "users" && tbl.Type == "table" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected 'users' table in results, got %v", tables)
	}
}

func TestSQLite_Name(t *testing.T) {
	db := testutil.NewTestDB(t)

	name := db.Name()
	if name != ":memory:" {
		t.Fatalf("expected ':memory:', got '%s'", name)
	}
}

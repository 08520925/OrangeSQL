package query_test

import (
	"testing"

	"OrangeSQL/internal/query"
	"OrangeSQL/internal/testutil"
)

func TestDetect_SimpleSelect(t *testing.T) {
	db := testutil.NewTestDB(t)
	table, pks := query.DetectEditableTable(
		"SELECT id, name, email FROM users",
		[]string{"id", "name", "email"},
		db,
	)
	if table != "users" {
		t.Fatalf("expected 'users', got '%s'", table)
	}
	if len(pks) != 1 || pks[0] != "id" {
		t.Fatalf("expected ['id'], got %v", pks)
	}
}

func TestDetect_SelectStar(t *testing.T) {
	db := testutil.NewTestDB(t)
	table, pks := query.DetectEditableTable(
		"SELECT * FROM users LIMIT 100",
		[]string{"id", "name", "email"},
		db,
	)
	if table != "users" {
		t.Fatalf("expected 'users', got '%s'", table)
	}
	if len(pks) != 1 {
		t.Fatalf("expected 1 PK, got %d", len(pks))
	}
}

func TestDetect_Join(t *testing.T) {
	db := testutil.NewTestDB(t)
	table, _ := query.DetectEditableTable(
		"SELECT u.id FROM users u JOIN orders o ON u.id = o.user_id",
		[]string{"id"},
		db,
	)
	if table != "" {
		t.Fatalf("expected empty for JOIN, got '%s'", table)
	}
}

func TestDetect_Union(t *testing.T) {
	db := testutil.NewTestDB(t)
	table, _ := query.DetectEditableTable(
		"SELECT id FROM users UNION SELECT id FROM users",
		[]string{"id"},
		db,
	)
	if table != "" {
		t.Fatalf("expected empty for UNION, got '%s'", table)
	}
}

func TestDetect_GroupBy(t *testing.T) {
	db := testutil.NewTestDB(t)
	table, _ := query.DetectEditableTable(
		"SELECT name, COUNT(*) FROM users GROUP BY name",
		[]string{"name"},
		db,
	)
	if table != "" {
		t.Fatalf("expected empty for GROUP BY, got '%s'", table)
	}
}

func TestDetect_NonSelect(t *testing.T) {
	db := testutil.NewTestDB(t)
	table, _ := query.DetectEditableTable(
		"INSERT INTO users (name) VALUES ('test')",
		[]string{},
		db,
	)
	if table != "" {
		t.Fatalf("expected empty for INSERT, got '%s'", table)
	}
}

func TestDetect_PKNotInResult(t *testing.T) {
	db := testutil.NewTestDB(t)
	// PK (id) が結果カラムに含まれていない
	table, _ := query.DetectEditableTable(
		"SELECT name, email FROM users",
		[]string{"name", "email"},
		db,
	)
	if table != "" {
		t.Fatalf("expected empty when PK not in result, got '%s'", table)
	}
}

func TestDetect_NonexistentTable(t *testing.T) {
	db := testutil.NewTestDB(t)
	table, _ := query.DetectEditableTable(
		"SELECT * FROM nonexistent",
		[]string{"id"},
		db,
	)
	if table != "" {
		t.Fatalf("expected empty for nonexistent table, got '%s'", table)
	}
}

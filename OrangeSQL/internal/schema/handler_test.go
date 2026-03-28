package schema_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"OrangeSQL/internal/database"
	"OrangeSQL/internal/schema"
)

func TestSchema_Tables(t *testing.T) {
	db, err := database.NewSQLite(":memory:")
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	defer db.Close()

	db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)")
	db.Exec("CREATE TABLE orders (id INTEGER PRIMARY KEY, user_id INTEGER)")
	db.Exec("CREATE VIEW user_view AS SELECT id, name FROM users")

	handler := schema.TablesHandler(db)
	req := httptest.NewRequest(http.MethodGet, "/api/schema/tables", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp schema.Response
	json.Unmarshal(w.Body.Bytes(), &resp)

	if len(resp.Tables) != 3 {
		t.Fatalf("expected 3 entries, got %d: %v", len(resp.Tables), resp.Tables)
	}

	// name 順で返る（orders, user_view, users）
	expected := []struct {
		name string
		typ  string
	}{
		{"orders", "table"},
		{"user_view", "view"},
		{"users", "table"},
	}
	for i, e := range expected {
		if resp.Tables[i].Name != e.name || resp.Tables[i].Type != e.typ {
			t.Fatalf("entry %d: expected %s/%s, got %s/%s",
				i, e.name, e.typ, resp.Tables[i].Name, resp.Tables[i].Type)
		}
	}
}

func TestSchema_EmptyDB(t *testing.T) {
	db, err := database.NewSQLite(":memory:")
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	defer db.Close()

	handler := schema.TablesHandler(db)
	req := httptest.NewRequest(http.MethodGet, "/api/schema/tables", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp schema.Response
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Tables != nil && len(resp.Tables) != 0 {
		t.Fatalf("expected empty tables, got %v", resp.Tables)
	}
}

func TestSchema_Columns(t *testing.T) {
	db, err := database.NewSQLite(":memory:")
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	defer db.Close()

	db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT NOT NULL, email TEXT)")

	handler := schema.ColumnsHandler(db)
	req := httptest.NewRequest(http.MethodGet, "/api/schema/columns?table=users", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp schema.ColumnsResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if len(resp.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(resp.Columns))
	}
	if resp.Columns[0].Name != "id" || !resp.Columns[0].PK {
		t.Fatalf("unexpected id column: %+v", resp.Columns[0])
	}
	if resp.Columns[1].Name != "name" || !resp.Columns[1].NotNull {
		t.Fatalf("unexpected name column: %+v", resp.Columns[1])
	}
}

func TestSchema_Columns_NotFound(t *testing.T) {
	db, err := database.NewSQLite(":memory:")
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	defer db.Close()

	handler := schema.ColumnsHandler(db)
	req := httptest.NewRequest(http.MethodGet, "/api/schema/columns?table=nonexistent", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestSchema_Columns_MissingParam(t *testing.T) {
	db, err := database.NewSQLite(":memory:")
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	defer db.Close()

	handler := schema.ColumnsHandler(db)
	req := httptest.NewRequest(http.MethodGet, "/api/schema/columns", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

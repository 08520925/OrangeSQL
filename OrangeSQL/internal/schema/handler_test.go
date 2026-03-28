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

	handler := schema.Handler(db)
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

	handler := schema.Handler(db)
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

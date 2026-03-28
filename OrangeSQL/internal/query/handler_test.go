package query_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"OrangeSQL/internal/query"
	"OrangeSQL/internal/testutil"
)

func TestQuery_Success(t *testing.T) {
	db := testutil.NewTestDB(t)
	handler := query.Handler(db)

	body := `{"sql": "SELECT id, name, email FROM users ORDER BY id"}`
	req := httptest.NewRequest(http.MethodPost, "/api/query", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp query.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if len(resp.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(resp.Columns))
	}
	if resp.RowCount != 2 {
		t.Fatalf("expected 2 rows, got %d", resp.RowCount)
	}
	if resp.ExecutionTime < 0 {
		t.Fatal("execution time should be >= 0")
	}
}

func TestQuery_NullValue(t *testing.T) {
	db := testutil.NewTestDB(t)
	handler := query.Handler(db)

	body := `{"sql": "SELECT email FROM users WHERE name = 'Bob'"}`
	req := httptest.NewRequest(http.MethodPost, "/api/query", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp query.Response
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.RowCount != 1 {
		t.Fatalf("expected 1 row, got %d", resp.RowCount)
	}
	// Bob の email は NULL → JSON null
	if resp.Rows[0][0] != nil {
		t.Fatalf("expected null for Bob's email, got %v", resp.Rows[0][0])
	}
}

func TestQuery_EmptyResult(t *testing.T) {
	db := testutil.NewTestDB(t)
	handler := query.Handler(db)

	body := `{"sql": "SELECT * FROM users WHERE id = -1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/query", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp query.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.RowCount != 0 {
		t.Fatalf("expected 0 rows, got %d", resp.RowCount)
	}
}

func TestQuery_SQLError(t *testing.T) {
	db := testutil.NewTestDB(t)
	handler := query.Handler(db)

	body := `{"sql": "SELEC * FROM users"}`
	req := httptest.NewRequest(http.MethodPost, "/api/query", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	var resp query.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Error == "" {
		t.Fatal("expected error message")
	}
}

func TestQuery_EmptySQL(t *testing.T) {
	db := testutil.NewTestDB(t)
	handler := query.Handler(db)

	body := `{"sql": ""}`
	req := httptest.NewRequest(http.MethodPost, "/api/query", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestQuery_InvalidJSON(t *testing.T) {
	db := testutil.NewTestDB(t)
	handler := query.Handler(db)

	req := httptest.NewRequest(http.MethodPost, "/api/query", strings.NewReader("not json"))
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

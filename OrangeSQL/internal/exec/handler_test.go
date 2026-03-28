package exec_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"OrangeSQL/internal/exec"
	"OrangeSQL/internal/testutil"
)

func TestExec_Insert(t *testing.T) {
	db := testutil.NewTestDB(t)
	handler := exec.Handler(db)

	body := `{"sql": "INSERT INTO users (name, email) VALUES ('Charlie', 'c@example.com')"}`
	req := httptest.NewRequest(http.MethodPost, "/api/exec", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp exec.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.AffectedRows != 1 {
		t.Fatalf("expected 1 affected row, got %d", resp.AffectedRows)
	}
	if resp.ExecutionTime < 0 {
		t.Fatal("execution time should be >= 0")
	}
}

func TestExec_Update(t *testing.T) {
	db := testutil.NewTestDB(t)
	handler := exec.Handler(db)

	body := `{"sql": "UPDATE users SET email = 'updated@example.com' WHERE name = 'Alice'"}`
	req := httptest.NewRequest(http.MethodPost, "/api/exec", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp exec.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.AffectedRows != 1 {
		t.Fatalf("expected 1 affected row, got %d", resp.AffectedRows)
	}
}

func TestExec_Delete(t *testing.T) {
	db := testutil.NewTestDB(t)
	handler := exec.Handler(db)

	body := `{"sql": "DELETE FROM users WHERE name = 'Bob'"}`
	req := httptest.NewRequest(http.MethodPost, "/api/exec", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp exec.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.AffectedRows != 1 {
		t.Fatalf("expected 1 affected row, got %d", resp.AffectedRows)
	}
}

func TestExec_SQLError(t *testing.T) {
	db := testutil.NewTestDB(t)
	handler := exec.Handler(db)

	body := `{"sql": "INSERT INTO nonexistent (name) VALUES ('test')"}`
	req := httptest.NewRequest(http.MethodPost, "/api/exec", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	var resp exec.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Error == "" {
		t.Fatal("expected error message")
	}
}

func TestExec_EmptySQL(t *testing.T) {
	db := testutil.NewTestDB(t)
	handler := exec.Handler(db)

	body := `{"sql": "  "}`
	req := httptest.NewRequest(http.MethodPost, "/api/exec", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestExec_InvalidJSON(t *testing.T) {
	db := testutil.NewTestDB(t)
	handler := exec.Handler(db)

	req := httptest.NewRequest(http.MethodPost, "/api/exec", strings.NewReader("{bad"))
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

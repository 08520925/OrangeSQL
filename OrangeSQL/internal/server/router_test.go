package server_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"OrangeSQL/internal/profile"
	"OrangeSQL/internal/server"
	"OrangeSQL/internal/testutil"
)

func setupCM(t *testing.T) *profile.ConnectionManager {
	t.Helper()
	db := testutil.NewTestDB(t)
	dir := t.TempDir()
	s, err := profile.NewStorage(dir + "/profiles.json")
	if err != nil {
		t.Fatalf("new storage: %v", err)
	}
	store := profile.ProfileStore{
		Profiles: []profile.Profile{
			profile.NewProfile("1", profile.CreateRequest{Name: "Test", Driver: "sqlite", Path: ":memory:"}),
		},
		LastUsedID: "1",
	}
	if err := s.Save(store); err != nil {
		t.Fatalf("save: %v", err)
	}
	return profile.NewConnectionManager(s, db, "1")
}

func TestRouter_VerboseDisabled_NoLog(t *testing.T) {
	cm := setupCM(t)
	defer cm.Close()

	router := server.NewRouter(cm, nil)

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no log output, got: %s", buf.String())
	}
}

func TestRouter_VerboseEnabled_WritesLog(t *testing.T) {
	cm := setupCM(t)
	defer cm.Close()

	router := server.NewRouter(cm, nil, server.WithVerbose())

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(buf.String(), "GET /api/health") {
		t.Fatalf("expected log to contain 'GET /api/health', got: %s", buf.String())
	}
}

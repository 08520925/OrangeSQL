package profile_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"OrangeSQL/internal/profile"
)

func tempStoragePath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "profiles.json")
}

func TestStorage_LoadEmpty(t *testing.T) {
	s, err := profile.NewStorage(tempStoragePath(t))
	if err != nil {
		t.Fatalf("new storage: %v", err)
	}

	store, err := s.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(store.Profiles) != 0 {
		t.Fatalf("expected 0 profiles, got %d", len(store.Profiles))
	}
}

func TestStorage_SaveAndLoad(t *testing.T) {
	path := tempStoragePath(t)
	s, _ := profile.NewStorage(path)

	store := profile.ProfileStore{
		Profiles: []profile.Profile{
			profile.NewProfile("1", profile.CreateRequest{Name: "TestDB", Driver: "sqlite", Path: "/tmp/test.db"}),
		},
		LastUsedID: "1",
	}

	if err := s.Save(store); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := s.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(loaded.Profiles))
	}
	if loaded.Profiles[0].Name != "TestDB" {
		t.Fatalf("expected name 'TestDB', got '%s'", loaded.Profiles[0].Name)
	}
	if loaded.LastUsedID != "1" {
		t.Fatalf("expected lastUsedId '1', got '%s'", loaded.LastUsedID)
	}
}

func TestStorage_CreateDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "profiles.json")

	s, err := profile.NewStorage(path)
	if err != nil {
		t.Fatalf("new storage: %v", err)
	}

	store := profile.ProfileStore{
		Profiles: []profile.Profile{
			profile.NewProfile("1", profile.CreateRequest{Name: "DB", Driver: "sqlite", Path: "/tmp/db"}),
		},
	}
	if err := s.Save(store); err != nil {
		t.Fatalf("save: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected file to exist")
	}
}

func TestNextID(t *testing.T) {
	store := profile.ProfileStore{
		Profiles: []profile.Profile{
			{ID: "1"}, {ID: "3"}, {ID: "2"},
		},
	}
	id := profile.NextID(store)
	if id != "4" {
		t.Fatalf("expected '4', got '%s'", id)
	}
}

func TestNextID_Empty(t *testing.T) {
	store := profile.ProfileStore{}
	id := profile.NextID(store)
	if id != "1" {
		t.Fatalf("expected '1', got '%s'", id)
	}
}

func TestDataDir_UsesUserConfigDir(t *testing.T) {
	dir := profile.DataDir()
	if dir == "" {
		t.Fatal("DataDir returned empty string")
	}
	// OS 標準のデータディレクトリ配下に OrangeSQL が含まれること
	if !filepath.IsAbs(dir) {
		t.Fatalf("expected absolute path, got: %s", dir)
	}
	if filepath.Base(dir) != "OrangeSQL" {
		t.Fatalf("expected dir to end with 'OrangeSQL', got: %s", dir)
	}
}

func TestDefaultPath_UnderDataDir(t *testing.T) {
	p := profile.DefaultPath()
	dataDir := profile.DataDir()
	if !strings.HasPrefix(p, dataDir) {
		t.Fatalf("expected DefaultPath (%s) to be under DataDir (%s)", p, dataDir)
	}
	if filepath.Base(p) != "profiles.json" {
		t.Fatalf("expected filename 'profiles.json', got: %s", filepath.Base(p))
	}
}

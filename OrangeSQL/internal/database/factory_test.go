package database_test

import (
	"testing"

	"OrangeSQL/internal/database"
)

func TestNew_SQLite(t *testing.T) {
	db, err := database.New("sqlite", database.ConnectionParams{Path: ":memory:"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer db.Close()

	if db.Name() != ":memory:" {
		t.Fatalf("expected name ':memory:', got '%s'", db.Name())
	}
}

func TestNew_UnsupportedDriver(t *testing.T) {
	_, err := database.New("oracle", database.ConnectionParams{})
	if err == nil {
		t.Fatal("expected error for unsupported driver")
	}
}

func TestDefaultPort(t *testing.T) {
	tests := []struct {
		driver string
		want   int
	}{
		{"postgres", 5432},
		{"mysql", 3306},
		{"sqlserver", 1433},
		{"sqlite", 0},
		{"unknown", 0},
	}
	for _, tt := range tests {
		got := database.DefaultPort(tt.driver)
		if got != tt.want {
			t.Errorf("DefaultPort(%q) = %d, want %d", tt.driver, got, tt.want)
		}
	}
}

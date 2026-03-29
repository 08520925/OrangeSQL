package profile_test

import (
	"testing"

	"OrangeSQL/internal/profile"
)

func TestProfile_Validate_SQLite(t *testing.T) {
	p := profile.Profile{Name: "test", Driver: "sqlite", Path: "/tmp/test.db"}
	if err := p.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestProfile_Validate_SQLite_MissingPath(t *testing.T) {
	p := profile.Profile{Name: "test", Driver: "sqlite"}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestProfile_Validate_Postgres(t *testing.T) {
	p := profile.Profile{Name: "test", Driver: "postgres", Host: "localhost", DBName: "mydb"}
	if err := p.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestProfile_Validate_Postgres_MissingHost(t *testing.T) {
	p := profile.Profile{Name: "test", Driver: "postgres", DBName: "mydb"}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for missing host")
	}
}

func TestProfile_Validate_Postgres_MissingDBName(t *testing.T) {
	p := profile.Profile{Name: "test", Driver: "postgres", Host: "localhost"}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for missing dbName")
	}
}

func TestProfile_Validate_MySQL(t *testing.T) {
	p := profile.Profile{Name: "test", Driver: "mysql", Host: "localhost", DBName: "mydb"}
	if err := p.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestProfile_Validate_SQLServer(t *testing.T) {
	p := profile.Profile{Name: "test", Driver: "sqlserver", Host: "localhost", DBName: "mydb"}
	if err := p.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestProfile_Validate_MissingName(t *testing.T) {
	p := profile.Profile{Driver: "sqlite", Path: "/tmp/test.db"}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestProfile_Validate_UnsupportedDriver(t *testing.T) {
	p := profile.Profile{Name: "test", Driver: "oracle"}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for unsupported driver")
	}
}

func TestProfile_ToConnectionParams(t *testing.T) {
	p := profile.Profile{
		Name:     "test",
		Driver:   "postgres",
		Host:     "db.example.com",
		Port:     5433,
		User:     "admin",
		Password: "secret",
		DBName:   "testdb",
		SSLMode:  "require",
	}
	params := p.ToConnectionParams()
	if params.Host != "db.example.com" {
		t.Fatalf("expected host 'db.example.com', got '%s'", params.Host)
	}
	if params.Port != 5433 {
		t.Fatalf("expected port 5433, got %d", params.Port)
	}
	if params.DBName != "testdb" {
		t.Fatalf("expected dbName 'testdb', got '%s'", params.DBName)
	}
}

package profile

import (
	"fmt"
	"time"

	"OrangeSQL/internal/database"
)

// Profile は接続プロファイル。
type Profile struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Driver    string `json:"driver"`
	CreatedAt string `json:"createdAt"`

	// SQLite
	Path string `json:"path,omitempty"`

	// TCP 系共通（PostgreSQL / MySQL / SQL Server）
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	DBName   string `json:"dbName,omitempty"`
	SSLMode  string `json:"sslMode,omitempty"`
}

// ToConnectionParams は Profile を ConnectionParams に変換する。
func (p Profile) ToConnectionParams() database.ConnectionParams {
	return database.ConnectionParams{
		Path:     p.Path,
		Host:     p.Host,
		Port:     p.Port,
		User:     p.User,
		Password: p.Password,
		DBName:   p.DBName,
		SSLMode:  p.SSLMode,
	}
}

// Validate はドライバごとの必須フィールドを検証する。
func (p Profile) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	switch p.Driver {
	case "sqlite":
		if p.Path == "" {
			return fmt.Errorf("path is required for sqlite")
		}
	case "postgres", "mysql", "sqlserver":
		if p.Host == "" {
			return fmt.Errorf("host is required for %s", p.Driver)
		}
		if p.DBName == "" {
			return fmt.Errorf("dbName is required for %s", p.Driver)
		}
	default:
		return fmt.Errorf("unsupported driver: %s", p.Driver)
	}
	return nil
}

// ProfileStore は profiles.json の構造。
type ProfileStore struct {
	Profiles   []Profile `json:"profiles"`
	LastUsedID string    `json:"lastUsedId"`
}

// CreateRequest は POST /api/profiles のリクエスト。
type CreateRequest struct {
	Name     string `json:"name"`
	Driver   string `json:"driver"`
	Path     string `json:"path,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	DBName   string `json:"dbName,omitempty"`
	SSLMode  string `json:"sslMode,omitempty"`
}

// UpdateRequest は PUT /api/profiles/{id} のリクエスト。
type UpdateRequest struct {
	Name     string `json:"name,omitempty"`
	Path     string `json:"path,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	DBName   string `json:"dbName,omitempty"`
	SSLMode  string `json:"sslMode,omitempty"`
}

// ProfilesResponse は GET /api/profiles のレスポンス。
type ProfilesResponse struct {
	Profiles []Profile `json:"profiles"`
	ActiveID string    `json:"activeId"`
}

// ErrorResponse はエラー時の JSON レスポンス。
type ErrorResponse struct {
	Error string `json:"error"`
}

// NewProfile は新しい Profile を生成する。
func NewProfile(id string, req CreateRequest) Profile {
	driver := req.Driver
	if driver == "" {
		driver = "sqlite"
	}
	return Profile{
		ID:        id,
		Name:      req.Name,
		Driver:    driver,
		Path:      req.Path,
		Host:      req.Host,
		Port:      req.Port,
		User:      req.User,
		Password:  req.Password,
		DBName:    req.DBName,
		SSLMode:   req.SSLMode,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

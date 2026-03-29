package profile

import "time"

// Profile は接続プロファイル。
type Profile struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Driver    string `json:"driver"`
	Path      string `json:"path"`
	CreatedAt string `json:"createdAt"`
}

// ProfileStore は profiles.json の構造。
type ProfileStore struct {
	Profiles   []Profile `json:"profiles"`
	LastUsedID string    `json:"lastUsedId"`
}

// CreateRequest は POST /api/profiles のリクエスト。
type CreateRequest struct {
	Name   string `json:"name"`
	Driver string `json:"driver"`
	Path   string `json:"path"`
}

// UpdateRequest は PUT /api/profiles/{id} のリクエスト。
type UpdateRequest struct {
	Name string `json:"name"`
	Path string `json:"path"`
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
func NewProfile(id, name, driver, path string) Profile {
	return Profile{
		ID:        id,
		Name:      name,
		Driver:    driver,
		Path:      path,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

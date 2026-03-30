package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// Storage は profiles.json の読み書きを担当する。
type Storage struct {
	mu   sync.Mutex
	path string
}

// NewStorage は指定パスの profiles.json を管理する Storage を返す。
// ディレクトリが存在しない場合は自動作成する。
func NewStorage(path string) (*Storage, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create config dir: %w", err)
	}
	return &Storage{path: path}, nil
}

// DataDir は OS 標準のアプリケーションデータディレクトリを返す。
// Windows: %APPDATA%\OrangeSQL, macOS: ~/Library/Application Support/OrangeSQL, Linux: ~/.config/OrangeSQL
func DataDir() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		home, _ := os.UserHomeDir()
		dir = home
	}
	return filepath.Join(dir, "OrangeSQL")
}

// DefaultPath は OS に応じたデフォルトの profiles.json パスを返す。
func DefaultPath() string {
	return filepath.Join(DataDir(), "profiles.json")
}

// Load は profiles.json を読み込む。ファイルが存在しない場合は空の store を返す。
func (s *Storage) Load() (ProfileStore, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return ProfileStore{Profiles: []Profile{}}, nil
		}
		return ProfileStore{}, fmt.Errorf("read profiles: %w", err)
	}

	var store ProfileStore
	if err := json.Unmarshal(data, &store); err != nil {
		return ProfileStore{}, fmt.Errorf("parse profiles: %w", err)
	}
	if store.Profiles == nil {
		store.Profiles = []Profile{}
	}
	return store, nil
}

// Save は profiles.json に書き込む。
func (s *Storage) Save(store ProfileStore) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal profiles: %w", err)
	}
	return os.WriteFile(s.path, data, 0644)
}

// NextID は store 内の最大 ID + 1 を文字列で返す。
func NextID(store ProfileStore) string {
	max := 0
	for _, p := range store.Profiles {
		if n, err := strconv.Atoi(p.ID); err == nil && n > max {
			max = n
		}
	}
	return strconv.Itoa(max + 1)
}

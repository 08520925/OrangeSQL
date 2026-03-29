package profile

import (
	"fmt"
	"sync"

	"OrangeSQL/internal/database"
)

// ConnectionManager は接続プロファイルと現在の DB 接続を管理する。
type ConnectionManager struct {
	mu       sync.RWMutex
	current  database.Database
	activeID string
	storage  *Storage
}

// NewConnectionManager は Storage と初期 DB 接続から ConnectionManager を作成する。
func NewConnectionManager(storage *Storage, db database.Database, activeID string) *ConnectionManager {
	return &ConnectionManager{
		current:  db,
		activeID: activeID,
		storage:  storage,
	}
}

// DB は現在の DB 接続を返す。全ハンドラはこれ経由で DB にアクセスする。
func (m *ConnectionManager) DB() database.Database {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.current
}

// ActiveID は現在接続中のプロファイル ID を返す。
func (m *ConnectionManager) ActiveID() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.activeID
}

// Storage は Storage を返す。
func (m *ConnectionManager) Store() *Storage {
	return m.storage
}

// SwitchTo は指定プロファイルに接続を切り替える。
// 新 DB 接続に成功したら旧 DB を閉じて差し替える。失敗時は旧接続を維持する。
func (m *ConnectionManager) SwitchTo(p Profile) error {
	newDB, err := database.NewSQLite(p.Path)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	m.mu.Lock()
	oldDB := m.current
	m.current = newDB
	m.activeID = p.ID
	m.mu.Unlock()

	if oldDB != nil {
		oldDB.Close()
	}

	// lastUsedId を更新
	store, err := m.storage.Load()
	if err == nil {
		store.LastUsedID = p.ID
		m.storage.Save(store)
	}

	return nil
}

// Close は現在の DB 接続を閉じる。
func (m *ConnectionManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.current != nil {
		return m.current.Close()
	}
	return nil
}

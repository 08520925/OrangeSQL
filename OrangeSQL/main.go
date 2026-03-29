package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"OrangeSQL/internal/database"
	"OrangeSQL/internal/profile"
	"OrangeSQL/internal/server"
)

func main() {
	dbPath := flag.String("db", "./data.db", "SQLite database file path (used for initial profile)")
	flag.Parse()

	storage, err := profile.NewStorage(profile.DefaultPath())
	if err != nil {
		log.Fatalf("failed to init profile storage: %v", err)
	}

	store, err := storage.Load()
	if err != nil {
		log.Fatalf("failed to load profiles: %v", err)
	}

	// 初回起動: プロファイルが空なら -db フラグの値で初期プロファイルを作成
	if len(store.Profiles) == 0 {
		p := profile.NewProfile("1", "Default", "sqlite", *dbPath)
		store.Profiles = append(store.Profiles, p)
		store.LastUsedID = "1"
		if err := storage.Save(store); err != nil {
			log.Fatalf("failed to save initial profile: %v", err)
		}
	}

	// lastUsedId のプロファイルに接続
	var connectProfile profile.Profile
	found := false
	for _, p := range store.Profiles {
		if p.ID == store.LastUsedID {
			connectProfile = p
			found = true
			break
		}
	}
	if !found {
		connectProfile = store.Profiles[0]
	}

	db, err := database.NewSQLite(connectProfile.Path)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	cm := profile.NewConnectionManager(storage, db, connectProfile.ID)
	defer cm.Close()

	router := server.NewRouter(cm)

	addr := ":5522"
	fmt.Printf("OrangeSQL server starting on http://localhost%s\n", addr)
	fmt.Printf("Profile: %s (%s)\n", connectProfile.Name, connectProfile.Path)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

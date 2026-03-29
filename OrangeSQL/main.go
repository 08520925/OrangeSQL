package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os/exec"
	"runtime"

	"OrangeSQL/internal/database"
	"OrangeSQL/internal/profile"
	"OrangeSQL/internal/server"
)

//go:embed frontend/dist
var frontendFS embed.FS

func main() {
	dbPath := flag.String("db", "./data.db", "SQLite database file path (used for initial profile)")
	noBrowser := flag.Bool("no-browser", false, "disable auto-opening browser on startup")
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
		p := profile.NewProfile("1", profile.CreateRequest{
			Name:   "Default",
			Driver: "sqlite",
			Path:   *dbPath,
		})
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

	db, err := database.New(connectProfile.Driver, connectProfile.ToConnectionParams())
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	cm := profile.NewConnectionManager(storage, db, connectProfile.ID)
	defer cm.Close()

	// frontend/dist を fs.FS として取得
	distFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Fatalf("failed to load embedded frontend: %v", err)
	}

	router := server.NewRouter(cm, distFS)

	addr := ":5522"
	url := "http://localhost" + addr
	fmt.Printf("OrangeSQL server starting on %s\n", url)
	fmt.Printf("Profile: %s (%s)\n", connectProfile.Name, db.Name())

	// ブラウザ自動オープン
	if !*noBrowser {
		openBrowser(url)
	}

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// openBrowser はデフォルトブラウザで URL を開く。
func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	if err := cmd.Start(); err != nil {
		log.Printf("failed to open browser: %v", err)
	}
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"OrangeSQL/internal/database"
	"OrangeSQL/internal/server"
)

func main() {
	dbPath := flag.String("db", "./data.db", "SQLite database file path")
	flag.Parse()

	db, err := database.NewSQLite(*dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	router := server.NewRouter(db)

	addr := ":5522"
	fmt.Printf("OrangeSQL server starting on http://localhost%s\n", addr)
	fmt.Printf("Database: %s\n", db.Name())

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

package store

import (
	"database/sql"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// InitDatabase initializes the SQLite database and calls runMigrations
func InitDatabase(dbFilePath string, migrationsDir string) *sql.DB {
	absPath, _ := filepath.Abs(dbFilePath)
	log.Printf("Using SQLite database at: %s", absPath)

	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to SQLite database: %v", err)
	}

	log.Println("SQLite database connected successfully")
	runMigrations(db, migrationsDir)

	return db
}

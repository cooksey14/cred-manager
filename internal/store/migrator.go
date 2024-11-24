package store

import (
	"database/sql"
	"log"

	"github.com/pressly/goose/v3"
)

// runMigrations applies database migrations on startup
func runMigrations(db *sql.DB, migrationsDir string) {
	goose.SetDialect("sqlite3")

	if err := goose.Up(db, migrationsDir); err != nil {
		log.Fatalf("Migration failed: %v", err)
	} else {
		log.Println("Migrations applied successfully")
	}
}

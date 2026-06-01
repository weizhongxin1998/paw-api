package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
	appDataDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir: %w", err)
	}
	dbDir := filepath.Join(appDataDir, ".paw-api")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create db dir: %w", err)
	}
	dbPath := filepath.Join(dbDir, "paw-api.db")
	DB, err = sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	DB.SetMaxOpenConns(1)
	if err := runMigrations(DB); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

func OpenDatabase(dbPath string) (*sql.DB, string, error) {
	if dbPath != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
			return nil, "", fmt.Errorf("create database directory: %w", err)
		}
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, "", fmt.Errorf("open sqlite database: %w", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		db.Close()
		return nil, "", fmt.Errorf("enable foreign keys: %w", err)
	}

	return db, dbPath, nil
}

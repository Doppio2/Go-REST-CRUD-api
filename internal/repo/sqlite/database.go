package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

const defaultDatabasePath = "db/go_rest_crud.db"

func ResolveDatabasePath() string {
	if path := os.Getenv("SQLITE_PATH"); path != "" {
		return path
	}

	return defaultDatabasePath
}

func OpenDatabase() (*sql.DB, string, error) {
	dbPath := ResolveDatabasePath()
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

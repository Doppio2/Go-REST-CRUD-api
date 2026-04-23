package sqlite

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

const migrationsDir = "migrations"

//go:embed migrations/*.sql
var migrationsFS embed.FS

func Migrate(db *sql.DB) error {
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("sqlite"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}

func RunMigrationCommand(db *sql.DB, command string, args ...string) error {
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("sqlite"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	if err := goose.Run(command, db, migrationsDir, args...); err != nil {
		return fmt.Errorf("run goose command %q: %w", command, err)
	}

	return nil
}

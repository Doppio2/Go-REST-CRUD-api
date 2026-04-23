package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	defaultAppPort    = "8080"
	defaultSQLitePath = "db/go_rest_crud.db"
)

type Config struct {
	AppPort     string
	SQLitePath  string
	AutoMigrate bool
}

func Load() (Config, error) {
	autoMigrate, err := parseBoolEnv("AUTO_MIGRATE", true)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		AppPort:     getEnv("APP_PORT", defaultAppPort),
		SQLitePath:  getEnv("SQLITE_PATH", defaultSQLitePath),
		AutoMigrate: autoMigrate,
	}

	return cfg, nil
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func parseBoolEnv(key string, fallback bool) (bool, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%s must be a boolean value: %w", key, err)
	}

	return parsed, nil
}

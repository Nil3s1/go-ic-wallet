package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

const (
	defaultMaxOpenConns    = 25
	defaultMaxIdleConns    = 25
	defaultConnMaxLifetime = 5 * time.Minute
	defaultConnMaxIdleTime = 1 * time.Minute
)

type Config struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// NewMSSQL builds the connection pool and verifies connectivity once at startup.
func NewMSSQL(ctx context.Context, cfg Config) (*sql.DB, error) {
	if cfg.DSN == "" {
		return nil, fmt.Errorf("database: dsn is empty")
	}

	db, err := sql.Open("sqlserver", cfg.DSN)
	if err != nil {
		// Never wrap the DSN into the error, it carries credentials.
		return nil, fmt.Errorf("database: open mssql: %w", err)
	}

	db.SetMaxOpenConns(orDefaultInt(cfg.MaxOpenConns, defaultMaxOpenConns))
	db.SetMaxIdleConns(orDefaultInt(cfg.MaxIdleConns, defaultMaxIdleConns))
	db.SetConnMaxLifetime(orDefaultDuration(cfg.ConnMaxLifetime, defaultConnMaxLifetime))
	db.SetConnMaxIdleTime(orDefaultDuration(cfg.ConnMaxIdleTime, defaultConnMaxIdleTime))

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("database: ping mssql: %w", err)
	}

	return db, nil
}

func orDefaultInt(value, fallback int) int {
	if value > 0 {
		return value
	}
	return fallback
}

func orDefaultDuration(value, fallback time.Duration) time.Duration {
	if value > 0 {
		return value
	}
	return fallback
}

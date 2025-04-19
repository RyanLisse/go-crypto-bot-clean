// +build nosqlite

package factory

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/rs/zerolog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// NewDBConnection creates a new database connection based on the configuration
func NewDBConnection(cfg *config.Config, log *zerolog.Logger) (*gorm.DB, error) {
	// Create a GORM logger that uses our zerolog instance
	gLogger := NewGormLogger(log, cfg.Database.EnableLogging)

	// Create GORM config
	gormConfig := &gorm.Config{
		Logger: gLogger,
	}

	// Connect to the database based on the configuration
	var db *gorm.DB
	var err error

	switch cfg.Database.Type {
	case "mysql":
		db, err = gorm.Open(mysql.Open(cfg.Database.DSN), gormConfig)
	case "postgres":
		// Implement PostgreSQL connection if needed
		return nil, fmt.Errorf("postgres database not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetimeMinutes) * time.Minute)

	log.Info().
		Str("type", cfg.Database.Type).
		Str("dsn", maskDSN(cfg.Database.DSN)).
		Int("max_idle_conns", cfg.Database.MaxIdleConns).
		Int("max_open_conns", cfg.Database.MaxOpenConns).
		Int("conn_max_lifetime_minutes", cfg.Database.ConnMaxLifetimeMinutes).
		Msg("Connected to database")

	return db, nil
}

// NewGormLogger creates a new GORM logger that uses zerolog
func NewGormLogger(log *zerolog.Logger, enableLogging bool) gormLogger.Interface {
	var logLevel gormLogger.LogLevel
	if enableLogging {
		logLevel = gormLogger.Info
	} else {
		logLevel = gormLogger.Silent
	}

	return gormLogger.New(
		&zerologWriter{log: log},
		gormLogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}

// zerologWriter implements the io.Writer interface for zerolog
type zerologWriter struct {
	log *zerolog.Logger
}

// Write implements the io.Writer interface
func (w *zerologWriter) Write(p []byte) (n int, err error) {
	w.log.Debug().Msg(string(p))
	return len(p), nil
}

// maskDSN masks sensitive information in the DSN for logging
func maskDSN(dsn string) string {
	// For SQLite, just return the path
	return dsn
}

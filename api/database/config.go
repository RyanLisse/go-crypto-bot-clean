// Package database provides database functionality for the API
package database

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds database configuration
type Config struct {
	Path            string
	Debug           bool
	Logger          *zap.Logger
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

// DefaultConfig returns the default database configuration
func DefaultConfig() Config {
	return Config{
		Path:            "data/api.db",
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Hour,
	}
}

// NewDatabase creates a new GORM database connection
func NewDatabase(config Config) (*gorm.DB, error) {
	// Ensure directory exists
	dbDir := filepath.Dir(config.Path)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Configure GORM logger
	gormLogLevel := logger.Silent
	if config.Debug {
		gormLogLevel = logger.Info
	}

	gormLogger := logger.New(
		logger.Writer{}, // Use default writer
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gormLogLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Open database connection
	db, err := gorm.Open(sqlite.Open(config.Path), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	maxIdle := 10
	if config.MaxIdleConns > 0 {
		maxIdle = config.MaxIdleConns
	}
	sqlDB.SetMaxIdleConns(maxIdle)

	maxOpen := 100
	if config.MaxOpenConns > 0 {
		maxOpen = config.MaxOpenConns
	}
	sqlDB.SetMaxOpenConns(maxOpen)

	lifetime := time.Hour
	if config.ConnMaxLifetime > 0 {
		lifetime = config.ConnMaxLifetime
	}
	sqlDB.SetConnMaxLifetime(lifetime)

	// Ping to verify connection
	if err = sqlDB.Ping(); err != nil {
		sqlDB.Close() // Close if ping fails
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Get logger if not provided
	zapLogger := config.Logger
	if zapLogger == nil {
		var loggerErr error
		zapLogger, loggerErr = zap.NewProduction()
		if loggerErr != nil {
			sqlDB.Close()
			return nil, fmt.Errorf("failed to create logger: %w", loggerErr)
		}
	}
	zapLogger.Info("Database connection established successfully", zap.String("path", config.Path))

	return db, nil
}

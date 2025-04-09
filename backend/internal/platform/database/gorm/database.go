package gorm

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database represents a GORM database connection
type Database struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

// Config holds database configuration
type Config struct {
	Path   string
	Debug  bool
	Logger *zap.Logger
}

// NewDatabase creates a new GORM database connection
func NewDatabase(config Config) (*Database, error) {
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
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel: gormLogLevel,
			Colorful: true,
		},
	)

	// Open database connection
	db, err := gorm.Open(sqlite.Open(config.Path), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set logger
	zapLogger := config.Logger
	if zapLogger == nil {
		var loggerErr error
		zapLogger, loggerErr = zap.NewProduction()
		if loggerErr != nil {
			return nil, fmt.Errorf("failed to create logger: %w", loggerErr)
		}
	}

	return &Database{
		DB:     db,
		Logger: zapLogger,
	}, nil
}

// Migrate runs auto-migrations for all models
func (d *Database) Migrate() error {
	d.Logger.Info("Running database migrations")

	// Add all models to be migrated here
	err := d.DB.AutoMigrate(
		&models.BoughtCoin{},
		&models.NewCoin{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	d.Logger.Info("Database migrations completed successfully")
	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB: %w", err)
	}
	return sqlDB.Close()
}

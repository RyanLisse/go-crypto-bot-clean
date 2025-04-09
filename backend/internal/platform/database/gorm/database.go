package gorm

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
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

// NewDatabase creates a new GORM database connection and returns the *gorm.DB instance.
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

	// Use a more sophisticated GORM logger if needed, perhaps integrating with zap
	gormLogger := logger.New(
		log.New(os.Stdout, "\\r\\n", log.LstdFlags), // Simple logger for now
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Log slow queries
			LogLevel:                  gormLogLevel,
			IgnoreRecordNotFoundError: true, // Don't log ErrRecordNotFound
			Colorful:                  true,
		},
	)

	// Open database connection
	db, err := gorm.Open(sqlite.Open(config.Path), &gorm.Config{
		Logger: gormLogger,
		// Add other GORM configs if needed, e.g., NamingStrategy
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
		// Consider using NewDevelopment() if config.Debug is true
		zapLogger, loggerErr = zap.NewProduction()
		if loggerErr != nil {
			sqlDB.Close()
			return nil, fmt.Errorf("failed to create logger: %w", loggerErr)
		}
	}
	zapLogger.Info("GORM Database connection established successfully", zap.String("path", config.Path))

	return db, nil
}

// RunMigrations runs auto-migrations for all models using the provided GORM DB instance.
func RunMigrations(db *gorm.DB, logger *zap.Logger) error {
	if logger == nil {
		// Fallback logger if none provided
		var loggerErr error
		logger, loggerErr = zap.NewProduction()
		if loggerErr != nil {
			return fmt.Errorf("failed to create fallback logger for migrations: %w", loggerErr)
		}
	}

	logger.Info("Running GORM AutoMigrate")

	// Add all models to be migrated here
	err := db.AutoMigrate(
		&models.BoughtCoin{},
		&models.NewCoin{},
		&models.BalanceHistory{},
		// Add other models as needed:
		// &models.PurchaseDecision{},
		// &models.LogEvent{},
	)
	if err != nil {
		logger.Error("GORM AutoMigrate failed", zap.Error(err))
		return fmt.Errorf("failed to run GORM AutoMigrate: %w", err)
	}

	logger.Info("GORM AutoMigrate completed successfully")
	return nil
}

// CloseDatabase closes the GORM database connection
func CloseDatabase(db *gorm.DB) error {
	if db == nil {
		return nil // Nothing to close
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB for closing: %w", err)
	}
	return sqlDB.Close()
}

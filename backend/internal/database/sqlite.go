// Package database provides database connectivity and operations
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

// SQLiteConfig represents the configuration for SQLite database
type SQLiteConfig struct {
	Path                   string
	MaxOpenConns           int
	MaxIdleConns           int
	ConnMaxLifetimeSeconds int
	Debug                  bool
}

// SQLiteManager manages the SQLite database connection
type SQLiteManager struct {
	config SQLiteConfig
	db     *gorm.DB
	logger *zap.Logger
}

// NewSQLiteManager creates a new SQLite database manager
func NewSQLiteManager(config SQLiteConfig, logger *zap.Logger) *SQLiteManager {
	return &SQLiteManager{
		config: config,
		logger: logger,
	}
}

// Connect establishes a connection to the SQLite database
func (m *SQLiteManager) Connect() error {
	// Ensure the directory exists
	dbDir := filepath.Dir(m.config.Path)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Configure GORM logger
	gormLogger := logger.New(
		&zapAdapter{logger: m.logger},
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  getGormLogLevel(m.config.Debug),
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Open database connection
	db, err := gorm.Open(sqlite.Open(m.config.Path), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	sqlDB.SetMaxOpenConns(m.config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(m.config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(m.config.ConnMaxLifetimeSeconds) * time.Second)

	m.db = db
	m.logger.Info("Connected to SQLite database",
		zap.String("path", m.config.Path),
		zap.Int("maxOpenConns", m.config.MaxOpenConns),
		zap.Int("maxIdleConns", m.config.MaxIdleConns),
		zap.Int("connMaxLifetimeSeconds", m.config.ConnMaxLifetimeSeconds),
		zap.Bool("debug", m.config.Debug),
	)

	return nil
}

// DB returns the GORM database instance
func (m *SQLiteManager) DB() *gorm.DB {
	return m.db
}

// Close closes the database connection
func (m *SQLiteManager) Close() error {
	if m.db == nil {
		return nil
	}

	sqlDB, err := m.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	m.logger.Info("Closed SQLite database connection", zap.String("path", m.config.Path))
	return nil
}

// AutoMigrate runs auto migration for the given models
func (m *SQLiteManager) AutoMigrate(models ...interface{}) error {
	if m.db == nil {
		return fmt.Errorf("database not connected")
	}

	if err := m.db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to run auto migration: %w", err)
	}

	m.logger.Info("Successfully ran auto migration")
	return nil
}

// getGormLogLevel converts debug flag to GORM log level
func getGormLogLevel(debug bool) logger.LogLevel {
	if debug {
		return logger.Info
	}
	return logger.Error
}

// zapAdapter adapts zap.Logger to GORM's logger interface
type zapAdapter struct {
	logger *zap.Logger
}

// Printf implements GORM's logger interface
func (a *zapAdapter) Printf(format string, args ...interface{}) {
	a.logger.Sugar().Debugf(format, args...)
}

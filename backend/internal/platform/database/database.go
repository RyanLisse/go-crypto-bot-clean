package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// DB is the global database connection
var DB *sqlx.DB

// Initialize sets up the database connection
func Initialize(dsn string) error {
	// Create directories if they don't exist
	dir := filepath.Dir(dsn)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// Run migrations
	if err := RunMigrations(db); err != nil {
		db.Close()
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	DB = db
	fmt.Println("Database initialized successfully")
	return nil
}

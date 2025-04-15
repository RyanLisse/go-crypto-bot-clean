package scripts

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// DEPRECATED: This script is no longer used. All database migrations are now handled via GORM AutoMigrate.
// See docs/database-migrations.md for details.

// RunMigrations applies all pending database migrations
func RunMigrations() {
	// Database path
	dbPath := "../data/crypto_bot.db"

	// Ensure the database directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	// Connect to the database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create migrations table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	// Get list of applied migrations
	rows, err := db.Query("SELECT name FROM migrations")
	if err != nil {
		log.Fatalf("Failed to query migrations: %v", err)
	}
	defer rows.Close()

	appliedMigrations := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatalf("Failed to scan migration row: %v", err)
		}
		appliedMigrations[name] = true
	}

	// Get list of migration files
	migrationsDir := "../migrations"
	// Make sure the migrations directory exists
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		log.Fatalf("Migrations directory does not exist: %s", migrationsDir)
	}

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	// Filter and sort migration files
	var migrations []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") && strings.Contains(file.Name(), "_sqlite") {
			migrations = append(migrations, file.Name())
		}
	}
	sort.Strings(migrations)

	// Apply migrations
	for _, migration := range migrations {
		migrationName := strings.TrimSuffix(migration, ".up.sql")
		if appliedMigrations[migrationName] {
			fmt.Printf("Migration %s already applied, skipping\n", migrationName)
			continue
		}

		fmt.Printf("Applying migration %s...\n", migrationName)
		migrationPath := filepath.Join(migrationsDir, migration)
		migrationSQL, err := os.ReadFile(migrationPath)
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", migrationPath, err)
		}

		// Begin transaction
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("Failed to begin transaction: %v", err)
		}

		// Execute migration
		_, err = tx.Exec(string(migrationSQL))
		if err != nil {
			tx.Rollback()
			log.Fatalf("Failed to execute migration %s: %v", migrationName, err)
		}

		// Record migration
		_, err = tx.Exec("INSERT INTO migrations (name) VALUES (?)", migrationName)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Failed to record migration %s: %v", migrationName, err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			log.Fatalf("Failed to commit transaction: %v", err)
		}

		fmt.Printf("Migration %s applied successfully\n", migrationName)
	}

	fmt.Println("All migrations applied successfully")
}

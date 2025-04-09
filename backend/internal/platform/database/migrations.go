package database

import (
	"embed"
	"fmt"

	"github.com/jmoiron/sqlx"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// RunMigrations executes all database migrations in order
func RunMigrations(db *sqlx.DB) error {
	// Create migrations table if it doesn't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// List of migrations to apply (in order)
	migrations := []string{
		"01_create_bought_coins.sql",
		"02_create_new_coins.sql",
		"03_create_purchase_decisions.sql",
		"04_create_log_events.sql",
		"05_update_new_coins.sql",
		"06_add_status_to_new_coins.sql",
	}

	// Check and apply each migration
	for _, migration := range migrations {
		var count int
		err = db.Get(&count, "SELECT COUNT(*) FROM migrations WHERE name = ?", migration)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if count == 0 {
			// Migration hasn't been applied yet
			sqlBytes, err := migrationFiles.ReadFile("migrations/" + migration)
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", migration, err)
			}

			// Execute migration within a transaction
			tx, err := db.Beginx()
			if err != nil {
				return fmt.Errorf("failed to start transaction: %w", err)
			}

			_, err = tx.Exec(string(sqlBytes))
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to apply migration %s: %w", migration, err)
			}

			// Record the migration
			_, err = tx.Exec("INSERT INTO migrations (name) VALUES (?)", migration)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to record migration %s: %w", migration, err)
			}

			err = tx.Commit()
			if err != nil {
				return fmt.Errorf("failed to commit migration %s: %w", migration, err)
			}

			fmt.Printf("Applied migration: %s\n", migration)
		}
	}

	return nil
}

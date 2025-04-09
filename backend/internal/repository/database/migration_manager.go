package database

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// MigrationManager handles database migrations
type MigrationManager struct {
	repo Repository
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(repo Repository) *MigrationManager {
	return &MigrationManager{
		repo: repo,
	}
}

// Migration represents a database migration
type Migration struct {
	Version     string
	Description string
	SQL         string
}

// AppliedMigration represents a migration that has been applied to the database
type AppliedMigration struct {
	ID          int64
	Version     string
	AppliedAt   time.Time
	Description string
}

// MigrateDatabase applies all pending migrations to the database
func (m *MigrationManager) MigrateDatabase(ctx context.Context) error {
	// Ensure migration_history table exists
	err := m.ensureMigrationTable(ctx)
	if err != nil {
		return fmt.Errorf("failed to ensure migration table exists: %w", err)
	}

	// Get applied migrations
	appliedMigrations, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Get available migrations
	availableMigrations, err := m.getAvailableMigrations()
	if err != nil {
		return fmt.Errorf("failed to get available migrations: %w", err)
	}

	// Apply pending migrations
	for _, migration := range availableMigrations {
		// Check if migration has already been applied
		if m.isMigrationApplied(migration.Version, appliedMigrations) {
			log.Printf("Migration %s already applied, skipping", migration.Version)
			continue
		}

		log.Printf("Applying migration %s: %s", migration.Version, migration.Description)

		// Begin transaction
		tx, err := m.repo.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %s: %w", migration.Version, err)
		}

		// Apply migration
		_, err = tx.Exec(migration.SQL)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
		}

		// Record migration
		_, err = tx.Exec(
			"INSERT INTO migration_history (version, description) VALUES (?, ?)",
			migration.Version,
			migration.Description,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", migration.Version, err)
		}

		// Commit transaction
		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", migration.Version, err)
		}

		log.Printf("Migration %s applied successfully", migration.Version)
	}

	return nil
}

// ensureMigrationTable ensures the migration_history table exists
func (m *MigrationManager) ensureMigrationTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS migration_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			version TEXT NOT NULL,
			applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			description TEXT NOT NULL
		)
	`

	_, err := m.repo.Execute(ctx, query)
	return err
}

// getAppliedMigrations gets all migrations that have been applied to the database
func (m *MigrationManager) getAppliedMigrations(ctx context.Context) ([]AppliedMigration, error) {
	query := `
		SELECT id, version, applied_at, description
		FROM migration_history
		ORDER BY id ASC
	`

	rows, err := m.repo.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []AppliedMigration
	for rows.Next() {
		var migration AppliedMigration
		err := rows.Scan(
			&migration.ID,
			&migration.Version,
			&migration.AppliedAt,
			&migration.Description,
		)
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, migration)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return migrations, nil
}

// getAvailableMigrations gets all available migrations from the embedded filesystem
func (m *MigrationManager) getAvailableMigrations() ([]Migration, error) {
	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return nil, err
	}

	var migrations []Migration
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		// Parse version and description from filename
		// Expected format: 001_description.sql
		name := strings.TrimSuffix(entry.Name(), ".sql")
		parts := strings.SplitN(name, "_", 2)
		if len(parts) != 2 {
			log.Printf("Skipping migration with invalid name format: %s", entry.Name())
			continue
		}

		version := parts[0]
		description := strings.ReplaceAll(parts[1], "_", " ")

		// Read migration SQL
		content, err := fs.ReadFile(migrationsFS, filepath.Join("migrations", entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read migration %s: %w", entry.Name(), err)
		}

		migrations = append(migrations, Migration{
			Version:     version,
			Description: description,
			SQL:         string(content),
		})
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// isMigrationApplied checks if a migration has already been applied
func (m *MigrationManager) isMigrationApplied(version string, appliedMigrations []AppliedMigration) bool {
	for _, migration := range appliedMigrations {
		if migration.Version == version {
			return true
		}
	}
	return false
}

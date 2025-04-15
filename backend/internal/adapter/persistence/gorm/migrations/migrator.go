package migrations

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// MigrationFunc is a function that performs a database migration
type MigrationFunc func(*gorm.DB) error

// Migration represents a database migration
type Migration struct {
	Name     string
	Function MigrationFunc
}

// ObjectMigration defines the interface for object-based migrations
type ObjectMigration interface {
	Name() string
	Up(ctx context.Context, db *gorm.DB) error
	Down(ctx context.Context, db *gorm.DB) error
}

// Migrator manages database migrations
type Migrator struct {
	db               *gorm.DB
	logger           *zerolog.Logger
	migrations       []Migration
	objectMigrations []ObjectMigration
}

// NewMigrator creates a new Migrator
func NewMigrator(db *gorm.DB, logger *zerolog.Logger) *Migrator {
	return &Migrator{
		db:               db,
		logger:           logger,
		migrations:       []Migration{},
		objectMigrations: []ObjectMigration{},
	}
}

// AddMigration adds a migration to the migrator
func (m *Migrator) AddMigration(name string, function MigrationFunc) {
	m.migrations = append(m.migrations, Migration{
		Name:     name,
		Function: function,
	})
}

// AddObjectMigration adds an object-based migration to the migrator
func (m *Migrator) AddObjectMigration(migration ObjectMigration) {
	m.objectMigrations = append(m.objectMigrations, migration)
}

// MigrationRecord represents a record of a migration that has been applied
type MigrationRecord struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"uniqueIndex"`
	AppliedAt time.Time `gorm:"autoCreateTime"`
}

// TableName sets the table name for MigrationRecord
func (MigrationRecord) TableName() string {
	return "migrations"
}

// RunMigrations runs all pending migrations
func (m *Migrator) RunMigrations() error {
	// Create migrations table if it doesn't exist
	if err := m.db.AutoMigrate(&MigrationRecord{}); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	var appliedMigrations []MigrationRecord
	if err := m.db.Find(&appliedMigrations).Error; err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Create a map of applied migrations for quick lookup
	appliedMap := make(map[string]bool)
	for _, migration := range appliedMigrations {
		appliedMap[migration.Name] = true
	}

	// Run pending migrations
	for _, migration := range m.migrations {
		if !appliedMap[migration.Name] {
			m.logger.Info().Str("migration", migration.Name).Msg("Running migration")

			// Run migration in a transaction
			err := m.db.Transaction(func(tx *gorm.DB) error {
				// Run the migration
				if err := migration.Function(tx); err != nil {
					m.logger.Error().Err(err).Str("migration", migration.Name).Msg("Migration failed")
					return err
				}

				// Record the migration
				record := MigrationRecord{
					Name: migration.Name,
				}
				if err := tx.Create(&record).Error; err != nil {
					m.logger.Error().Err(err).Str("migration", migration.Name).Msg("Failed to record migration")
					return err
				}

				return nil
			})

			if err != nil {
				return fmt.Errorf("migration %s failed: %w", migration.Name, err)
			}

			m.logger.Info().Str("migration", migration.Name).Msg("Migration completed successfully")
		} else {
			m.logger.Debug().Str("migration", migration.Name).Msg("Migration already applied")
		}
	}

	// Run object-based migrations
	ctx := context.Background()
	for _, migration := range m.objectMigrations {
		migrationName := migration.Name()
		if !appliedMap[migrationName] {
			m.logger.Info().Str("migration", migrationName).Msg("Running object-based migration")

			// Run migration in a transaction
			err := m.db.Transaction(func(tx *gorm.DB) error {
				// Run the migration
				if err := migration.Up(ctx, tx); err != nil {
					m.logger.Error().Err(err).Str("migration", migrationName).Msg("Object-based migration failed")
					return err
				}

				// Record the migration
				record := MigrationRecord{
					Name: migrationName,
				}
				if err := tx.Create(&record).Error; err != nil {
					m.logger.Error().Err(err).Str("migration", migrationName).Msg("Failed to record object-based migration")
					return err
				}

				return nil
			})

			if err != nil {
				return fmt.Errorf("object-based migration %s failed: %w", migrationName, err)
			}

			m.logger.Info().Str("migration", migrationName).Msg("Object-based migration completed successfully")
		} else {
			m.logger.Debug().Str("migration", migrationName).Msg("Object-based migration already applied")
		}
	}

	return nil
}

package migrations

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// RunAllMigrations runs all defined database migrations in sequence.
func RunAllMigrations(db *gorm.DB, logger *zerolog.Logger) error {
	log := logger.With().Str("component", "migrations").Logger()
	log.Info().Msg("Starting database migrations...")

	// Add migration functions to this list in the desired execution order
	migrationFuncs := []func(*gorm.DB, *zerolog.Logger) error{
		CreateWalletsTable, // Add the wallet table migration
		// Add other migration functions here, e.g.:
		// CreateOrdersTable,
		// CreateUsersTable,
	}

	for _, migrateFunc := range migrationFuncs {
		if err := migrateFunc(db, logger); err != nil {
			// Log the specific error within the migration function is usually sufficient
			log.Error().Err(err).Msg("Migration failed")
			return err // Stop migrations if one fails
		}
	}

	log.Info().Msg("Database migrations completed successfully.")
	return nil
}

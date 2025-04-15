package migrations

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/repo"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// CreateStatusTable creates the status table
func CreateStatusTable(db *gorm.DB) error {
	logger := log.With().Str("migration", "create_status_table").Logger()
	logger.Info().Msg("Running migration: Create status table")

	// Create the status table
	if err := db.AutoMigrate(&repo.StatusRecord{}); err != nil {
		logger.Error().Err(err).Msg("Failed to create status table")
		return err
	}

	logger.Info().Msg("Status table created successfully")
	return nil
}

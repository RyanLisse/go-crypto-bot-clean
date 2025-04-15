package migrations

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func CreateMexcApiCredentialsTable(db *gorm.DB) error {
	logger := log.With().Str("migration", "create_mexc_api_credentials_table").Logger()
	logger.Info().Msg("Running migration: Create MEXC API credentials table")
	if err := db.AutoMigrate(&entity.MexcApiCredential{}); err != nil {
		logger.Error().Err(err).Msg("Failed to create MEXC API credentials table")
		return err
	}
	logger.Info().Msg("MEXC API credentials table created successfully")
	return nil
}

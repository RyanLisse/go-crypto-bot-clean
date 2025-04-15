package migrations

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func CreateAccountTable(db *gorm.DB) error {
	logger := log.With().Str("migration", "create_account_table").Logger()
	logger.Info().Msg("Running migration: Create account table")
	if err := db.AutoMigrate(&entity.AccountEntity{}); err != nil {
		logger.Error().Err(err).Msg("Failed to create account table")
		return err
	}
	logger.Info().Msg("Account table created successfully")
	return nil
}

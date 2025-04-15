package migrations

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func CreateWalletsTable(db *gorm.DB) error {
	logger := log.With().Str("migration", "create_wallets_table").Logger()
	logger.Info().Msg("Running migration: Create wallets table")
	if err := db.AutoMigrate(&entity.WalletEntity{}); err != nil {
		logger.Error().Err(err).Msg("Failed to create wallets table")
		return err
	}
	logger.Info().Msg("Wallets table created successfully")
	return nil
}

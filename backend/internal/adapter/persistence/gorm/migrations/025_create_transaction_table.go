package migrations

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func CreateTransactionsTable(db *gorm.DB) error {
	logger := log.With().Str("migration", "create_transactions_table").Logger()
	logger.Info().Msg("Running migration: Create transactions table")
	if err := db.AutoMigrate(&entity.TransactionEntity{}); err != nil {
		logger.Error().Err(err).Msg("Failed to create transactions table")
		return err
	}
	logger.Info().Msg("Transactions table created successfully")
	return nil
}

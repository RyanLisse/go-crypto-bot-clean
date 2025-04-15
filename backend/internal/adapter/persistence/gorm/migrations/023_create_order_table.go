package migrations

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func CreateOrdersTable(db *gorm.DB) error {
	logger := log.With().Str("migration", "create_orders_table").Logger()
	logger.Info().Msg("Running migration: Create orders table")
	if err := db.AutoMigrate(&entity.OrderEntity{}); err != nil {
		logger.Error().Err(err).Msg("Failed to create orders table")
		return err
	}
	logger.Info().Msg("Orders table created successfully")
	return nil
}

package migrations

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// AddOrderExtraFields adds additional fields to the orders table
func AddOrderExtraFields(db *gorm.DB) error {
	logger := log.With().Str("migration", "add_order_extra_fields").Logger()
	logger.Info().Msg("Running migration: Add order extra fields")

	// Add columns to the orders table
	if err := db.Exec(`
		ALTER TABLE orders 
		ADD COLUMN IF NOT EXISTS time_in_force VARCHAR(20),
		ADD COLUMN IF NOT EXISTS stop_price DECIMAL(24,8),
		ADD COLUMN IF NOT EXISTS iceberg_qty DECIMAL(24,8),
		ADD COLUMN IF NOT EXISTS order_list_id VARCHAR(100),
		ADD COLUMN IF NOT EXISTS executed_price DECIMAL(24,8),
		ADD COLUMN IF NOT EXISTS executed_qty DECIMAL(24,8),
		ADD COLUMN IF NOT EXISTS executed_quote_qty DECIMAL(24,8),
		ADD COLUMN IF NOT EXISTS executed_at TIMESTAMP,
		ADD COLUMN IF NOT EXISTS canceled_at TIMESTAMP,
		ADD COLUMN IF NOT EXISTS rejected_reason TEXT
	`).Error; err != nil {
		logger.Error().Err(err).Msg("Failed to add extra fields to orders table")
		return err
	}

	logger.Info().Msg("Order extra fields added successfully")
	return nil
}

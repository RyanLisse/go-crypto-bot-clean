package migrations

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// AddPositionOrderPrice adds order price columns to the positions table
func AddPositionOrderPrice(db *gorm.DB) error {
	logger := log.With().Str("migration", "add_position_order_price").Logger()
	logger.Info().Msg("Running migration: Add position order price")

	// Add columns to the positions table
	if err := db.Exec(`
		ALTER TABLE positions 
		ADD COLUMN IF NOT EXISTS exit_price DECIMAL(24,8),
		ADD COLUMN IF NOT EXISTS realized_pnl DECIMAL(24,8)
	`).Error; err != nil {
		logger.Error().Err(err).Msg("Failed to add order price columns to positions table")
		return err
	}

	logger.Info().Msg("Position order price columns added successfully")
	return nil
}

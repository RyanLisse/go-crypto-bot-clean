package migrations

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// AddPositionOrderId adds order ID columns to the positions table
func AddPositionOrderId(db *gorm.DB) error {
	logger := log.With().Str("migration", "add_position_order_id").Logger()
	logger.Info().Msg("Running migration: Add position order ID")

	// Add columns to the positions table
	if err := db.Exec(`
		ALTER TABLE positions 
		ADD COLUMN IF NOT EXISTS entry_order_id VARCHAR(100),
		ADD COLUMN IF NOT EXISTS exit_order_id VARCHAR(100)
	`).Error; err != nil {
		logger.Error().Err(err).Msg("Failed to add order ID columns to positions table")
		return err
	}

	logger.Info().Msg("Position order ID columns added successfully")
	return nil
}

package migrations

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// AddOrderSymbolIndex adds an index on the symbol column in the orders table
func AddOrderSymbolIndex(db *gorm.DB) error {
	logger := log.With().Str("migration", "add_order_symbol_index").Logger()
	logger.Info().Msg("Running migration: Add order symbol index")

	// Add index to the orders table
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_order_symbol ON orders(symbol)
	`).Error; err != nil {
		logger.Error().Err(err).Msg("Failed to add symbol index to orders table")
		return err
	}

	logger.Info().Msg("Order symbol index added successfully")
	return nil
}

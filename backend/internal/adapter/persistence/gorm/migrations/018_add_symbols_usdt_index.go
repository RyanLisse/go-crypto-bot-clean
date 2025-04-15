package migrations

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// AddSymbolsUsdtIndex adds an index for USDT symbols
func AddSymbolsUsdtIndex(db *gorm.DB) error {
	logger := log.With().Str("migration", "add_symbols_usdt_index").Logger()
	logger.Info().Msg("Running migration: Add symbols USDT index")

	// Add index to the symbols table
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_symbols_usdt ON symbols(quote_asset) WHERE quote_asset = 'USDT'
	`).Error; err != nil {
		logger.Error().Err(err).Msg("Failed to add USDT index to symbols table")
		return err
	}

	logger.Info().Msg("Symbols USDT index added successfully")
	return nil
}

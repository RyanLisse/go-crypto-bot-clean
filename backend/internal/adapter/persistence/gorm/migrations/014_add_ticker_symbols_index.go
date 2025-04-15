package migrations

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// AddTickerSymbolsIndex adds an index on the symbol column in the tickers table
func AddTickerSymbolsIndex(db *gorm.DB) error {
	logger := log.With().Str("migration", "add_ticker_symbols_index").Logger()
	logger.Info().Msg("Running migration: Add ticker symbols index")

	// Add index to the tickers table
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_ticker_symbol ON tickers(symbol)
	`).Error; err != nil {
		logger.Error().Err(err).Msg("Failed to add symbol index to tickers table")
		return err
	}

	logger.Info().Msg("Ticker symbols index added successfully")
	return nil
}

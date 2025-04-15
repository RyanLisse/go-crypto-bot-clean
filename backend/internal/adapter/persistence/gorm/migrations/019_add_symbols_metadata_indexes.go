package migrations

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// AddSymbolsMetadataIndexes adds indexes for symbol metadata
func AddSymbolsMetadataIndexes(db *gorm.DB) error {
	logger := log.With().Str("migration", "add_symbols_metadata_indexes").Logger()
	logger.Info().Msg("Running migration: Add symbols metadata indexes")

	// Add indexes to the symbols table
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_symbols_base_asset ON symbols(base_asset);
		CREATE INDEX IF NOT EXISTS idx_symbols_status ON symbols(status);
		CREATE INDEX IF NOT EXISTS idx_symbols_listing_status ON symbols(listing_status);
	`).Error; err != nil {
		logger.Error().Err(err).Msg("Failed to add metadata indexes to symbols table")
		return err
	}

	logger.Info().Msg("Symbols metadata indexes added successfully")
	return nil
}

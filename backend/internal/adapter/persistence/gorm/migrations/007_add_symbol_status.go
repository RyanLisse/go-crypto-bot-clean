package migrations

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// AddSymbolStatus adds a status column to the symbols table
func AddSymbolStatus(db *gorm.DB) error {
	logger := log.With().Str("migration", "add_symbol_status").Logger()
	logger.Info().Msg("Running migration: Add symbol status")

	// Check if the column already exists
	var columnExists bool
	db.Raw(`
		SELECT COUNT(*) > 0 
		FROM pragma_table_info('symbols') 
		WHERE name = 'listing_status'
	`).Scan(&columnExists)

	if !columnExists {
		// Add the listing_status column
		if err := db.Exec(`
			ALTER TABLE symbols 
			ADD COLUMN listing_status VARCHAR(20) DEFAULT 'LISTED'
		`).Error; err != nil {
			logger.Error().Err(err).Msg("Failed to add listing_status column to symbols table")
			return err
		}
	}

	logger.Info().Msg("Symbol status added successfully")
	return nil
}

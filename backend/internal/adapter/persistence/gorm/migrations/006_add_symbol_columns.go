package migrations

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// AddSymbolColumns adds additional columns to the symbols table
func AddSymbolColumns(db *gorm.DB) error {
	logger := log.With().Str("migration", "add_symbol_columns").Logger()
	logger.Info().Msg("Running migration: Add symbol columns")

	// Use GORM's AutoMigrate to add new columns in a DB-agnostic way
	if err := db.AutoMigrate(&entity.MexcSymbolEntity{}); err != nil {
		logger.Error().Err(err).Msg("Failed to add columns to symbols table via AutoMigrate")
		return err
	}

	logger.Info().Msg("Symbol columns added successfully via AutoMigrate")
	return nil
}

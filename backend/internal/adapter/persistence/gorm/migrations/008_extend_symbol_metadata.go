package migrations

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// ExtendSymbolMetadata adds additional metadata columns to the symbols table
func ExtendSymbolMetadata(db *gorm.DB) error {
	logger := log.With().Str("migration", "extend_symbol_metadata").Logger()
	logger.Info().Msg("Running migration: Extend symbol metadata")

	// Use GORM's AutoMigrate to add new columns in a DB-agnostic way
	if err := db.AutoMigrate(&entity.MexcSymbolEntity{}); err != nil {
		logger.Error().Err(err).Msg("Failed to add metadata columns to symbols table via AutoMigrate")
		return err
	}

	logger.Info().Msg("Symbol metadata extended successfully via AutoMigrate")
	return nil
}

package migrations

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func CreatePositionsTable(db *gorm.DB) error {
	logger := log.With().Str("migration", "create_positions_table").Logger()
	logger.Info().Msg("Running migration: Create positions table")
	if err := db.AutoMigrate(&entity.PositionEntity{}); err != nil {
		logger.Error().Err(err).Msg("Failed to create positions table")
		return err
	}
	logger.Info().Msg("Positions table created successfully")
	return nil
}

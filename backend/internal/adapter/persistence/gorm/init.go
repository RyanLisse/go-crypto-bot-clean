package gorm

import (
	"github.com/rs/zerolog"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
)

// NewDB is a helper to create a DB connection or log fatal on error.
import "gorm.io/gorm"

func NewDB(cfg *config.Config, logger *zerolog.Logger) *gorm.DB {
	db, err := NewDBConnection(cfg, *logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	return db
}

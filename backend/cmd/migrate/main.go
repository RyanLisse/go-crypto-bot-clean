package main

import (
	"fmt"
	"os"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/infrastructure/database"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Configure logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	logger := log.With().Str("component", "migration").Logger()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Connect to database
	db, err := database.Connect(cfg, &logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Run migrations
	logger.Info().Msg("Starting database migrations")
	if err := database.RunMigrations(db, &logger); err != nil {
		logger.Fatal().Err(err).Msg("Failed to run migrations")
	}

	fmt.Println("Migrations completed successfully")
}

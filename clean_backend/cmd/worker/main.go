package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/di"
)

func main() {
	// Setup logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	logger := log.With().Str("component", "worker").Logger()
	logger.Info().Msg("Starting worker...")

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		logger.Warn().Err(err).Msg("Error loading .env file")
	}

	// Initialize dependency container
	_, err := di.NewContainer()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize dependency container")
	}

	// TODO: Initialize and start worker tasks
	logger.Info().Msg("Worker initialized, but no tasks are configured yet")

	// Graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	logger.Info().Msg("Shutting down worker...")
	// Create a context with timeout for graceful shutdown
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// TODO: Implement worker shutdown logic
	logger.Info().Msg("Worker shutdown complete")
}

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
	logger := log.With().Str("component", "server").Logger()
	logger.Info().Msg("Starting server...")

	// Kill any process using port 8080 before starting
	if err := killProcessOnPort(8080, logger); err != nil {
		logger.Warn().Err(err).Int("port", 8080).Msg("Failed to kill process on port")
	}

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		logger.Warn().Err(err).Msg("Error loading .env file")
	}

	// Initialize dependency container
	container, err := di.NewContainer()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize dependency container")
	}

	// Get server port from config
	port := container.Config.Server.Port

	// Start server in a goroutine
	go func() {
		logger.Info().Int("port", port).Msg("HTTP server starting")
		if err := container.Server.Start(port); err != nil {
			logger.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	logger.Info().Msg("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := container.Server.Stop(ctx); err != nil {
		logger.Error().Err(err).Msg("Server shutdown error")
	}

	logger.Info().Msg("Server shutdown complete")
}

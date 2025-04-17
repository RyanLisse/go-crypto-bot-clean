package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/server"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/rs/zerolog"
)

func main() {
	// Set up logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger.Info().Msg("Starting error handling example application")

	// Create a simple config
	cfg := &config.Config{
		Server: struct {
			Port               int           `mapstructure:"port"`
			Host               string        `mapstructure:"host"`
			ReadTimeout        time.Duration `mapstructure:"read_timeout"`
			WriteTimeout       time.Duration `mapstructure:"write_timeout"`
			IdleTimeout        time.Duration `mapstructure:"idle_timeout"`
			FrontendURL        string        `mapstructure:"frontend_url"`
			CORSAllowedOrigins []string      `mapstructure:"cors_allowed_origins"`
		}{
			Port: 8085,
		},
	}

	// Create and configure the example server
	srv := server.NewExampleServer(cfg, &logger)
	if err := srv.SetupRoutes(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to set up routes")
	}

	// Start the server in a goroutine
	go func() {
		logger.Info().Int("port", cfg.Server.Port).Msg("Server starting")
		if err := srv.Start(cfg.Server.Port); err != nil {
			logger.Error().Err(err).Msg("Server failed")
		}
	}()

	// Print available endpoints
	logger.Info().Msg("Available error example endpoints:")
	logger.Info().Msg("GET  /health                  - Health check")
	logger.Info().Msg("GET  /errors/not-found        - 404 Not Found example")
	logger.Info().Msg("GET  /errors/unauthorized     - 401 Unauthorized example")
	logger.Info().Msg("GET  /errors/forbidden        - 403 Forbidden example")
	logger.Info().Msg("GET  /errors/internal         - 500 Internal Server Error example")
	logger.Info().Msg("GET  /errors/validation/single - Single field validation error")
	logger.Info().Msg("GET  /errors/validation/multiple - Multiple field validation errors")
	logger.Info().Msg("POST /errors/validation       - Input validation (requires JSON body)")
	logger.Info().Msg("GET  /errors/wrapped          - Wrapped error chain example")
	logger.Info().Msg("GET  /errors/external-api     - External API error example")
	logger.Info().Msg("GET  /errors/panic            - Panic recovery example")

	// Example curl commands
	logger.Info().Msg("\nExample curl commands:")
	logger.Info().Msgf("curl http://localhost:%d/health", cfg.Server.Port)
	logger.Info().Msgf("curl http://localhost:%d/errors/not-found", cfg.Server.Port)
	logger.Info().Msgf("curl http://localhost:%d/errors/validation/multiple", cfg.Server.Port)
	logger.Info().Msgf("curl -X POST -H \"Content-Type: application/json\" -d '{\"email\":\"\",\"username\":\"u\",\"age\":16}' http://localhost:%d/errors/validation", cfg.Server.Port)
	logger.Info().Msgf("curl http://localhost:%d/errors/panic", cfg.Server.Port)

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will wait until the timeout deadline
	if err := srv.Stop(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("Server exited properly")
}

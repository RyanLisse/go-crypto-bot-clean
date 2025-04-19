package logger

import (
	"os"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/rs/zerolog"
)

// TODO: Setup zerolog or chosen logger.
func NewLogger(cfg config.LogConfig) *zerolog.Logger {
	logLevel, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	// Use ConsoleWriter for development for pretty printing
	// logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	return &logger
}

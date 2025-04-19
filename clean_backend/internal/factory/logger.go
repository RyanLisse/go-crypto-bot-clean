package factory

import (
	"os"
	"strings"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/rs/zerolog"
)

// NewLogger creates a new zerolog logger with the specified configuration
func NewLogger(cfg config.LogConfig) *zerolog.Logger {
	// Set global log level
	level := zerolog.InfoLevel
	if cfg.Level != "" {
		switch strings.ToLower(cfg.Level) {
		case "debug":
			level = zerolog.DebugLevel
		case "info":
			level = zerolog.InfoLevel
		case "warn":
			level = zerolog.WarnLevel
		case "error":
			level = zerolog.ErrorLevel
		case "fatal":
			level = zerolog.FatalLevel
		case "panic":
			level = zerolog.PanicLevel
		case "trace":
			level = zerolog.TraceLevel
		}
	}
	zerolog.SetGlobalLevel(level)

	// Configure output format
	var output zerolog.ConsoleWriter
	if cfg.Pretty {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	} else {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			NoColor:    true,
			TimeFormat: time.RFC3339,
		}
	}

	// Create logger
	logger := zerolog.New(output).
		With().
		Timestamp().
		Str("service", "crypto-bot").
		Logger()

	return &logger
}

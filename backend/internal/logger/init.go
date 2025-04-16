package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Configure global zerolog settings
func init() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

// NewLogger returns a zerolog.Logger instance with timestamp, caller info, and appropriate output format.
func NewLogger() *zerolog.Logger {
	return NewLoggerWithLevel("info")
}

// NewLoggerWithLevel creates a new zerolog Logger with the specified log level
func NewLoggerWithLevel(level string) *zerolog.Logger {
	// Set the logging level
	setLogLevel(level)

	// Create console writer
	var output io.Writer = os.Stdout
	// If development environment, use pretty console output
	if os.Getenv("ENV") == "development" {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			NoColor:    false,
		}
	}

	// Create the logger
	logger := zerolog.New(output).With().Timestamp().Caller().Logger()
	return &logger
}

// setLogLevel sets the global log level from a string
func setLogLevel(level string) {
	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

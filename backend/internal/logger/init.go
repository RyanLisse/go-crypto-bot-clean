package logger

import (
	"os"
	"github.com/rs/zerolog"
)

// NewLogger returns a zerolog.Logger instance with timestamp.
func NewLogger() *zerolog.Logger {
	l := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &l
}

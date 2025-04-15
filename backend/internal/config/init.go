package config

import (
	"github.com/rs/zerolog"
)

// LoadConfig is a helper to load configuration and log fatal on error.
func LoadConfig(logger *zerolog.Logger) *Config {
	cfg, err := Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load configuration")
	}
	return cfg
}

package factory

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/pkg/platform/mexc"
	"github.com/rs/zerolog"
)

// NewMEXCClient creates a new MEXC client
func NewMEXCClient(cfg *config.Config, logger *zerolog.Logger) port.MEXCClient {
	return mexc.NewClient(cfg.MEXC.APIKey, cfg.MEXC.APISecret, logger)
}

package factory

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/pkg/platform/mexc"
	"github.com/rs/zerolog"
)

// MEXCFactory creates MEXC-related components
type MEXCFactory struct {
	cfg    *config.Config
	logger *zerolog.Logger
}

// NewMEXCFactory creates a new MEXCFactory
func NewMEXCFactory(cfg *config.Config, logger *zerolog.Logger) *MEXCFactory {
	return &MEXCFactory{
		cfg:    cfg,
		logger: logger,
	}
}

// CreateMEXCClient creates a MEXC client
func (f *MEXCFactory) CreateMEXCClient() port.MEXCClient {
	apiKey := f.cfg.MEXC.APIKey
	apiSecret := f.cfg.MEXC.APISecret

	f.logger.Debug().Msg("Creating MEXC client")
	return mexc.NewClient(apiKey, apiSecret, f.logger)
}

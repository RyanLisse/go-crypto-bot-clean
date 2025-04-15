package factory

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/cache/standard"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
)

// WalletCacheFactory creates wallet cache instances
// Similar to CacheFactory but for wallet data

type WalletCacheFactory struct {
	logger *zerolog.Logger
}

func NewWalletCacheFactory(logger *zerolog.Logger) *WalletCacheFactory {
	return &WalletCacheFactory{logger: logger}
}

func (f *WalletCacheFactory) CreateWalletCache() port.WalletCache {
	defaultTTL := 5 * time.Minute
	cleanupInterval := 10 * time.Minute
	f.logger.Info().Dur("defaultTTL", defaultTTL).Dur("cleanupInterval", cleanupInterval).Msg("Creating wallet cache")
	return standard.NewWalletCache(defaultTTL, cleanupInterval)
}

package di

import (
	"github.com/rs/zerolog"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/util/crypto"
)

// provideKeyManager creates and returns a new KeyManager
func provideKeyManager(cfg *config.Config, logger *zerolog.Logger) (*crypto.KeyManager, error) {
	return crypto.NewKeyManager(cfg.Encryption.CurrentKeyID, cfg.Encryption.Keys, logger)
}

// provideEncryptionService creates and returns a new EncryptionService
func provideEncryptionService(keyManager *crypto.KeyManager, logger *zerolog.Logger) port.EncryptionService {
	return crypto.NewEnhancedEncryptionService(keyManager, logger)
}

// provideEnhancedEncryptionService creates and returns a new EnhancedEncryptionService
func provideEnhancedEncryptionService(keyManager *crypto.KeyManager, logger *zerolog.Logger) port.EnhancedEncryptionService {
	return crypto.NewEnhancedEncryptionService(keyManager, logger)
}

package factory

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/wallet"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/util/crypto"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// APICredentialManagerFactory creates API credential manager components
type APICredentialManagerFactory struct {
	cfg    *config.Config
	logger *zerolog.Logger
	db     *gorm.DB
}

// NewAPICredentialManagerFactory creates a new APICredentialManagerFactory
func NewAPICredentialManagerFactory(cfg *config.Config, logger *zerolog.Logger, db *gorm.DB) *APICredentialManagerFactory {
	return &APICredentialManagerFactory{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}
}

// CreateAPICredentialManagerService creates an API credential manager service
func (f *APICredentialManagerFactory) CreateAPICredentialManagerService(
	repo port.APICredentialRepository,
	encryptionSvc crypto.EncryptionService,
	providerRegistry *wallet.ProviderRegistry,
) usecase.APICredentialManagerService {
	return usecase.NewAPICredentialManagerService(
		repo,
		encryptionSvc,
		providerRegistry,
		f.logger,
	)
}

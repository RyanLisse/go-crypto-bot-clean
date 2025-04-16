package factory

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/delivery/http/handler"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/repo"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/wallet"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/util/crypto"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// APICredentialFactory creates API credential-related components
type APICredentialFactory struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewAPICredentialFactory creates a new APICredentialFactory
func NewAPICredentialFactory(db *gorm.DB, logger *zerolog.Logger) *APICredentialFactory {
	return &APICredentialFactory{
		db:     db,
		logger: logger,
	}
}

// CreateAPICredentialRepository creates a new API credential repository
func (f *APICredentialFactory) CreateAPICredentialRepository() *repo.APICredentialRepository {
	// Create encryption service
	encryptionService, err := crypto.NewAESEncryptionService()
	if err != nil {
		f.logger.Error().Err(err).Msg("Failed to create AESEncryptionService")
		return nil
	}

	// Create repository
	return repo.NewAPICredentialRepository(f.db, encryptionService, f.logger)
}

// CreateAPICredentialHandler creates a new API credential handler
func (f *APICredentialFactory) CreateAPICredentialHandler() *handler.APICredentialHandler {
	// Create repository
	repository := f.CreateAPICredentialRepository()

	// Create use case
	useCase := usecase.NewAPICredentialUseCase(repository, f.logger)

	// Create handler
	return handler.NewAPICredentialHandler(useCase, f.logger)
}

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

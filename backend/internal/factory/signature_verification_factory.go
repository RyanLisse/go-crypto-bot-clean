package factory

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/delivery/http/handler"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/wallet"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// SignatureVerificationFactory creates signature verification components
type SignatureVerificationFactory struct {
	cfg    *config.Config
	logger *zerolog.Logger
	db     *gorm.DB
}

// NewSignatureVerificationFactory creates a new SignatureVerificationFactory
func NewSignatureVerificationFactory(cfg *config.Config, logger *zerolog.Logger, db *gorm.DB) *SignatureVerificationFactory {
	return &SignatureVerificationFactory{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}
}

// CreateSignatureVerificationService creates a signature verification service
func (f *SignatureVerificationFactory) CreateSignatureVerificationService(
	providerRegistry *wallet.ProviderRegistry,
	walletRepo port.WalletRepository,
) usecase.SignatureVerificationService {
	return usecase.NewSignatureVerificationService(
		providerRegistry,
		walletRepo,
		f.logger,
	)
}

// CreateSignatureVerificationHandler creates a signature verification handler
func (f *SignatureVerificationFactory) CreateSignatureVerificationHandler(
	verificationService usecase.SignatureVerificationService,
) *handler.SignatureVerificationHandler {
	return handler.NewSignatureVerificationHandler(
		verificationService,
		f.logger,
	)
}

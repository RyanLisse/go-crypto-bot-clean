package factory

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/wallet"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// WalletDataSyncFactory creates wallet data synchronization components
type WalletDataSyncFactory struct {
	cfg    *config.Config
	logger *zerolog.Logger
	db     *gorm.DB
}

// NewWalletDataSyncFactory creates a new WalletDataSyncFactory
func NewWalletDataSyncFactory(cfg *config.Config, logger *zerolog.Logger, db *gorm.DB) *WalletDataSyncFactory {
	return &WalletDataSyncFactory{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}
}

// CreateWalletDataSyncService creates a wallet data synchronization service
func (f *WalletDataSyncFactory) CreateWalletDataSyncService(
	walletRepo port.WalletRepository,
	apiCredentialManager usecase.APICredentialManagerService,
	providerRegistry *wallet.ProviderRegistry,
) usecase.WalletDataSyncService {
	return usecase.NewWalletDataSyncService(
		walletRepo,
		apiCredentialManager,
		providerRegistry,
		f.logger,
	)
}

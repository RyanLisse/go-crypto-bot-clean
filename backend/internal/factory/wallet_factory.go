package factory

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/delivery/http/handler"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/repo"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// WalletFactory creates wallet-related components
type WalletFactory struct {
	cfg    *config.Config
	logger *zerolog.Logger
	db     *gorm.DB
}

// NewWalletFactory creates a new WalletFactory
func NewWalletFactory(cfg *config.Config, logger *zerolog.Logger, db *gorm.DB) *WalletFactory {
	return &WalletFactory{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}
}

// CreateWalletRepository creates a wallet repository
func (f *WalletFactory) CreateWalletRepository() port.WalletRepository {
	return repo.NewConsolidatedWalletRepository(f.db, f.logger)
}

// CreateWalletService creates a wallet service
func (f *WalletFactory) CreateWalletService(mexcClient port.MEXCClient) usecase.WalletService {
	walletRepo := f.CreateWalletRepository()
	return usecase.NewWalletService(walletRepo, mexcClient, f.logger)
}

// CreateWalletHandler creates a wallet handler
func (f *WalletFactory) CreateWalletHandler(walletService usecase.WalletService) *handler.WalletHandler {
	return handler.NewWalletHandler(walletService, f.logger)
}

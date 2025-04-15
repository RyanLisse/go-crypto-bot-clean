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

// Web3WalletFactory creates Web3 wallet components
type Web3WalletFactory struct {
	cfg    *config.Config
	logger *zerolog.Logger
	db     *gorm.DB
}

// NewWeb3WalletFactory creates a new Web3WalletFactory
func NewWeb3WalletFactory(cfg *config.Config, logger *zerolog.Logger, db *gorm.DB) *Web3WalletFactory {
	return &Web3WalletFactory{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}
}

// CreateWeb3WalletService creates a Web3 wallet service
func (f *Web3WalletFactory) CreateWeb3WalletService(
	walletRepo port.WalletRepository,
	providerRegistry *wallet.ProviderRegistry,
) usecase.Web3WalletService {
	return usecase.NewWeb3WalletService(
		walletRepo,
		providerRegistry,
		f.logger,
	)
}

// CreateWeb3WalletHandler creates a Web3 wallet handler
func (f *Web3WalletFactory) CreateWeb3WalletHandler(
	web3WalletService usecase.Web3WalletService,
) *handler.Web3WalletHandler {
	return handler.NewWeb3WalletHandler(
		web3WalletService,
		f.logger,
	)
}

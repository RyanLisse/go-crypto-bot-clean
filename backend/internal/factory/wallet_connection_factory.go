package factory

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/wallet"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// WalletConnectionFactory creates wallet connection components
type WalletConnectionFactory struct {
	cfg    *config.Config
	logger *zerolog.Logger
	db     *gorm.DB
}

// NewWalletConnectionFactory creates a new WalletConnectionFactory
func NewWalletConnectionFactory(cfg *config.Config, logger *zerolog.Logger, db *gorm.DB) *WalletConnectionFactory {
	return &WalletConnectionFactory{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}
}

// CreateProviderRegistry creates a wallet provider registry
func (f *WalletConnectionFactory) CreateProviderRegistry(mexcClient port.MEXCClient) *wallet.ProviderRegistry {
	registry := wallet.NewProviderRegistry()

	// Register MEXC provider
	mexcProvider := wallet.NewMEXCProvider(mexcClient, f.logger)
	registry.RegisterProvider(mexcProvider)

	// Register Ethereum provider
	ethereumProvider := wallet.NewEthereumProvider(
		1, // Ethereum Mainnet
		"Ethereum",
		"https://mainnet.infura.io/v3/" + f.cfg.InfuraAPIKey,
		"https://etherscan.io",
		f.logger,
	)
	registry.RegisterProvider(ethereumProvider)

	return registry
}

// CreateWalletConnectionService creates a wallet connection service
func (f *WalletConnectionFactory) CreateWalletConnectionService(
	providerRegistry *wallet.ProviderRegistry,
	walletRepo port.WalletRepository,
) usecase.WalletConnectionService {
	return usecase.NewWalletConnectionService(
		providerRegistry,
		walletRepo,
		f.logger,
	)
}

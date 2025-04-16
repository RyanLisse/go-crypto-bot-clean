package port

import "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"

// ProviderRegistry defines the interface for wallet provider registry
type ProviderRegistry interface {
	// RegisterProvider registers a wallet provider
	RegisterProvider(provider WalletProvider)

	// GetProvider gets a wallet provider by name
	GetProvider(name string) (WalletProvider, error)

	// GetExchangeProvider gets an exchange wallet provider by name
	GetExchangeProvider(name string) (ExchangeWalletProvider, error)

	// GetWeb3Provider gets a Web3 wallet provider by name
	GetWeb3Provider(name string) (Web3WalletProvider, error)

	// GetProviderByType gets wallet providers by type
	GetProviderByType(typ model.WalletType) ([]WalletProvider, error)

	// GetAllProviders gets all wallet providers
	GetAllProviders() []WalletProvider

	// GetAllExchangeProviders gets all exchange wallet providers
	GetAllExchangeProviders() []ExchangeWalletProvider

	// GetAllWeb3Providers gets all Web3 wallet providers
	GetAllWeb3Providers() []Web3WalletProvider
}

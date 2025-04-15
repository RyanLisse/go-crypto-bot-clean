package wallet

import (
	"errors"
	"sync"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
)

// ProviderRegistry manages wallet providers
type ProviderRegistry struct {
	providers     map[string]port.WalletProvider
	exchangeProviders map[string]port.ExchangeWalletProvider
	web3Providers     map[string]port.Web3WalletProvider
	mu           sync.RWMutex
}

// NewProviderRegistry creates a new wallet provider registry
func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		providers:     make(map[string]port.WalletProvider),
		exchangeProviders: make(map[string]port.ExchangeWalletProvider),
		web3Providers:     make(map[string]port.Web3WalletProvider),
	}
}

// RegisterProvider registers a wallet provider
func (r *ProviderRegistry) RegisterProvider(provider port.WalletProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.providers[provider.GetName()] = provider

	// Register in type-specific maps
	if exchangeProvider, ok := provider.(port.ExchangeWalletProvider); ok {
		// Ensure MEXCProvider implements port.ExchangeWalletProvider
		var _ port.ExchangeWalletProvider = (*MEXCProvider)(nil)
		// Ensure CoinbaseProvider implements port.ExchangeWalletProvider
		var _ port.ExchangeWalletProvider = (*CoinbaseProvider)(nil)
		r.exchangeProviders[provider.GetName()] = exchangeProvider
	}

	if web3Provider, ok := provider.(port.Web3WalletProvider); ok {
		r.web3Providers[provider.GetName()] = web3Provider
	}
}

// GetProvider gets a wallet provider by name
func (r *ProviderRegistry) GetProvider(name string) (port.WalletProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, ok := r.providers[name]
	if !ok {
		return nil, errors.New("provider not found")
	}

	return provider, nil
}

// GetExchangeProvider gets an exchange wallet provider by name
func (r *ProviderRegistry) GetExchangeProvider(name string) (port.ExchangeWalletProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, ok := r.exchangeProviders[name]
	if !ok {
		return nil, errors.New("exchange provider not found")
	}

	return provider, nil
}

// GetWeb3Provider gets a Web3 wallet provider by name
func (r *ProviderRegistry) GetWeb3Provider(name string) (port.Web3WalletProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, ok := r.web3Providers[name]
	if !ok {
		return nil, errors.New("web3 provider not found")
	}

	return provider, nil
}

// GetProviderByType gets a wallet provider by type
func (r *ProviderRegistry) GetProviderByType(typ model.WalletType) ([]port.WalletProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var providers []port.WalletProvider
	for _, provider := range r.providers {
		if provider.GetType() == typ {
			providers = append(providers, provider)
		}
	}

	if len(providers) == 0 {
		return nil, errors.New("no providers found for type")
	}

	return providers, nil
}

// GetAllProviders gets all wallet providers
func (r *ProviderRegistry) GetAllProviders() []port.WalletProvider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var providers []port.WalletProvider
	for _, provider := range r.providers {
		providers = append(providers, provider)
	}

	return providers
}

// GetAllExchangeProviders gets all exchange wallet providers
func (r *ProviderRegistry) GetAllExchangeProviders() []port.ExchangeWalletProvider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var providers []port.ExchangeWalletProvider
	for _, provider := range r.exchangeProviders {
		providers = append(providers, provider)
	}

	return providers
}

// GetAllWeb3Providers gets all Web3 wallet providers
func (r *ProviderRegistry) GetAllWeb3Providers() []port.Web3WalletProvider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var providers []port.Web3WalletProvider
	for _, provider := range r.web3Providers {
		providers = append(providers, provider)
	}

	return providers
}

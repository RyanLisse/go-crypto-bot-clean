package usecase

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/wallet"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// AddressValidatorService defines the interface for wallet address validation
type AddressValidatorService interface {
	// ValidateAddress validates a wallet address for a specific network
	ValidateAddress(ctx context.Context, network, address string) (bool, error)

	// GetAddressInfo returns information about a wallet address
	GetAddressInfo(ctx context.Context, network, address string) (*AddressInfo, error)

	// GetSupportedNetworks returns a list of supported networks
	GetSupportedNetworks(ctx context.Context) ([]string, error)
}

// AddressInfo contains information about a wallet address
type AddressInfo struct {
	Network     string `json:"network"`      // Network name (e.g., "Ethereum", "Bitcoin")
	Address     string `json:"address"`      // The wallet address
	IsValid     bool   `json:"is_valid"`     // Whether the address is valid
	AddressType string `json:"address_type"` // Type of address (e.g., "EOA", "Contract", "P2PKH")
	ChainID     int64  `json:"chain_id"`     // Chain ID for the network
	Explorer    string `json:"explorer"`     // Block explorer URL
}

// addressValidatorService implements the AddressValidatorService interface
type addressValidatorService struct {
	providerRegistry *wallet.ProviderRegistry
	logger           *zerolog.Logger
}

// NewAddressValidatorService creates a new AddressValidatorService
func NewAddressValidatorService(
	providerRegistry *wallet.ProviderRegistry,
	logger *zerolog.Logger,
) AddressValidatorService {
	return &addressValidatorService{
		providerRegistry: providerRegistry,
		logger:           logger,
	}
}

// ValidateAddress validates a wallet address for a specific network
func (s *addressValidatorService) ValidateAddress(ctx context.Context, network, address string) (bool, error) {
	// Trim whitespace
	address = strings.TrimSpace(address)
	if address == "" {
		return false, errors.New("address cannot be empty")
	}

	// Get the provider for the network
	provider, err := s.providerRegistry.GetProvider(network)
	if err != nil {
		s.logger.Error().Err(err).Str("network", network).Msg("Failed to get provider")
		return false, fmt.Errorf("unsupported network: %s", network)
	}

	// Validate the address using the provider
	return provider.IsValidAddress(ctx, address)
}

// GetAddressInfo returns information about a wallet address
func (s *addressValidatorService) GetAddressInfo(ctx context.Context, network, address string) (*AddressInfo, error) {
	// Trim whitespace
	address = strings.TrimSpace(address)
	if address == "" {
		return nil, errors.New("address cannot be empty")
	}

	// Get the provider for the network
	provider, err := s.providerRegistry.GetProvider(network)
	if err != nil {
		s.logger.Error().Err(err).Str("network", network).Msg("Failed to get provider")
		return nil, fmt.Errorf("unsupported network: %s", network)
	}

	// Validate the address using the provider
	isValid, err := provider.IsValidAddress(ctx, address)
	if err != nil {
		s.logger.Error().Err(err).Str("network", network).Str("address", address).Msg("Failed to validate address")
		return nil, err
	}

	// Create address info
	info := &AddressInfo{
		Network: network,
		Address: address,
		IsValid: isValid,
	}

	// Add additional information based on the network
	switch network {
	case "Ethereum":
		// Check if the address is a contract or EOA
		// This is a simplified check - in a real implementation, we would use the provider to check
		if isValid {
			// For now, assume it's an EOA (Externally Owned Account)
			info.AddressType = "EOA"
		}

		// Get chain ID and explorer URL from the provider
		if web3Provider, ok := provider.(port.Web3WalletProvider); ok {
			info.ChainID = web3Provider.GetChainID()
			info.Explorer = fmt.Sprintf("https://etherscan.io/address/%s", address)
		} else {
			// For testing purposes
			info.Explorer = fmt.Sprintf("https://etherscan.io/address/%s", address)
		}
	case "Bitcoin":
		// Determine Bitcoin address type
		if isValid {
			info.AddressType = determineBitcoinAddressType(address)
		}
		info.Explorer = fmt.Sprintf("https://www.blockchain.com/explorer/addresses/btc/%s", address)
	default:
		// For other networks, just set a generic address type
		if isValid {
			info.AddressType = "Unknown"
		}
	}

	return info, nil
}

// GetSupportedNetworks returns a list of supported networks
func (s *addressValidatorService) GetSupportedNetworks(ctx context.Context) ([]string, error) {
	// Get all providers
	providers := s.providerRegistry.GetAllProviders()

	// Extract network names
	networks := make([]string, 0, len(providers))
	for _, provider := range providers {
		networks = append(networks, provider.GetName())
	}

	return networks, nil
}

// determineBitcoinAddressType determines the type of Bitcoin address
func determineBitcoinAddressType(address string) string {
	// P2PKH addresses start with 1
	if regexp.MustCompile(`^1[a-km-zA-HJ-NP-Z1-9]{25,34}$`).MatchString(address) {
		return "P2PKH"
	}

	// P2SH addresses start with 3
	if regexp.MustCompile(`^3[a-km-zA-HJ-NP-Z1-9]{25,34}$`).MatchString(address) {
		return "P2SH"
	}

	// Bech32 addresses start with bc1
	if regexp.MustCompile(`^bc1[a-zA-HJ-NP-Z0-9]{25,89}$`).MatchString(address) {
		return "Bech32"
	}

	return "Unknown"
}

package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/wallet"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// Web3WalletService defines the interface for Web3 wallet operations
type Web3WalletService interface {
	// ConnectWallet connects a Web3 wallet
	ConnectWallet(ctx context.Context, userID, network, address string) (*model.Wallet, error)

	// DisconnectWallet disconnects a Web3 wallet
	DisconnectWallet(ctx context.Context, walletID string) error

	// GetWalletBalance gets the balance of a Web3 wallet
	GetWalletBalance(ctx context.Context, walletID string) (*model.Wallet, error)

	// GetWalletByAddress gets a wallet by its address
	GetWalletByAddress(ctx context.Context, network, address string) (*model.Wallet, error)

	// IsValidAddress checks if an address is valid for the given network
	IsValidAddress(ctx context.Context, network, address string) (bool, error)

	// GetSupportedNetworks gets the list of supported networks
	GetSupportedNetworks(ctx context.Context) ([]string, error)
}

// web3WalletService implements the Web3WalletService interface
type web3WalletService struct {
	walletRepo       port.WalletRepository
	providerRegistry *wallet.ProviderRegistry
	logger           *zerolog.Logger
}

// NewWeb3WalletService creates a new Web3WalletService
func NewWeb3WalletService(
	walletRepo port.WalletRepository,
	providerRegistry *wallet.ProviderRegistry,
	logger *zerolog.Logger,
) Web3WalletService {
	return &web3WalletService{
		walletRepo:       walletRepo,
		providerRegistry: providerRegistry,
		logger:           logger,
	}
}

// ConnectWallet connects a Web3 wallet
func (s *web3WalletService) ConnectWallet(ctx context.Context, userID, network, address string) (*model.Wallet, error) {
	// Get the provider for the network
	provider, err := s.providerRegistry.GetWeb3Provider(network)
	if err != nil {
		s.logger.Error().Err(err).Str("network", network).Msg("Failed to get Web3 provider")
		return nil, fmt.Errorf("unsupported network: %s", network)
	}

	// Check if the address is valid
	valid, err := provider.IsValidAddress(ctx, address)
	if err != nil {
		s.logger.Error().Err(err).Str("address", address).Msg("Failed to validate address")
		return nil, err
	}
	if !valid {
		return nil, errors.New("invalid address")
	}

	// Check if the wallet already exists
	existingWallet, err := s.GetWalletByAddress(ctx, network, address)
	if err == nil && existingWallet != nil {
		// If the wallet exists but belongs to a different user, return an error
		if existingWallet.UserID != userID {
			return nil, errors.New("wallet already connected to another user")
		}
		// If the wallet exists and belongs to the same user, return it
		return existingWallet, nil
	}

	// Connect the wallet
	params := map[string]interface{}{
		"user_id": userID,
		"address": address,
	}
	wallet, err := provider.Connect(ctx, params)
	if err != nil {
		s.logger.Error().Err(err).Str("network", network).Str("address", address).Msg("Failed to connect wallet")
		return nil, err
	}

	// Save the wallet
	err = s.walletRepo.Save(ctx, wallet)
	if err != nil {
		s.logger.Error().Err(err).Str("id", wallet.ID).Msg("Failed to save wallet")
		return nil, err
	}

	return wallet, nil
}

// DisconnectWallet disconnects a Web3 wallet
func (s *web3WalletService) DisconnectWallet(ctx context.Context, walletID string) error {
	// Get the wallet
	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get wallet")
		return err
	}
	if wallet == nil {
		return errors.New("wallet not found")
	}

	// Check if the wallet is a Web3 wallet
	if wallet.Type != model.WalletTypeWeb3 {
		return errors.New("not a Web3 wallet")
	}

	// Get the provider for the network
	provider, err := s.providerRegistry.GetWeb3Provider(wallet.Metadata.Network)
	if err != nil {
		s.logger.Error().Err(err).Str("network", wallet.Metadata.Network).Msg("Failed to get Web3 provider")
		return err
	}

	// Disconnect the wallet
	err = provider.Disconnect(ctx, walletID)
	if err != nil {
		s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to disconnect wallet")
		return err
	}

	// Delete the wallet
	err = s.walletRepo.DeleteWallet(ctx, walletID)
	if err != nil {
		s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to delete wallet")
		return err
	}

	return nil
}

// GetWalletBalance gets the balance of a Web3 wallet
func (s *web3WalletService) GetWalletBalance(ctx context.Context, walletID string) (*model.Wallet, error) {
	// Get the wallet
	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get wallet")
		return nil, err
	}
	if wallet == nil {
		return nil, errors.New("wallet not found")
	}

	// Check if the wallet is a Web3 wallet
	if wallet.Type != model.WalletTypeWeb3 {
		return nil, errors.New("not a Web3 wallet")
	}

	// Get the provider for the network
	provider, err := s.providerRegistry.GetWeb3Provider(wallet.Metadata.Network)
	if err != nil {
		s.logger.Error().Err(err).Str("network", wallet.Metadata.Network).Msg("Failed to get Web3 provider")
		return nil, err
	}

	// Get the balance
	updatedWallet, err := provider.GetBalance(ctx, wallet)
	if err != nil {
		s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get wallet balance")
		return nil, err
	}

	// Update the wallet
	updatedWallet.LastUpdated = time.Now()
	updatedWallet.LastSyncAt = time.Now()
	err = s.walletRepo.Save(ctx, updatedWallet)
	if err != nil {
		s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to save wallet")
		return nil, err
	}

	return updatedWallet, nil
}

// GetWalletByAddress gets a wallet by its address
func (s *web3WalletService) GetWalletByAddress(ctx context.Context, network, address string) (*model.Wallet, error) {
	// Get all wallets for the user
	wallets, err := s.walletRepo.GetWalletsByUserID(ctx, "")
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get wallets")
		return nil, err
	}

	// Find the wallet with the matching address and network
	for _, wallet := range wallets {
		if wallet.Type == model.WalletTypeWeb3 &&
			wallet.Metadata != nil &&
			wallet.Metadata.Network == network &&
			wallet.Metadata.Address == address {
			return wallet, nil
		}
	}

	return nil, errors.New("wallet not found")
}

// IsValidAddress checks if an address is valid for the given network
func (s *web3WalletService) IsValidAddress(ctx context.Context, network, address string) (bool, error) {
	// Get the provider for the network
	provider, err := s.providerRegistry.GetWeb3Provider(network)
	if err != nil {
		s.logger.Error().Err(err).Str("network", network).Msg("Failed to get Web3 provider")
		return false, fmt.Errorf("unsupported network: %s", network)
	}

	// Check if the address is valid
	return provider.IsValidAddress(ctx, address)
}

// GetSupportedNetworks gets the list of supported networks
func (s *web3WalletService) GetSupportedNetworks(ctx context.Context) ([]string, error) {
	// Get all Web3 providers
	providers := s.providerRegistry.GetAllWeb3Providers()
	
	// Extract network names
	networks := make([]string, 0, len(providers))
	for _, provider := range providers {
		networks = append(networks, provider.GetName())
	}

	return networks, nil
}

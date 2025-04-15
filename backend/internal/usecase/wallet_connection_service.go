package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/wallet"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// WalletConnectionService defines the interface for wallet connection operations
type WalletConnectionService interface {
	// Connect connects to a wallet provider
	Connect(ctx context.Context, userID, providerName string, params map[string]interface{}) (*model.Wallet, error)

	// Disconnect disconnects from a wallet provider
	Disconnect(ctx context.Context, walletID string) error

	// Verify verifies a wallet connection using a signature
	Verify(ctx context.Context, walletID, message, signature string) (bool, error)

	// RefreshWallet refreshes a wallet's balance
	RefreshWallet(ctx context.Context, walletID string) (*model.Wallet, error)

	// GetProviders gets all available wallet providers
	GetProviders(ctx context.Context) ([]string, error)

	// GetProvidersByType gets all available wallet providers of a specific type
	GetProvidersByType(ctx context.Context, typ model.WalletType) ([]string, error)

	// IsValidAddress checks if an address is valid for a specific provider
	IsValidAddress(ctx context.Context, providerName, address string) (bool, error)
}

// walletConnectionService implements the WalletConnectionService interface
type walletConnectionService struct {
	providerRegistry *wallet.ProviderRegistry
	walletRepo       port.WalletRepository
	logger           *zerolog.Logger
}

// NewWalletConnectionService creates a new wallet connection service
func NewWalletConnectionService(
	providerRegistry *wallet.ProviderRegistry,
	walletRepo port.WalletRepository,
	logger *zerolog.Logger,
) WalletConnectionService {
	return &walletConnectionService{
		providerRegistry: providerRegistry,
		walletRepo:       walletRepo,
		logger:           logger,
	}
}

// Connect connects to a wallet provider
func (s *walletConnectionService) Connect(ctx context.Context, userID, providerName string, params map[string]interface{}) (*model.Wallet, error) {
	// Get provider
	provider, err := s.providerRegistry.GetProvider(providerName)
	if err != nil {
		s.logger.Error().Err(err).Str("provider", providerName).Msg("Failed to get provider")
		return nil, err
	}

	// Add userID to params
	params["user_id"] = userID

	// Connect to provider
	wallet, err := provider.Connect(ctx, params)
	if err != nil {
		s.logger.Error().Err(err).Str("provider", providerName).Msg("Failed to connect to provider")
		return nil, err
	}

	// Save wallet
	if err := s.walletRepo.Save(ctx, wallet); err != nil {
		s.logger.Error().Err(err).Str("id", wallet.ID).Msg("Failed to save wallet")
		return nil, err
	}

	return wallet, nil
}

// Disconnect disconnects from a wallet provider
func (s *walletConnectionService) Disconnect(ctx context.Context, walletID string) error {
	// Get wallet
	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get wallet")
		return err
	}
	if wallet == nil {
		return errors.New("wallet not found")
	}

	// Get provider
	var providerName string
	if wallet.Type == model.WalletTypeExchange {
		providerName = wallet.Exchange
	} else if wallet.Type == model.WalletTypeWeb3 {
		providerName = wallet.Metadata.Network
	} else {
		return errors.New("unsupported wallet type")
	}

	provider, err := s.providerRegistry.GetProvider(providerName)
	if err != nil {
		s.logger.Error().Err(err).Str("provider", providerName).Msg("Failed to get provider")
		return err
	}

	// Disconnect from provider
	if err := provider.Disconnect(ctx, walletID); err != nil {
		s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to disconnect from provider")
		return err
	}

	// Delete wallet
	if err := s.walletRepo.DeleteWallet(ctx, walletID); err != nil {
		s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to delete wallet")
		return err
	}

	return nil
}

// Verify verifies a wallet connection using a signature
func (s *walletConnectionService) Verify(ctx context.Context, walletID, message, signature string) (bool, error) {
	// Get wallet
	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get wallet")
		return false, err
	}
	if wallet == nil {
		return false, errors.New("wallet not found")
	}

	// Get provider
	var providerName string
	var address string
	if wallet.Type == model.WalletTypeExchange {
		providerName = wallet.Exchange
		address = wallet.Exchange // For exchange wallets, address is not applicable
	} else if wallet.Type == model.WalletTypeWeb3 {
		providerName = wallet.Metadata.Network
		address = wallet.Metadata.Address
	} else {
		return false, errors.New("unsupported wallet type")
	}

	provider, err := s.providerRegistry.GetProvider(providerName)
	if err != nil {
		s.logger.Error().Err(err).Str("provider", providerName).Msg("Failed to get provider")
		return false, err
	}

	// Verify signature
	verified, err := provider.Verify(ctx, address, message, signature)
	if err != nil {
		s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to verify signature")
		return false, err
	}

	// Update wallet verification status
	if verified {
		wallet.Status = model.WalletStatusVerified
		wallet.LastUpdated = time.Now()
		if err := s.walletRepo.Save(ctx, wallet); err != nil {
			s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to update wallet verification status")
			return false, err
		}
	}

	return verified, nil
}

// RefreshWallet refreshes a wallet's balance
func (s *walletConnectionService) RefreshWallet(ctx context.Context, walletID string) (*model.Wallet, error) {
	// Get wallet
	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get wallet")
		return nil, err
	}
	if wallet == nil {
		return nil, errors.New("wallet not found")
	}

	// Get provider
	var providerName string
	if wallet.Type == model.WalletTypeExchange {
		providerName = wallet.Exchange
	} else if wallet.Type == model.WalletTypeWeb3 {
		providerName = wallet.Metadata.Network
	} else {
		return nil, errors.New("unsupported wallet type")
	}

	provider, err := s.providerRegistry.GetProvider(providerName)
	if err != nil {
		s.logger.Error().Err(err).Str("provider", providerName).Msg("Failed to get provider")
		return nil, err
	}

	// Refresh wallet
	updatedWallet, err := provider.GetBalance(ctx, wallet)
	if err != nil {
		s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to refresh wallet")
		return nil, err
	}

	// Save wallet
	if err := s.walletRepo.Save(ctx, updatedWallet); err != nil {
		s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to save wallet")
		return nil, err
	}

	return updatedWallet, nil
}

// GetProviders gets all available wallet providers
func (s *walletConnectionService) GetProviders(ctx context.Context) ([]string, error) {
	providers := s.providerRegistry.GetAllProviders()
	providerNames := make([]string, len(providers))
	for i, provider := range providers {
		providerNames[i] = provider.GetName()
	}
	return providerNames, nil
}

// GetProvidersByType gets all available wallet providers of a specific type
func (s *walletConnectionService) GetProvidersByType(ctx context.Context, typ model.WalletType) ([]string, error) {
	providers, err := s.providerRegistry.GetProviderByType(typ)
	if err != nil {
		s.logger.Error().Err(err).Str("type", string(typ)).Msg("Failed to get providers by type")
		return nil, err
	}

	providerNames := make([]string, len(providers))
	for i, provider := range providers {
		providerNames[i] = provider.GetName()
	}
	return providerNames, nil
}

// IsValidAddress checks if an address is valid for a specific provider
func (s *walletConnectionService) IsValidAddress(ctx context.Context, providerName, address string) (bool, error) {
	// Get provider
	provider, err := s.providerRegistry.GetProvider(providerName)
	if err != nil {
		s.logger.Error().Err(err).Str("provider", providerName).Msg("Failed to get provider")
		return false, err
	}

	// Check if address is valid
	return provider.IsValidAddress(ctx, address)
}

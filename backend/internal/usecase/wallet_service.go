package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// WalletService defines the interface for wallet operations
type WalletService interface {
	// Fetch real account data from the exchange
	GetRealAccountData(ctx context.Context) (*model.Wallet, error)
	// Core wallet operations
	CreateWallet(ctx context.Context, userID, exchange string, walletType model.WalletType) (*model.Wallet, error)
	GetWallet(ctx context.Context, id string) (*model.Wallet, error)
	GetWalletByUserID(ctx context.Context, userID string) (*model.Wallet, error)
	GetWalletsByUserID(ctx context.Context, userID string) ([]*model.Wallet, error)
	UpdateWallet(ctx context.Context, wallet *model.Wallet) error
	DeleteWallet(ctx context.Context, id string) error

	// Wallet metadata operations
	SetWalletMetadata(ctx context.Context, id string, name, description string, tags []string) error
	SetPrimaryWallet(ctx context.Context, userID, walletID string) error
	AddCustomMetadata(ctx context.Context, id string, key, value string) error

	// Balance operations
	UpdateBalance(ctx context.Context, walletID string, asset model.Asset, free, locked, usdValue float64) error
	GetBalance(ctx context.Context, walletID string, asset model.Asset) (*model.Balance, error)
	GetBalanceHistory(ctx context.Context, userID string, asset model.Asset, from, to time.Time) ([]*model.BalanceHistory, error)

	// Wallet refresh operations
	RefreshWallet(ctx context.Context, id string) error
}

// walletService implements the WalletService interface
type walletService struct {
	walletRepo port.WalletRepository
	mexcClient port.MEXCClient
	logger     *zerolog.Logger
}

// NewWalletService creates a new wallet service
func NewWalletService(
	walletRepo port.WalletRepository,
	mexcClient port.MEXCClient,
	logger *zerolog.Logger,
) WalletService {
	return &walletService{
		walletRepo: walletRepo,
		mexcClient: mexcClient,
		logger:     logger,
	}
}

// CreateWallet creates a new wallet
func (s *walletService) CreateWallet(ctx context.Context, userID, exchange string, walletType model.WalletType) (*model.Wallet, error) {
	s.logger.Debug().
		Str("userID", userID).
		Str("exchange", exchange).
		Str("type", string(walletType)).
		Msg("Creating wallet")

	// Create wallet based on type
	var wallet *model.Wallet
	if walletType == model.WalletTypeExchange {
		wallet = model.NewExchangeWallet(userID, exchange)
	} else if walletType == model.WalletTypeWeb3 {
		wallet = model.NewWeb3Wallet(userID, "", "") // Network and address will be set later
	} else {
		wallet = model.NewWallet(userID)
		wallet.Type = walletType
	}

	// Validate wallet
	if err := wallet.Validate(); err != nil {
		s.logger.Error().Err(err).Msg("Invalid wallet")
		return nil, err
	}

	// Save wallet
	if err := s.walletRepo.Save(ctx, wallet); err != nil {
		s.logger.Error().Err(err).Msg("Failed to save wallet")
		return nil, err
	}

	return wallet, nil
}

// GetWallet gets a wallet by ID
func (s *walletService) GetWallet(ctx context.Context, id string) (*model.Wallet, error) {
	s.logger.Debug().Str("id", id).Msg("Getting wallet")
	return s.walletRepo.GetByID(ctx, id)
}

// GetWalletByUserID gets a wallet by user ID
func (s *walletService) GetWalletByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	s.logger.Debug().Str("userID", userID).Msg("Getting wallet by user ID")
	return s.walletRepo.GetByUserID(ctx, userID)
}

// GetWalletsByUserID gets all wallets for a user
func (s *walletService) GetWalletsByUserID(ctx context.Context, userID string) ([]*model.Wallet, error) {
	s.logger.Debug().Str("userID", userID).Msg("Getting wallets by user ID")
	return s.walletRepo.GetWalletsByUserID(ctx, userID)
}

// UpdateWallet updates a wallet
func (s *walletService) UpdateWallet(ctx context.Context, wallet *model.Wallet) error {
	s.logger.Debug().
		Str("id", wallet.ID).
		Str("userID", wallet.UserID).
		Msg("Updating wallet")

	// Validate wallet
	if err := wallet.Validate(); err != nil {
		s.logger.Error().Err(err).Msg("Invalid wallet")
		return err
	}

	// Save wallet
	return s.walletRepo.Save(ctx, wallet)
}

// DeleteWallet deletes a wallet
func (s *walletService) DeleteWallet(ctx context.Context, id string) error {
	s.logger.Debug().Str("id", id).Msg("Deleting wallet")
	return s.walletRepo.DeleteWallet(ctx, id)
}

// SetWalletMetadata sets the metadata for a wallet
func (s *walletService) SetWalletMetadata(ctx context.Context, id string, name, description string, tags []string) error {
	s.logger.Debug().
		Str("id", id).
		Str("name", name).
		Str("description", description).
		Strs("tags", tags).
		Msg("Setting wallet metadata")

	// Get wallet
	wallet, err := s.walletRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to get wallet")
		return err
	}
	if wallet == nil {
		return errors.New("wallet not found")
	}

	// Set metadata
	wallet.SetMetadata(name, description, tags)

	// Save wallet
	return s.walletRepo.Save(ctx, wallet)
}

// AddCustomMetadata adds a custom metadata key-value pair to a wallet
func (s *walletService) AddCustomMetadata(ctx context.Context, id string, key, value string) error {
	s.logger.Debug().
		Str("id", id).
		Str("key", key).
		Str("value", value).
		Msg("Adding custom metadata")

	// Get wallet
	wallet, err := s.walletRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to get wallet")
		return err
	}
	if wallet == nil {
		return errors.New("wallet not found")
	}

	// Add custom metadata
	wallet.AddCustomMetadata(key, value)

	// Save wallet
	return s.walletRepo.Save(ctx, wallet)
}

// SetPrimaryWallet sets a wallet as the primary wallet for a user
func (s *walletService) SetPrimaryWallet(ctx context.Context, userID, walletID string) error {
	s.logger.Debug().
		Str("userID", userID).
		Str("walletID", walletID).
		Msg("Setting primary wallet")

	// Get all wallets for the user
	wallets, err := s.walletRepo.GetWalletsByUserID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get wallets")
		return err
	}

	// Find the wallet to set as primary
	var primaryWallet *model.Wallet
	for _, wallet := range wallets {
		if wallet.ID == walletID {
			primaryWallet = wallet
		}
	}

	if primaryWallet == nil {
		return errors.New("wallet not found")
	}

	// Update all wallets
	for _, wallet := range wallets {
		isPrimary := wallet.ID == walletID
		wallet.SetPrimary(isPrimary)
		if err := s.walletRepo.Save(ctx, wallet); err != nil {
			s.logger.Error().Err(err).Str("id", wallet.ID).Msg("Failed to save wallet")
			return err
		}
	}

	return nil
}

// UpdateBalance updates a balance for a wallet
func (s *walletService) UpdateBalance(ctx context.Context, walletID string, asset model.Asset, free, locked, usdValue float64) error {
	s.logger.Debug().
		Str("walletID", walletID).
		Str("asset", string(asset)).
		Float64("free", free).
		Float64("locked", locked).
		Float64("usdValue", usdValue).
		Msg("Updating balance")

	// Get wallet
	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get wallet")
		return err
	}
	if wallet == nil {
		return errors.New("wallet not found")
	}

	// Update balance
	wallet.UpdateBalance(asset, free, locked, usdValue)

	// Save wallet
	return s.walletRepo.Save(ctx, wallet)
}

// GetBalance gets a balance for a wallet
func (s *walletService) GetBalance(ctx context.Context, walletID string, asset model.Asset) (*model.Balance, error) {
	s.logger.Debug().
		Str("walletID", walletID).
		Str("asset", string(asset)).
		Msg("Getting balance")

	// Get wallet
	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		s.logger.Error().Err(err).Str("id", walletID).Msg("Failed to get wallet")
		return nil, err
	}
	if wallet == nil {
		return nil, errors.New("wallet not found")
	}

	// Get balance
	return wallet.GetBalance(asset), nil
}

// GetBalanceHistory gets balance history for a user and asset
func (s *walletService) GetBalanceHistory(ctx context.Context, userID string, asset model.Asset, from, to time.Time) ([]*model.BalanceHistory, error) {
	s.logger.Debug().
		Str("userID", userID).
		Str("asset", string(asset)).
		Time("from", from).
		Time("to", to).
		Msg("Getting balance history")

	return s.walletRepo.GetBalanceHistory(ctx, userID, asset, from, to)
}

// GetRealAccountData fetches the real account data from the exchange via MEXCClient
func (s *walletService) GetRealAccountData(ctx context.Context) (*model.Wallet, error) {
	return s.mexcClient.GetAccount(ctx)
}

// RefreshWallet refreshes a wallet from the exchange
func (s *walletService) RefreshWallet(ctx context.Context, id string) error {
	s.logger.Debug().Str("id", id).Msg("Refreshing wallet")

	// Get wallet
	wallet, err := s.walletRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to get wallet")
		return err
	}
	if wallet == nil {
		return errors.New("wallet not found")
	}

	// Refresh wallet from exchange
	if wallet.Exchange == "MEXC" {
		// Get account from MEXC
		account, err := s.mexcClient.GetAccount(ctx)
		if err != nil {
			s.logger.Error().Err(err).Str("id", id).Msg("Failed to get account from MEXC")
			return err
		}

		// Update wallet balances
		for asset, balance := range account.Balances {
			wallet.UpdateBalance(asset, balance.Free, balance.Locked, balance.USDValue)
		}

		// Update wallet
		wallet.LastSyncAt = time.Now()
		wallet.LastUpdated = time.Now()

		// Save wallet
		return s.walletRepo.Save(ctx, wallet)
	}

	return errors.New("unsupported exchange")
}

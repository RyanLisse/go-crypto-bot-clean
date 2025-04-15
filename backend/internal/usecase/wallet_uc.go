package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// WalletUseCase defines the interface for wallet operations
type WalletUseCase interface {
	// Core wallet operations
	CreateWallet(ctx context.Context, wallet *model.Wallet) error
	GetWallet(ctx context.Context, id string) (*model.Wallet, error)
	GetWalletByUserID(ctx context.Context, userID string) (*model.Wallet, error)
	GetWalletsByUserID(ctx context.Context, userID string) ([]*model.Wallet, error)
	UpdateWallet(ctx context.Context, wallet *model.Wallet) error
	DeleteWallet(ctx context.Context, id string) error
	
	// Wallet metadata operations
	SetWalletMetadata(ctx context.Context, id string, name, description string, tags []string) error
	SetPrimaryWallet(ctx context.Context, userID, walletID string) error
	
	// Balance operations
	UpdateBalance(ctx context.Context, walletID string, asset model.Asset, free, locked, usdValue float64) error
	GetBalance(ctx context.Context, walletID string, asset model.Asset) (*model.Balance, error)
	GetBalanceHistory(ctx context.Context, userID string, asset model.Asset, from, to time.Time) ([]*model.BalanceHistory, error)
	
	// Wallet refresh operations
	RefreshWallet(ctx context.Context, id string) error
}

// walletUseCase implements the WalletUseCase interface
type walletUseCase struct {
	walletRepo port.WalletRepository
	mexcClient port.MEXCClient
	logger     *zerolog.Logger
}

// NewWalletUseCase creates a new wallet use case
func NewWalletUseCase(
	walletRepo port.WalletRepository,
	mexcClient port.MEXCClient,
	logger *zerolog.Logger,
) WalletUseCase {
	return &walletUseCase{
		walletRepo: walletRepo,
		mexcClient: mexcClient,
		logger:     logger,
	}
}

// CreateWallet creates a new wallet
func (uc *walletUseCase) CreateWallet(ctx context.Context, wallet *model.Wallet) error {
	// Validate wallet
	if err := wallet.Validate(); err != nil {
		uc.logger.Error().Err(err).Msg("Invalid wallet")
		return err
	}

	// Save wallet
	return uc.walletRepo.Save(ctx, wallet)
}

// GetWallet gets a wallet by ID
func (uc *walletUseCase) GetWallet(ctx context.Context, id string) (*model.Wallet, error) {
	return uc.walletRepo.GetByID(ctx, id)
}

// GetWalletByUserID gets a wallet by user ID
func (uc *walletUseCase) GetWalletByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	return uc.walletRepo.GetByUserID(ctx, userID)
}

// GetWalletsByUserID gets all wallets for a user
func (uc *walletUseCase) GetWalletsByUserID(ctx context.Context, userID string) ([]*model.Wallet, error) {
	return uc.walletRepo.GetWalletsByUserID(ctx, userID)
}

// UpdateWallet updates a wallet
func (uc *walletUseCase) UpdateWallet(ctx context.Context, wallet *model.Wallet) error {
	// Validate wallet
	if err := wallet.Validate(); err != nil {
		uc.logger.Error().Err(err).Msg("Invalid wallet")
		return err
	}

	// Get existing wallet
	existingWallet, err := uc.walletRepo.GetByID(ctx, wallet.ID)
	if err != nil {
		return err
	}
	if existingWallet == nil {
		return errors.New("wallet not found")
	}

	// Update wallet
	return uc.walletRepo.Save(ctx, wallet)
}

// DeleteWallet deletes a wallet
func (uc *walletUseCase) DeleteWallet(ctx context.Context, id string) error {
	return uc.walletRepo.DeleteWallet(ctx, id)
}

// SetWalletMetadata sets wallet metadata
func (uc *walletUseCase) SetWalletMetadata(ctx context.Context, id string, name, description string, tags []string) error {
	// Get wallet
	wallet, err := uc.walletRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if wallet == nil {
		return errors.New("wallet not found")
	}

	// Update metadata
	wallet.SetMetadata(name, description, tags)

	// Save wallet
	return uc.walletRepo.Save(ctx, wallet)
}

// SetPrimaryWallet sets a wallet as the primary wallet for a user
func (uc *walletUseCase) SetPrimaryWallet(ctx context.Context, userID, walletID string) error {
	// Get all wallets for the user
	wallets, err := uc.walletRepo.GetWalletsByUserID(ctx, userID)
	if err != nil {
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
		if err := uc.walletRepo.Save(ctx, wallet); err != nil {
			return err
		}
	}

	return nil
}

// UpdateBalance updates a balance for a wallet
func (uc *walletUseCase) UpdateBalance(ctx context.Context, walletID string, asset model.Asset, free, locked, usdValue float64) error {
	// Get wallet
	wallet, err := uc.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		return err
	}
	if wallet == nil {
		return errors.New("wallet not found")
	}

	// Update balance
	wallet.UpdateBalance(asset, free, locked, usdValue)

	// Save wallet
	return uc.walletRepo.Save(ctx, wallet)
}

// GetBalance gets a balance for a wallet
func (uc *walletUseCase) GetBalance(ctx context.Context, walletID string, asset model.Asset) (*model.Balance, error) {
	// Get wallet
	wallet, err := uc.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, errors.New("wallet not found")
	}

	// Get balance
	return wallet.GetBalance(asset), nil
}

// GetBalanceHistory gets balance history for a user and asset
func (uc *walletUseCase) GetBalanceHistory(ctx context.Context, userID string, asset model.Asset, from, to time.Time) ([]*model.BalanceHistory, error) {
	return uc.walletRepo.GetBalanceHistory(ctx, userID, asset, from, to)
}

// RefreshWallet refreshes a wallet from the exchange
func (uc *walletUseCase) RefreshWallet(ctx context.Context, id string) error {
	// Get wallet
	wallet, err := uc.walletRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if wallet == nil {
		return errors.New("wallet not found")
	}

	// Only refresh exchange wallets
	if wallet.Type != model.WalletTypeExchange {
		return errors.New("only exchange wallets can be refreshed")
	}

	// Refresh wallet from exchange
	if wallet.Exchange == "MEXC" {
		// Get account from MEXC
		account, err := uc.mexcClient.GetAccount(ctx)
		if err != nil {
			uc.logger.Error().Err(err).Str("id", id).Msg("Failed to get account from MEXC")
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
		return uc.walletRepo.Save(ctx, wallet)
	}

	return errors.New("unsupported exchange")
}

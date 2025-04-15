package usecase

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// AccountUsecase defines the interface for account operations
type AccountUsecase interface {
	GetWallet(ctx context.Context, userID string) (*model.Wallet, error)
	GetBalanceHistory(ctx context.Context, userID string, asset model.Asset, from, to time.Time) ([]*model.BalanceHistory, error)
	RefreshWallet(ctx context.Context, userID string) error
}

// accountUsecase implements the AccountUsecase interface
type accountUsecase struct {
	mexcClient port.MEXCClient // Changed mexcAPI to mexcClient
	walletRepo port.WalletRepository
	logger     zerolog.Logger
}

// NewAccountUsecase creates a new account usecase
func NewAccountUsecase(
	mexcClient port.MEXCClient, // Changed mexcAPI to mexcClient
	walletRepo port.WalletRepository,
	logger zerolog.Logger,
) AccountUsecase {
	return &accountUsecase{
		mexcClient: mexcClient, // Changed mexcAPI to mexcClient
		walletRepo: walletRepo,
		logger:     logger.With().Str("component", "account_usecase").Logger(),
	}
}

// GetWallet gets the user's wallet
func (uc *accountUsecase) GetWallet(ctx context.Context, userID string) (*model.Wallet, error) {
	// First try to get from DB
	wallet, err := uc.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get wallet from DB")
	}

	// If not found or error, get from API
	if wallet == nil {
		uc.logger.Debug().Str("userID", userID).Msg("No wallet found in DB, getting from API")
		wallet, err = uc.mexcClient.GetAccount(ctx) // Changed mexcAPI to mexcClient
		if err != nil {
			uc.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get wallet from API")
			return nil, err
		}

		// Set the user ID
		wallet.UserID = userID

		// Save to DB
		if err := uc.walletRepo.Save(ctx, wallet); err != nil {
			uc.logger.Error().Err(err).Str("userID", userID).Msg("Failed to save wallet to DB")
			// Continue anyway since we have the wallet data
		}
	}

	return wallet, nil
}

// GetBalanceHistory gets the user's balance history for a specific asset
func (uc *accountUsecase) GetBalanceHistory(ctx context.Context, userID string, asset model.Asset, from, to time.Time) ([]*model.BalanceHistory, error) {
	// Get from repository
	history, err := uc.walletRepo.GetBalanceHistory(ctx, userID, asset, from, to)
	if err != nil {
		uc.logger.Error().Err(err).Str("userID", userID).Str("asset", string(asset)).Time("from", from).Time("to", to).Msg("Failed to get balance history")
		return nil, err
	}

	return history, nil
}

// RefreshWallet refreshes the user's wallet from the exchange
func (uc *accountUsecase) RefreshWallet(ctx context.Context, userID string) error {
	// Get from API
	wallet, err := uc.mexcClient.GetAccount(ctx) // Changed mexcAPI to mexcClient
	if err != nil {
		uc.logger.Error().Err(err).Str("userID", userID).Msg("Failed to refresh wallet from API")
		return err
	}

	// Set the user ID
	wallet.UserID = userID

	// Save to DB
	if err := uc.walletRepo.Save(ctx, wallet); err != nil {
		uc.logger.Error().Err(err).Str("userID", userID).Msg("Failed to save refreshed wallet to DB")
		return err
	}

	return nil
}

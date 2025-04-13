package usecase

import (
	"context"
	"time"

	"github.com/neo/crypto-bot/internal/domain/model"
	"github.com/neo/crypto-bot/internal/domain/port"
	"github.com/rs/zerolog"
)

// AccountUsecase defines the interface for account operations
type AccountUsecase interface {
	GetWallet(ctx context.Context, userID string) (*model.Wallet, error)
	GetBalanceHistory(ctx context.Context, userID string, asset model.Asset, days int) ([]*model.BalanceHistory, error)
	RefreshWallet(ctx context.Context, userID string) error
}

// accountUsecase implements the AccountUsecase interface
type accountUsecase struct {
	mexcAPI    port.MexcAPI
	walletRepo port.WalletRepository
	logger     zerolog.Logger
}

// NewAccountUsecase creates a new account usecase
func NewAccountUsecase(
	mexcAPI port.MexcAPI,
	walletRepo port.WalletRepository,
	logger zerolog.Logger,
) AccountUsecase {
	return &accountUsecase{
		mexcAPI:    mexcAPI,
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

		// If not found or other error, try to get from API
		wallet, err = uc.mexcAPI.GetAccount(ctx)
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
func (uc *accountUsecase) GetBalanceHistory(ctx context.Context, userID string, asset model.Asset, days int) ([]*model.BalanceHistory, error) {
	// Calculate time range
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -days)

	// Get from repository
	history, err := uc.walletRepo.GetBalanceHistory(ctx, userID, asset, startTime, endTime)
	if err != nil {
		uc.logger.Error().Err(err).Str("userID", userID).Str("asset", string(asset)).Int("days", days).Msg("Failed to get balance history")
		return nil, err
	}

	return history, nil
}

// RefreshWallet refreshes the user's wallet from the exchange
func (uc *accountUsecase) RefreshWallet(ctx context.Context, userID string) error {
	// Get from API
	wallet, err := uc.mexcAPI.GetAccount(ctx)
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

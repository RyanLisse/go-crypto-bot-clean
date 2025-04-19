package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/google/uuid"
)

// WalletService defines the interface for wallet business operations
type WalletService interface {
	// CreateWeb3Wallet creates a new Web3 wallet for a user
	CreateWeb3Wallet(ctx context.Context, userID uuid.UUID, address, network string, isPrimary bool) (*model.Wallet, error)

	// CreateExchangeWallet creates a new exchange wallet for a user
	CreateExchangeWallet(ctx context.Context, userID uuid.UUID, exchangeID string, exchangeName string, isPrimary bool) (*model.Wallet, error)

	// GetWallet retrieves a wallet by its ID
	GetWallet(ctx context.Context, id uuid.UUID) (*model.Wallet, error)

	// GetUserWallets retrieves all wallets for a user
	GetUserWallets(ctx context.Context, userID uuid.UUID) ([]*model.Wallet, error)

	// GetPrimaryWallet retrieves a user's primary wallet
	GetPrimaryWallet(ctx context.Context, userID uuid.UUID) (*model.Wallet, error)

	// SetPrimaryWallet sets a wallet as the primary wallet for a user
	SetPrimaryWallet(ctx context.Context, walletID, userID uuid.UUID) error

	// UpdateBalances updates the balances for a wallet
	UpdateBalances(ctx context.Context, walletID uuid.UUID, balances []*model.Balance) error

	// VerifyWallet initiates the verification process for a wallet
	VerifyWallet(ctx context.Context, walletID uuid.UUID) error

	// CheckVerificationStatus checks the verification status of a wallet
	CheckVerificationStatus(ctx context.Context, walletID uuid.UUID) (model.VerificationStatus, error)

	// DeleteWallet removes a wallet
	DeleteWallet(ctx context.Context, walletID, userID uuid.UUID) error
}

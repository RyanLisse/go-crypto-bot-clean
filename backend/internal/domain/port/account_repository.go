package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// AccountRepository defines the interface for account-related repository operations
type AccountRepository interface {
	// GetWallet returns a wallet by user ID
	GetWallet(ctx context.Context, userID string) (*model.Wallet, error)

	// SaveWallet saves a wallet
	SaveWallet(ctx context.Context, wallet *model.Wallet) error

	// GetBalanceHistory returns balance history for a user and asset over a number of days
	GetBalanceHistory(ctx context.Context, userID string, asset model.Asset, days int) ([]*model.BalanceHistory, error)

	// GetTransactions returns transactions for a user with pagination
	GetTransactions(ctx context.Context, userID string, limit, offset int) ([]model.Transaction, error)
}

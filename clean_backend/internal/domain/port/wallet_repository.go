package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/google/uuid"
)

// WalletRepository defines the interface for wallet persistence operations
type WalletRepository interface {
	// Create/Update Operations
	Create(ctx context.Context, w *model.Wallet) error
	Update(ctx context.Context, w *model.Wallet) error
	UpdateStatus(ctx context.Context, walletID uuid.UUID, status model.WalletStatus) error
	UpdateBalance(ctx context.Context, walletID uuid.UUID, balances []model.Balance) error
	SetPrimary(ctx context.Context, userID uuid.UUID, walletID uuid.UUID, walletType model.WalletType) error

	// Query Operations
	GetByID(ctx context.Context, walletID uuid.UUID) (*model.Wallet, error)
	GetByAddress(ctx context.Context, address string) (*model.Wallet, error)
	GetByExchangeID(ctx context.Context, exchangeID string) (*model.Wallet, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Wallet, error)
	FindPrimary(ctx context.Context, userID uuid.UUID, walletType model.WalletType) (*model.Wallet, error)
	GetBalance(ctx context.Context, walletID uuid.UUID, asset model.Asset) (*model.Balance, error)
	GetTotalBalance(ctx context.Context, userID uuid.UUID, asset model.Asset) (float64, error)

	// Delete Operations
	Delete(ctx context.Context, walletID uuid.UUID) error
}

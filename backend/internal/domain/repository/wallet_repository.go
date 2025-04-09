package repository

import (
	"context"
	"errors"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

var (
	// ErrNotFound is returned when a requested entity is not found
	ErrNotFound = errors.New("entity not found")
)

// WalletRepository defines the interface for wallet data access
type WalletRepository interface {
	// GetWallet retrieves the wallet from the database
	GetWallet(ctx context.Context) (*models.Wallet, error)
	
	// SaveWallet saves the wallet to the database
	SaveWallet(ctx context.Context, wallet *models.Wallet) (*models.Wallet, error)
}

package port

import (
	"context"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// CoinbaseClient defines the interface for interacting with Coinbase API
// This is a scaffold for Coinbase integration

type CoinbaseClient interface {
	GetAccount(ctx context.Context) (*model.Wallet, error)
}

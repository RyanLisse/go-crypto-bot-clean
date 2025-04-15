package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// WalletProvider defines the interface for wallet providers
type WalletProvider interface {
	// GetName returns the name of the wallet provider
	GetName() string

	// GetType returns the type of wallet provider
	GetType() model.WalletType

	// Connect connects to the wallet provider
	Connect(ctx context.Context, params map[string]interface{}) (*model.Wallet, error)

	// Disconnect disconnects from the wallet provider
	Disconnect(ctx context.Context, walletID string) error

	// Verify verifies a wallet connection using a signature
	Verify(ctx context.Context, address string, message string, signature string) (bool, error)

	// GetBalance gets the balance for a wallet
	GetBalance(ctx context.Context, wallet *model.Wallet) (*model.Wallet, error)

	// IsValidAddress checks if an address is valid for this provider
	IsValidAddress(ctx context.Context, address string) (bool, error)
}

// ExchangeWalletProvider defines the interface for exchange wallet providers
type ExchangeWalletProvider interface {
	WalletProvider

	// SetAPICredentials sets the API credentials for the exchange
	SetAPICredentials(ctx context.Context, apiKey, apiSecret string) error
}

// Web3WalletProvider defines the interface for Web3 wallet providers
type Web3WalletProvider interface {
	WalletProvider

	// GetChainID returns the chain ID for the provider
	GetChainID() int64

	// GetNetwork returns the network for the provider
	GetNetwork() string

	// SignMessage signs a message with the wallet's private key
	SignMessage(ctx context.Context, message string) (string, error)
}

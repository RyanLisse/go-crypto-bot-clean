package usecase

import "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"

// NewCoinUseCase defines methods for detecting and managing new coin listings
type NewCoinUseCase interface {
	// DetectNewCoins checks for newly listed coins on MEXC
	DetectNewCoins() error

	// UpdateCoinStatus updates a coin's status and creates an event
	UpdateCoinStatus(coinID string, newStatus model.CoinStatus) error

	// SubscribeToEvents subscribes to new coin events
	SubscribeToEvents(handler func(*model.CoinEvent)) error

	// GetCoinDetails retrieves details for a specific coin
	GetCoinDetails(coinID string) (*model.Coin, error)
}

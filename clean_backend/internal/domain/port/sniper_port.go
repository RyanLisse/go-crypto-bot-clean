package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
)

// SniperShotServicePort defines the interface for sniper shot operations,
// using domain models for requests and results.
type SniperShotServicePort interface {
	// ExecuteSniper executes a sniper shot trade according to the provided request.
	ExecuteSniper(ctx context.Context, req *model.SniperShotRequest) (*model.SniperShotResult, error)

	// CancelSniper attempts to cancel a previously executed sniper shot order.
	CancelSniper(ctx context.Context, symbol, orderID string) error

	// GetSniperOrderStatus retrieves the current status of a sniper shot order.
	// It should return the order details using the domain model.
	GetSniperOrderStatus(ctx context.Context, symbol, orderID string) (*model.Order, error)
}

// PriceChecker defines an interface for getting the current price of a symbol.
// This is often needed by services implementing SniperShotServicePort.
type PriceChecker interface {
	// GetCurrentPrice returns the current price for a given symbol.
	GetCurrentPrice(ctx context.Context, symbol string) (float64, error)
}

package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// SniperConfig defines configuration parameters for the sniper service
type SniperConfig struct {
	// MaxBuyAmount is the maximum amount in quote currency to spend on a new listing
	MaxBuyAmount float64

	// MaxPricePerToken is the maximum price per token to pay
	MaxPricePerToken float64

	// EnablePartialFills allows the sniper to execute partial fills if full amount can't be filled
	EnablePartialFills bool

	// MaxSlippagePercent is the maximum allowed slippage percentage
	MaxSlippagePercent float64

	// BypassRiskChecks determines whether to bypass risk checks for faster execution
	BypassRiskChecks bool

	// PreferredOrderType specifies the preferred order type (market or limit)
	PreferredOrderType model.OrderType

	// MaxConcurrentOrders is the maximum number of concurrent orders to place
	MaxConcurrentOrders int

	// RetryAttempts is the number of retry attempts for failed orders
	RetryAttempts int

	// RetryDelayMs is the delay between retries in milliseconds
	RetryDelayMs int

	// EnableTakeProfit enables automatic take-profit orders
	EnableTakeProfit bool

	// TakeProfitPercent is the percentage gain at which to take profit
	TakeProfitPercent float64

	// EnableStopLoss enables automatic stop-loss orders
	EnableStopLoss bool

	// StopLossPercent is the percentage loss at which to stop loss
	StopLossPercent float64

	// PriceCacheExpiryMs is the expiry time for price cache entries in milliseconds
	PriceCacheExpiryMs int

	// RateLimitPerSec is the maximum number of API calls per second
	RateLimitPerSec int
}

// SniperService defines the interface for high-speed trading on new listings
type SniperService interface {
	// ExecuteSnipe executes a high-speed buy on a newly listed token
	ExecuteSnipe(ctx context.Context, symbol string) (*model.Order, error)

	// ExecuteSnipeWithConfig executes a high-speed buy with custom configuration
	ExecuteSnipeWithConfig(ctx context.Context, symbol string, config *SniperConfig) (*model.Order, error)

	// PrevalidateSymbol checks if a symbol is valid for sniping without executing a trade
	PrevalidateSymbol(ctx context.Context, symbol string) (bool, error)

	// GetConfig returns the current sniper configuration
	GetConfig() *SniperConfig

	// UpdateConfig updates the sniper configuration
	UpdateConfig(config *SniperConfig) error

	// GetStatus returns the current status of the sniper service
	GetStatus() string

	// Start starts the sniper service
	Start() error

	// Stop stops the sniper service
	Stop() error
}

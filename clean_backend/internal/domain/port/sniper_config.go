package port

import (
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
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

package model

import (
	"context"
	"time"
)

// TriggerType defines the different conditions that can trigger an AutoBuyRule
type TriggerType string

// Available trigger types
const (
	TriggerTypeNewListing    TriggerType = "NEW_LISTING"
	TriggerTypePriceIncrease TriggerType = "PRICE_INCREASE"
	TriggerTypePriceDecrease TriggerType = "PRICE_DECREASE"
	TriggerTypeVolumeSurge   TriggerType = "VOLUME_SURGE"
	TriggerTypePriceBelow    TriggerType = "PRICE_BELOW"
	TriggerTypePriceAbove    TriggerType = "PRICE_ABOVE"
	TriggerTypePercentDrop   TriggerType = "PERCENT_DROP"
	TriggerTypePercentRise   TriggerType = "PERCENT_RISE"
)

// AutoBuyRule defines a rule for automatic buying based on specific conditions
type AutoBuyRule struct {
	ID                  string
	UserID              string      // ID of the user who owns the rule
	Name                string      `json:"name"` // User-defined name for the rule
	Symbol              string      // Trading symbol (e.g., BTCUSDT) or "*" for any new listing
	IsEnabled           bool        `json:"is_enabled"`        // Whether the rule is active (renamed from IsActive for consistency)
	TriggerType         TriggerType `json:"trigger_type"`      // Condition to trigger the buy
	TriggerValue        float64     `json:"trigger_value"`     // Value associated with the trigger type (price, percentage, volume)
	QuoteAsset          string      `json:"quote_asset"`       // e.g., "USDT" - only buy pairs with this quote asset (for new listings)
	BuyAmountQuote      float64     `json:"buy_amount_quote"`  // Amount of quote asset to spend per buy (supersedes FixedAmount/PercentageAmount)
	MaxBuyPrice         *float64    `json:"max_buy_price"`     // Optional: Maximum price to buy at
	MinBaseAssetVolume  *float64    `json:"min_base_volume"`   // Optional: Minimum 24h volume for the base asset
	MinQuoteAssetVolume *float64    `json:"min_quote_volume"`  // Optional: Minimum 24h volume for the quote asset
	AllowPreTrading     bool        `json:"allow_pre_trading"` // Whether to buy during pre-trading phase if possible (for new listings)
	CooldownMinutes     int         `json:"cooldown_minutes"`  // Minimum minutes between triggers for the same rule/symbol
	OrderType           OrderType   `json:"order_type"`        // Type of order to place (MARKET, LIMIT)
	EnableRiskCheck     bool        `json:"enable_risk_check"` // Whether to perform risk checks before buying
	ExecutionCount      int         `json:"execution_count"`   // How many times this rule has been successfully executed
	LastTriggered       *time.Time  `json:"last_triggered"`    // Timestamp of the last successful execution
	LastPrice           float64     `json:"last_price"`        // Price at the time of the last trigger
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`
}

// AutoBuyExecution records an instance where an AutoBuyRule was triggered and an order placed
type AutoBuyExecution struct {
	ID        string
	RuleID    string
	UserID    string
	Symbol    string
	OrderID   string  // The ID of the order placed by the trade service
	Price     float64 // Execution price
	Quantity  float64 // Executed quantity
	Amount    float64 // Total amount in quote currency
	Timestamp time.Time
}

// AutoBuyRuleRepository defines the interface for AutoBuyRule persistence
type AutoBuyRuleRepository interface {
	Create(ctx context.Context, rule *AutoBuyRule) error
	Update(ctx context.Context, rule *AutoBuyRule) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*AutoBuyRule, error)
	GetByUserID(ctx context.Context, userID string) ([]*AutoBuyRule, error)
	GetBySymbol(ctx context.Context, symbol string) ([]*AutoBuyRule, error) // Can be multiple rules per symbol
	GetActive(ctx context.Context) ([]*AutoBuyRule, error)                  // Get all enabled rules
	ListAll(ctx context.Context, limit, offset int) ([]*AutoBuyRule, error)
	// GetByTriggerType is likely not needed if evaluation logic fetches active rules first
	// GetByTriggerType(ctx context.Context, triggerType TriggerType) ([]*AutoBuyRule, error)
}

// AutoBuyExecutionRepository defines the interface for AutoBuyExecution persistence
type AutoBuyExecutionRepository interface {
	Create(ctx context.Context, execution *AutoBuyExecution) error
	GetByID(ctx context.Context, id string) (*AutoBuyExecution, error)
	GetByRuleID(ctx context.Context, ruleID string, limit, offset int) ([]*AutoBuyExecution, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*AutoBuyExecution, error)
	GetBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*AutoBuyExecution, error)
	GetByTimeRange(ctx context.Context, from, to time.Time, limit, offset int) ([]*AutoBuyExecution, error)
}

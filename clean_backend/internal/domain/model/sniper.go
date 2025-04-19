package model

import "time"

// ComparisonType defines how to compare prices for triggers
type ComparisonType string

const (
	// Above triggers when price goes above threshold
	Above ComparisonType = "ABOVE"
	// Below triggers when price goes below threshold
	Below ComparisonType = "BELOW"
	// Equal triggers when price equals threshold
	Equal ComparisonType = "EQUAL" // Added for completeness
	// NotEqual triggers when price is not equal to threshold
	NotEqual ComparisonType = "NOTEQUAL" // Added for completeness
	// GreaterOrEqual triggers when price is >= threshold
	GreaterOrEqual ComparisonType = "GE" // Using common abbreviations
	// LessOrEqual triggers when price is <= threshold
	LessOrEqual ComparisonType = "LE" // Using common abbreviations
)

// TriggerCondition defines conditions that must be met before executing a sniper shot
type TriggerCondition struct {
	TargetPrice     float64               // The price threshold to compare against
	Operator        string                // Comparison operator (e.g., ">", "<=", "==") - Kept string for flexibility from service
	Comparison      ComparisonType        // Enum based comparison type (preferred)
	MaxTimeoutSecs  int                   // Maximum seconds to wait for the condition
	CheckIntervalMs int                   // Milliseconds between price checks
	PriceBufferPct  float64               // Optional price buffer percentage for limit orders
	Callbacks       []func(price float64) // Optional functions to call when condition met
	// Consider adding Symbol here if condition is tied to a specific symbol price check
}

// SniperShotRequest represents a request for a sniper shot trade, using domain types
type SniperShotRequest struct {
	UserID    string            // ID of the user making the request
	Symbol    string            // Symbol to trade (e.g., "BTCUSDT")
	Side      OrderSide         // BUY or SELL (from model.OrderSide)
	Quantity  float64           // Amount to buy or sell
	Price     float64           // Price limit (0 for market orders)
	Type      OrderType         // LIMIT or MARKET (from model.OrderType)
	TimeLimit time.Duration     // Maximum time to try executing the order
	Condition *TriggerCondition // Optional condition that must be met
}

// SniperShotResult represents the result of a sniper shot execution, using domain types
type SniperShotResult struct {
	Success   bool          // Whether the shot was successful
	Order     *Order        // Order details if successful (from model.Order)
	Error     error         // Error details if unsuccessful
	Timestamp time.Time     // When the shot was executed
	Latency   time.Duration // How long it took to execute
}

package account

import (
	"time"
)

// Balance represents an account balance
type Balance struct {
	Fiat      float64             `json:"fiat"`
	Available map[string]float64  `json:"available"`
}

// AssetBalance represents the balance of a single asset
type AssetBalance struct {
	Asset  string  `json:"asset"`
	Free   float64 `json:"free"`
	Locked float64 `json:"locked"`
	Total  float64 `json:"total"`
}

// Wallet represents a wallet with multiple asset balances
type Wallet struct {
	Balances  map[string]*AssetBalance `json:"balances"`
	UpdatedAt time.Time                `json:"updatedAt"`
}

// Transaction represents a transaction
type Transaction struct {
	ID        int64     `json:"id"`
	Amount    float64   `json:"amount"`
	Balance   float64   `json:"balance"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}

// BalanceSummary represents a summary of the account balance
type BalanceSummary struct {
	CurrentBalance   float64 `json:"currentBalance"`
	Deposits         float64 `json:"deposits"`
	Withdrawals      float64 `json:"withdrawals"`
	NetChange        float64 `json:"netChange"`
	TransactionCount int     `json:"transactionCount"`
	Period           int     `json:"period"` // days
}

// TransactionAnalysis represents an analysis of transactions
type TransactionAnalysis struct {
	TotalDeposits    float64 `json:"totalDeposits"`
	TotalWithdrawals float64 `json:"totalWithdrawals"`
	NetChange        float64 `json:"netChange"`
	TransactionCount int     `json:"transactionCount"`
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
}

// PositionRisk represents the risk of a position
type PositionRisk struct {
	Symbol           string  `json:"symbol"`
	PositionAmount   float64 `json:"positionAmount"`
	EntryPrice       float64 `json:"entryPrice"`
	MarkPrice        float64 `json:"markPrice"`
	UnrealizedProfit float64 `json:"unrealizedProfit"`
	LiquidationPrice float64 `json:"liquidationPrice"`
	Leverage         float64 `json:"leverage"`
	MaxNotionalValue float64 `json:"maxNotionalValue"`
	MarginType       string  `json:"marginType"`
	PositionSide     string  `json:"positionSide"`
	UpdateTime       int64   `json:"updateTime"`
}

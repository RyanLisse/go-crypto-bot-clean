package models

import (
	"time"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID        int64     `json:"id" db:"id"`
	Amount    float64   `json:"amount" db:"amount"`
	Balance   float64   `json:"balance" db:"balance"`
	Reason    string    `json:"reason" db:"reason"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

// BalanceSummary provides an overview of wallet activity
type BalanceSummary struct {
	CurrentBalance   float64   `json:"current_balance"`
	Deposits         float64   `json:"deposits"`
	Withdrawals      float64   `json:"withdrawals"`
	NetChange        float64   `json:"net_change"`
	TransactionCount int       `json:"transaction_count"`
	Period           int       `json:"period_days"`
	GeneratedAt      time.Time `json:"generated_at"`
}

// TransactionAnalysis provides analysis of transaction history
type TransactionAnalysis struct {
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	TotalCount  int       `json:"total_count"`
	BuyCount    int       `json:"buy_count"`
	SellCount   int       `json:"sell_count"`
	TotalVolume float64   `json:"total_volume"`
	BuyVolume   float64   `json:"buy_volume"`
	SellVolume  float64   `json:"sell_volume"`
}

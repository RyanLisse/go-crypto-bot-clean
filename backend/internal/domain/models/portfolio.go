package models

import "time"

// Portfolio represents a user's portfolio
type Portfolio struct {
	UserID    int       `json:"user_id"`
	Balance   float64   `json:"balance"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PortfolioSummary represents a summary of a user's portfolio
type PortfolioSummary struct {
	TotalValue         float64                `json:"total_value"`
	AvailableBalance   float64                `json:"available_balance"`
	TotalProfitLoss    float64                `json:"total_profit_loss"`
	ProfitLossPercent  float64                `json:"profit_loss_percent"`
	ActivePositions    int                    `json:"active_positions"`
	OpenOrders         int                    `json:"open_orders"`
	LastUpdated        time.Time              `json:"last_updated"`
	TopHoldings        []PortfolioHolding     `json:"top_holdings"`
	RecentTransactions []PortfolioTransaction `json:"recent_transactions"`
}

// PortfolioHolding represents a single holding in the portfolio
type PortfolioHolding struct {
	Symbol    string  `json:"symbol"`
	Quantity  float64 `json:"quantity"`
	Value     float64 `json:"value"`
	PriceUSD  float64 `json:"price_usd"`
	Change24h float64 `json:"change_24h"`
}

// PortfolioTransaction represents a recent transaction in the portfolio
type PortfolioTransaction struct {
	Symbol    string    `json:"symbol"`
	Type      string    `json:"type"` // "BUY" or "SELL"
	Quantity  float64   `json:"quantity"`
	Price     float64   `json:"price"`
	Total     float64   `json:"total"`
	Timestamp time.Time `json:"timestamp"`
}

// Asset represents a single asset in a portfolio
type Asset struct {
	Symbol    string  `json:"symbol"`
	Quantity  float64 `json:"quantity"`
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"`
}

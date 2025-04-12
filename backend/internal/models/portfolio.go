package models

import (
	"time"
)

// Portfolio represents a user's portfolio
type Portfolio struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	UserID     string    `json:"user_id" gorm:"index"`
	TotalValue float64   `json:"total_value"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Position represents a position in the portfolio
type Position struct {
	ID           string  `json:"id" gorm:"primaryKey"`
	PortfolioID  string  `json:"portfolio_id" gorm:"index"`
	Symbol       string  `json:"symbol"`
	Quantity     float64 `json:"quantity"`
	AveragePrice float64 `json:"average_price"`
	CurrentPrice float64 `json:"current_price"`
	ProfitLoss   float64 `json:"profit_loss"`
}

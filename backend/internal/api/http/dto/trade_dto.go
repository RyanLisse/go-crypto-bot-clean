package dto

import (
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// TradeResponse represents the response body for trade operations
type TradeResponse struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	Symbol    string    `json:"symbol"`
	Side      string    `json:"side"`
	Quantity  float64   `json:"quantity"`
	Price     float64   `json:"price"`
	Fee       float64   `json:"fee"`
	TradeTime time.Time `json:"trade_time"`
}

// FromModel converts a domain Trade model to a TradeResponse
func TradeResponseFromModel(trade *models.Trade) *TradeResponse {
	return &TradeResponse{
		ID:        trade.ID,
		OrderID:   trade.OrderID,
		Symbol:    trade.Symbol,
		Side:      trade.Side,
		Quantity:  trade.Amount, // Use Amount instead of Quantity
		Price:     trade.Price,
		Fee:       0, // Fee is not in the model, set to 0
		TradeTime: trade.TradeTime,
	}
}

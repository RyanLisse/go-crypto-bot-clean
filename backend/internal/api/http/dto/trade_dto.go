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
	var side string
	if trade.IsBuyer {
		side = "buy"
	} else {
		side = "sell"
	}

	return &TradeResponse{
		ID:        trade.ID,
		OrderID:   trade.OrderID,
		Symbol:    trade.Symbol,
		Side:      side,
		Quantity:  trade.Quantity,
		Price:     trade.Price,
		Fee:       trade.Fee,
		TradeTime: trade.TradeTime,
	}
}

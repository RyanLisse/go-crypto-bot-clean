package dto

import (
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// PositionResponse represents the response body for position operations
type PositionResponse struct {
	ID        string    `json:"id"`
	Symbol    string    `json:"symbol"`
	Side      string    `json:"side"`
	Quantity  float64   `json:"quantity"`
	EntryPrice float64  `json:"entry_price"`
	ExitPrice  float64  `json:"exit_price,omitempty"`
	Status    string    `json:"status"`
	PnL       float64   `json:"pnl"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FromModel converts a domain Position model to a PositionResponse
func PositionResponseFromModel(position *models.Position) *PositionResponse {
	var side string
	if position.Side == models.PositionSideLong {
		side = "long"
	} else {
		side = "short"
	}

	return &PositionResponse{
		ID:         position.ID,
		Symbol:     position.Symbol,
		Side:       side,
		Quantity:   position.Quantity,
		EntryPrice: position.EntryPrice,
		ExitPrice:  position.ExitPrice,
		Status:     string(position.Status),
		PnL:        position.PnL,
		CreatedAt:  position.CreatedAt,
		UpdatedAt:  position.UpdatedAt,
	}
}

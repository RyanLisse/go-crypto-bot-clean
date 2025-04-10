package dto

import (
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// OrderRequest represents the request body for order operations
type OrderRequest struct {
	Symbol   string  `json:"symbol"`
	Side     string  `json:"side"`
	Type     string  `json:"type"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price,omitempty"`
}

// OrderResponse represents the response body for order operations
type OrderResponse struct {
	ID        string    `json:"id"`
	Symbol    string    `json:"symbol"`
	Side      string    `json:"side"`
	Type      string    `json:"type"`
	Quantity  float64   `json:"quantity"`
	Price     float64   `json:"price"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToModel converts an OrderRequest to a domain Order model
func (r *OrderRequest) ToModel() *models.Order {
	var side models.OrderSide
	if r.Side == "buy" {
		side = models.OrderSideBuy
	} else {
		side = models.OrderSideSell
	}

	var orderType models.OrderType
	switch r.Type {
	case "market":
		orderType = models.OrderTypeMarket
	case "limit":
		orderType = models.OrderTypeLimit
	default:
		orderType = models.OrderTypeMarket
	}

	return &models.Order{
		Symbol:   r.Symbol,
		Side:     side,
		Type:     orderType,
		Quantity: r.Quantity,
		Price:    r.Price,
		Status:   models.OrderStatusNew,
	}
}

// FromModel converts a domain Order model to an OrderResponse
func OrderResponseFromModel(order *models.Order) *OrderResponse {
	var side string
	if order.Side == models.OrderSideBuy {
		side = "buy"
	} else {
		side = "sell"
	}

	var orderType string
	switch order.Type {
	case models.OrderTypeMarket:
		orderType = "market"
	case models.OrderTypeLimit:
		orderType = "limit"
	default:
		orderType = "unknown"
	}

	return &OrderResponse{
		ID:        order.ID,
		Symbol:    order.Symbol,
		Side:      side,
		Type:      orderType,
		Quantity:  order.Quantity,
		Price:     order.Price,
		Status:    string(order.Status),
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}

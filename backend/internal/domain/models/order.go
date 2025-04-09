package models

import "time"

// Order represents a trading order
type Order struct {
	ID        string      `json:"id"`
	OrderID   string      `json:"order_id"`  // Exchange-generated order ID
	ClientID  string      `json:"client_id"` // Client-generated order ID
	Symbol    string      `json:"symbol"`
	Side      OrderSide   `json:"side"`
	Type      OrderType   `json:"type"`
	Quantity  float64     `json:"quantity"`
	FilledQty float64     `json:"filled_qty"` // Amount that has been filled
	Price     float64     `json:"price"`
	Status    OrderStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Time      time.Time   `json:"time"` // Alias for CreatedAt for backward compatibility
}

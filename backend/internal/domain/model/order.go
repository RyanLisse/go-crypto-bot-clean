package model

import (
	"time"
)

// OrderSide represents the side of an order (buy or sell)
type OrderSide string

// OrderType represents the type of an order
type OrderType string

// OrderStatus represents the status of an order
type OrderStatus string

// TimeInForce represents how long an order will remain active
type TimeInForce string

// Order side constants
const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

// Order type constants
const (
	OrderTypeMarket OrderType = "MARKET"
	OrderTypeLimit  OrderType = "LIMIT"
)

// Order status constants
const (
	OrderStatusNew             OrderStatus = "NEW"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusCanceled        OrderStatus = "CANCELED"
	OrderStatusRejected        OrderStatus = "REJECTED"
)

// TimeInForce constants
const (
	TimeInForceGTC TimeInForce = "GTC" // Good Till Cancel
	TimeInForceIOC TimeInForce = "IOC" // Immediate Or Cancel
	TimeInForceFOK TimeInForce = "FOK" // Fill Or Kill
)

// Order represents a trading order
type Order struct {
	ID            string      `json:"id"`
	OrderID       string      `json:"order_id"` // Exchange-specific order ID
	ClientOrderID string      `json:"clientOrderId"`
	Symbol        string      `json:"symbol"`
	Side          OrderSide   `json:"side"`
	Type          OrderType   `json:"type"`
	Status        OrderStatus `json:"status"`
	TimeInForce   TimeInForce `json:"timeInForce,omitempty"`
	Price         float64     `json:"price"`
	Quantity      float64     `json:"quantity"`
	ExecutedQty   float64     `json:"executedQty"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
}

// OrderRequest represents the data required to create a new order
type OrderRequest struct {
	Symbol        string    `json:"symbol" binding:"required"`
	Side          OrderSide `json:"side" binding:"required,oneof=BUY SELL"`
	Type          OrderType `json:"type" binding:"required,oneof=MARKET LIMIT"`
	Quantity      float64   `json:"quantity" binding:"required,gt=0"`
	Price         float64   `json:"price"`
	ClientOrderID string    `json:"clientOrderId"`
}

// OrderResponse represents the data returned after creating/querying an order
type OrderResponse struct {
	Order
	AvgPrice float64 `json:"avgPrice"`
}

// IsComplete returns true if the order is in a terminal state
func (o *Order) IsComplete() bool {
	return o.Status == OrderStatusFilled ||
		o.Status == OrderStatusCanceled ||
		o.Status == OrderStatusRejected
}

// RemainingQuantity returns the quantity that has not been executed yet
func (o *Order) RemainingQuantity() float64 {
	return o.Quantity - o.ExecutedQty
}

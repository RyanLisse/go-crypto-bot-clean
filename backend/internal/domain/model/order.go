package model

import "time"

// OrderSide represents the side of an order (BUY or SELL)
type OrderSide string

const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

// OrderType represents the type of an order (LIMIT, MARKET, etc.)
type OrderType string

const (
	OrderTypeLimit  OrderType = "LIMIT"
	OrderTypeMarket OrderType = "MARKET"
	// Add other types like STOP_LOSS, TAKE_PROFIT, STOP_LOSS_LIMIT, etc. if needed
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusNew             OrderStatus = "NEW"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusCanceled        OrderStatus = "CANCELED"
	OrderStatusPendingCancel   OrderStatus = "PENDING_CANCEL" // Currently unused, but potentially useful
	OrderStatusRejected        OrderStatus = "REJECTED"
	OrderStatusExpired         OrderStatus = "EXPIRED"
)

// TimeInForce represents how long an order remains active before cancellation
type TimeInForce string

const (
	TimeInForceGTC TimeInForce = "GTC" // Good Til Canceled
	TimeInForceIOC TimeInForce = "IOC" // Immediate or Cancel
	TimeInForceFOK TimeInForce = "FOK" // Fill or Kill
)

// Order represents a trading order
type Order struct {
	ID              string      `json:"id"`               // Unique identifier for the order in our system
	OrderID         string      `json:"order_id"`         // Exchange's order ID
	ClientOrderID   string      `json:"client_order_id"`  // Optional client-provided ID
	UserID          string      `json:"user_id"`          // User who placed the order
	Symbol          string      `json:"symbol"`           // Trading pair (e.g., BTCUSDT)
	Side            OrderSide   `json:"side"`             // BUY or SELL
	Type            OrderType   `json:"type"`             // LIMIT, MARKET, etc.
	Status          OrderStatus `json:"status"`           // Current status of the order
	Price           float64     `json:"price"`            // Order price (0 for MARKET orders)
	Quantity        float64     `json:"quantity"`         // Original order quantity
	ExecutedQty     float64     `json:"executed_qty"`     // Quantity that has been filled
	AvgFillPrice    float64     `json:"avg_fill_price"`   // Average price of filled quantity
	Commission      float64     `json:"commission"`       // Trading commission paid
	CommissionAsset string      `json:"commission_asset"` // Asset used for commission
	TimeInForce     TimeInForce `json:"time_in_force"`    // Order duration policy
	CreatedAt       time.Time   `json:"created_at"`       // Time order was created in our system
	UpdatedAt       time.Time   `json:"updated_at"`       // Last time order was updated in our system
	Exchange        string      `json:"exchange"`         // Exchange where the order was placed
}

// IsComplete returns true if the order is in a terminal state (filled, canceled, rejected, or expired)
func (o *Order) IsComplete() bool {
	switch o.Status {
	case OrderStatusFilled, OrderStatusCanceled, OrderStatusRejected, OrderStatusExpired:
		return true
	default:
		return false
	}
}

// OrderRequest represents the data needed to place a new order
type OrderRequest struct {
	UserID      string      `json:"user_id"`
	Symbol      string      `json:"symbol"`
	Side        OrderSide   `json:"side"`
	Type        OrderType   `json:"type"`
	Quantity    float64     `json:"quantity"`
	Price       float64     `json:"price,omitempty"` // Required for LIMIT orders
	TimeInForce TimeInForce `json:"time_in_force,omitempty"`
	// Add other fields like StopPrice, ClientOrderID if needed
}

// PlaceOrderResponse represents the response after placing an order
type PlaceOrderResponse struct {
	Order          // Embed the Order struct
	IsSuccess bool `json:"is_success"` // Indicates if the placement was successful initially
	// Add any other relevant fields from the exchange response if needed
}

// OrderResponse is an alias for PlaceOrderResponse for interface compatibility
type OrderResponse = PlaceOrderResponse

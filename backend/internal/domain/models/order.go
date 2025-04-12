package models

import "time"

// Order status and type constants are defined in enums.go

// Order represents a trading order
type Order struct {
	ID        string      `gorm:"primaryKey" json:"id"`
	OrderID   string      `gorm:"uniqueIndex;size:50" json:"order_id"` // Exchange-generated order ID
	ClientID  string      `gorm:"index;size:50" json:"client_id"`      // Client-generated order ID
	Symbol    string      `gorm:"index;not null;size:20" json:"symbol"`
	Side      OrderSide   `gorm:"type:varchar(4);not null" json:"side"`
	Type      OrderType   `gorm:"type:varchar(20);not null" json:"type"`
	Quantity  float64     `gorm:"not null" json:"quantity"`
	FilledQty float64     `gorm:"not null;default:0" json:"filled_qty"` // Amount that has been filled
	Price     float64     `gorm:"not null" json:"price"`
	Status    OrderStatus `gorm:"type:varchar(20);not null;index" json:"status"`
	CreatedAt time.Time   `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	Time      time.Time   `gorm:"-" json:"time"` // Alias for CreatedAt for backward compatibility

	// Foreign key relationship
	PositionID string `gorm:"index" json:"position_id,omitempty"`
}

// CalculateValue returns the total value of the order
func (o *Order) CalculateValue() float64 {
	return o.Price * o.Quantity
}

// IsValid checks if the order has valid properties
func (o *Order) IsValid() bool {
	return o.Symbol != "" && o.Quantity > 0 && (o.Type == OrderTypeMarket || (o.Type == OrderTypeLimit && o.Price > 0))
}

// IsFilled checks if the order is filled
func (o *Order) IsFilled() bool {
	return o.Status == OrderStatusFilled
}

// IsCanceled checks if the order is canceled
func (o *Order) IsCanceled() bool {
	return o.Status == OrderStatusCanceled
}

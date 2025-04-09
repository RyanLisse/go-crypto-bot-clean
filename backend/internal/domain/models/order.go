package models

import "time"

// Order represents a trading order
type Order struct {
	ID        string      `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	OrderID   string      `gorm:"uniqueIndex;size:50" json:"order_id"`  // Exchange-generated order ID
	ClientID  string      `gorm:"index;size:50" json:"client_id"` // Client-generated order ID
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
	PositionID string     `gorm:"index" json:"position_id,omitempty"`
}

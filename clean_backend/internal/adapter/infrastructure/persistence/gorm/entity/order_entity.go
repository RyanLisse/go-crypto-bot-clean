package entity

import (
	"time"
)

// OrderEntity is the GORM model for order data
type OrderEntity struct {
	ID                  string `gorm:"primaryKey"`
	OrderID             string `gorm:"index:idx_order_order_id"`
	UserID              string `gorm:"index:idx_order_user_id"`
	Symbol              string `gorm:"index:idx_order_symbol"`
	Exchange            string
	Side                string
	Type                string
	Status              string `gorm:"index:idx_order_status"`
	TimeInForce         string
	Price               float64
	Quantity            float64
	ExecutedQty         float64
	CummulativeQuoteQty float64
	ClientOrderID       string `gorm:"index:idx_order_client_id"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// TableName sets the table name for OrderEntity
func (OrderEntity) TableName() string {
	return "orders"
}

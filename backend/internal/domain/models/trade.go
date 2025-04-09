package models

import "time"

// TradeSide represents the side of a trade (buy/sell)
type TradeSide string

const (
	TradeSideBuy  TradeSide = "BUY"
	TradeSideSell TradeSide = "SELL"
)

// Trade represents a trade execution on the exchange
type Trade struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	TradeID     string    `gorm:"uniqueIndex;size:50" json:"trade_id"` // Exchange-generated trade ID
	Symbol      string    `gorm:"index;not null;size:20" json:"symbol"`
	Price       float64   `gorm:"not null" json:"price"`
	Quantity    float64   `gorm:"not null" json:"quantity"`
	Side        TradeSide `gorm:"type:varchar(4);not null" json:"side"`
	Timestamp   time.Time `gorm:"index;not null" json:"timestamp"`
	OrderID     string    `gorm:"index;size:50" json:"order_id,omitempty"`
	PositionID  string    `gorm:"index" json:"position_id,omitempty"`
	Commission  float64   `gorm:"default:0" json:"commission,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

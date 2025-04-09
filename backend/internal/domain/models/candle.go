package models

import "time"

// Candle represents a candlestick in a chart
type Candle struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Symbol      string    `gorm:"index;not null;size:20" json:"symbol"`
	Interval    string    `gorm:"index;not null;size:10" json:"interval"`
	OpenTime    time.Time `gorm:"index;not null" json:"open_time"`
	CloseTime   time.Time `gorm:"index;not null" json:"close_time"`
	OpenPrice   float64   `gorm:"not null" json:"open_price"`
	HighPrice   float64   `gorm:"not null" json:"high_price"`
	LowPrice    float64   `gorm:"not null" json:"low_price"`
	ClosePrice  float64   `gorm:"not null" json:"close_price"`
	Volume      float64   `gorm:"not null" json:"volume"`
	QuoteVolume float64   `gorm:"not null" json:"quote_volume"`
	TradeCount  int       `gorm:"not null" json:"trade_count"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// PriceUpdate represents a real-time price update
type PriceUpdate struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Symbol    string    `gorm:"index;not null;size:20" json:"symbol"`
	Price     float64   `gorm:"not null" json:"price"`
	Timestamp time.Time `gorm:"index;not null" json:"timestamp"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// MarketTrade represents a trade that occurred on the exchange
type MarketTrade struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	TradeID   string    `gorm:"uniqueIndex;size:50" json:"trade_id"` // Exchange-generated trade ID
	Symbol    string    `gorm:"index;not null;size:20" json:"symbol"`
	Price     float64   `gorm:"not null" json:"price"`
	Quantity  float64   `gorm:"not null" json:"quantity"`
	Timestamp time.Time `gorm:"index;not null" json:"timestamp"`
	IsBuyer   bool      `gorm:"not null" json:"is_buyer"`
	IsMaker   bool      `gorm:"not null" json:"is_maker"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// CandleOrderBookEntry represents a single entry in the order book
type CandleOrderBookEntry struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	OrderBookID string    `gorm:"index;not null" json:"order_book_id"`
	Type        string    `gorm:"type:varchar(4);not null" json:"type"` // bid or ask
	Price       float64   `gorm:"not null" json:"price"`
	Quantity    float64   `gorm:"not null" json:"quantity"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// OrderBook represents the current state of the order book
type OrderBook struct {
	ID        string                 `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Symbol    string                 `gorm:"index;not null;size:20" json:"symbol"`
	Timestamp time.Time              `gorm:"index;not null" json:"timestamp"`
	Bids      []CandleOrderBookEntry `gorm:"foreignKey:OrderBookID" json:"bids"`
	Asks      []CandleOrderBookEntry `gorm:"foreignKey:OrderBookID" json:"asks"`
	CreatedAt time.Time              `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time              `gorm:"autoUpdateTime" json:"updated_at"`
}

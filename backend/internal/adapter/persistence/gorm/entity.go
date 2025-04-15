package gorm

import (
	"time"
)

// SymbolEntity is the GORM model for trading pair information
type SymbolEntity struct {
	Symbol            string `gorm:"primaryKey"`
	Exchange          string `gorm:"primaryKey;index:idx_symbol_exchange"`
	BaseAsset         string
	QuoteAsset        string
	Status            string
	MinPrice          float64
	MaxPrice          float64
	PricePrecision    int
	MinQty            float64
	MaxQty            float64
	QtyPrecision      int
	AllowedOrderTypes string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// TableName sets the table name for SymbolEntity
func (SymbolEntity) TableName() string {
	return "symbols"
}

// TickerEntity is the GORM model for ticker data
type TickerEntity struct {
	ID            string `gorm:"primaryKey"`
	Symbol        string `gorm:"index:idx_ticker_symbol"`
	Exchange      string `gorm:"index:idx_ticker_exchange"`
	Price         float64
	Volume        float64
	QuoteVolume   float64
	High24h       float64
	Low24h        float64
	PriceChange   float64
	PercentChange float64
	Bid           float64
	Ask           float64
	OpenPrice     float64
	ClosePrice    float64
	LastUpdated   time.Time `gorm:"index:idx_ticker_last_updated"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// TableName sets the table name for TickerEntity
func (TickerEntity) TableName() string {
	return "tickers"
}

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

// PositionEntity is the GORM model for position data
type PositionEntity struct {
	ID              string `gorm:"primaryKey"`
	UserID          string `gorm:"index:idx_position_user_id"`
	Symbol          string `gorm:"index:idx_position_symbol"`
	Exchange        string
	Side            string `gorm:"index"`
	Status          string `gorm:"index"`
	Type            string `gorm:"index"`
	EntryPrice      float64
	Quantity        float64
	CurrentPrice    float64
	PnL             float64
	PnLPercent      float64
	StopLoss        *float64
	TakeProfit      *float64
	StrategyID      *string
	EntryOrderIDs   string // Stored as JSON array
	ExitOrderIDs    string // Stored as JSON array
	OpenOrderIDs    string // Stored as JSON array
	Notes           string
	OpenedAt        time.Time
	ClosedAt        *time.Time
	LastUpdatedAt   time.Time
	MaxDrawdown     float64
	MaxProfit       float64
	RiskRewardRatio float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// TableName sets the table name for PositionEntity
func (PositionEntity) TableName() string {
	return "positions"
}

// WalletEntity is the GORM model for wallet data
type WalletEntity struct {
	ID            uint   `gorm:"primaryKey"`
	UserID        string `gorm:"size:50;not null;index"`
	Exchange      string
	Balances      []byte    `gorm:"type:json"`
	TotalUSDValue float64   `gorm:"type:decimal(18,8);not null"`
	LastUpdated   time.Time `gorm:"not null"`
	LastSyncAt    time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// TableName sets the table name for WalletEntity
func (WalletEntity) TableName() string {
	return "wallets"
}

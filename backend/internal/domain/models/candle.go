package models

import "time"

// Candle represents a price candle (OHLCV)
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

// MarketTrade represents a market trade (different from Trade in trade.go)
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

// CandleOrderBookEntry represents a single entry in the order book (different from OrderBookEntry in orderbook.go)
type CandleOrderBookEntry struct {
	ID           string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	OrderBookID  string    `gorm:"index;not null" json:"order_book_id"`
	Type         string    `gorm:"type:varchar(4);not null" json:"type"` // bid or ask
	Price        float64   `gorm:"not null" json:"price"`
	Quantity     float64   `gorm:"not null" json:"quantity"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// OrderBook represents the market depth
type OrderBook struct {
	ID        string                 `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Symbol    string                 `gorm:"index;not null;size:20" json:"symbol"`
	Timestamp time.Time              `gorm:"index;not null" json:"timestamp"`
	Bids      []CandleOrderBookEntry `gorm:"foreignKey:OrderBookID" json:"bids"`
	Asks      []CandleOrderBookEntry `gorm:"foreignKey:OrderBookID" json:"asks"`
	CreatedAt time.Time              `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time              `gorm:"autoUpdateTime" json:"updated_at"`
}

// BacktestResult represents the result of a strategy backtest
type BacktestResult struct {
	ID                 string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	StrategyID         string    `gorm:"index;not null" json:"strategy_id"`
	StartTime          time.Time `gorm:"not null" json:"start_time"`
	EndTime            time.Time `gorm:"not null" json:"end_time"`
	InitialBalance     float64   `gorm:"not null" json:"initial_balance"`
	FinalBalance       float64   `gorm:"not null" json:"final_balance"`
	ProfitLoss         float64   `gorm:"not null" json:"profit_loss"`
	ProfitLossPercent  float64   `gorm:"not null" json:"profit_loss_percent"`
	TotalTrades        int       `gorm:"not null" json:"total_trades"`
	WinningTrades      int       `gorm:"not null" json:"winning_trades"`
	LosingTrades       int       `gorm:"not null" json:"losing_trades"`
	WinRate            float64   `gorm:"not null" json:"win_rate"`
	AverageWin         float64   `gorm:"not null" json:"average_win"`
	AverageLoss        float64   `gorm:"not null" json:"average_loss"`
	MaxDrawdown        float64   `gorm:"not null" json:"max_drawdown"`
	MaxDrawdownPercent float64   `gorm:"not null" json:"max_drawdown_percent"`
	SharpeRatio        float64   `gorm:"not null" json:"sharpe_ratio"`
	SortinoRatio       float64   `gorm:"not null" json:"sortino_ratio"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

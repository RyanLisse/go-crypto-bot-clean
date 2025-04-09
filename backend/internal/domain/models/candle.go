package models

import "time"

// Candle represents a price candle (OHLCV)
type Candle struct {
	Symbol      string    `json:"symbol"`
	Interval    string    `json:"interval"`
	OpenTime    time.Time `json:"open_time"`
	CloseTime   time.Time `json:"close_time"`
	OpenPrice   float64   `json:"open_price"`
	HighPrice   float64   `json:"high_price"`
	LowPrice    float64   `json:"low_price"`
	ClosePrice  float64   `json:"close_price"`
	Volume      float64   `json:"volume"`
	QuoteVolume float64   `json:"quote_volume"`
	TradeCount  int       `json:"trade_count"`
}

// PriceUpdate represents a real-time price update
type PriceUpdate struct {
	Symbol    string    `json:"symbol"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

// MarketTrade represents a market trade (different from Trade in trade.go)
type MarketTrade struct {
	Symbol    string    `json:"symbol"`
	ID        string    `json:"id"`
	Price     float64   `json:"price"`
	Quantity  float64   `json:"quantity"`
	Timestamp time.Time `json:"timestamp"`
	IsBuyer   bool      `json:"is_buyer"`
	IsMaker   bool      `json:"is_maker"`
}

// CandleOrderBookEntry represents a single entry in the order book (different from OrderBookEntry in orderbook.go)
type CandleOrderBookEntry struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}

// OrderBook represents the market depth
type OrderBook struct {
	Symbol    string                 `json:"symbol"`
	Timestamp time.Time              `json:"timestamp"`
	Bids      []CandleOrderBookEntry `json:"bids"`
	Asks      []CandleOrderBookEntry `json:"asks"`
}

// BacktestResult represents the result of a strategy backtest
type BacktestResult struct {
	StartTime          time.Time `json:"start_time"`
	EndTime            time.Time `json:"end_time"`
	InitialBalance     float64   `json:"initial_balance"`
	FinalBalance       float64   `json:"final_balance"`
	ProfitLoss         float64   `json:"profit_loss"`
	ProfitLossPercent  float64   `json:"profit_loss_percent"`
	TotalTrades        int       `json:"total_trades"`
	WinningTrades      int       `json:"winning_trades"`
	LosingTrades       int       `json:"losing_trades"`
	WinRate            float64   `json:"win_rate"`
	AverageWin         float64   `json:"average_win"`
	AverageLoss        float64   `json:"average_loss"`
	MaxDrawdown        float64   `json:"max_drawdown"`
	MaxDrawdownPercent float64   `json:"max_drawdown_percent"`
	SharpeRatio        float64   `json:"sharpe_ratio"`
	SortinoRatio       float64   `json:"sortino_ratio"`
}

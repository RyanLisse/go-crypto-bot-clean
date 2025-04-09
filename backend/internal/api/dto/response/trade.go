package response

import "time"

// TradeResponse represents a single trading position
type TradeResponse struct {
	ID              uint      `json:"id"`
	Symbol          string    `json:"symbol"`
	PurchasePrice   float64   `json:"purchase_price"`
	CurrentPrice    float64   `json:"current_price"`
	Quantity        float64   `json:"quantity"`
	PurchaseTime    time.Time `json:"purchase_time"`
	ProfitPercent   float64   `json:"profit_percent"`
	CurrentValue    float64   `json:"current_value"`
	StopLossPrice   float64   `json:"stop_loss_price"`
	TakeProfitLevels []TakeProfitLevelResponse `json:"take_profit_levels"`
}

// TakeProfitLevelResponse represents a take profit level
type TakeProfitLevelResponse struct {
	Price    float64 `json:"price"`
	Percent  float64 `json:"percent"`
	Executed bool    `json:"executed"`
}

// TradeHistoryResponse represents a list of historical trades
type TradeHistoryResponse struct {
	Trades    []TradeHistoryItem `json:"trades"`
	Count     int                `json:"count"`
	Timestamp time.Time          `json:"timestamp"`
}

// TradeHistoryItem represents a historical trade
type TradeHistoryItem struct {
	ID           uint      `json:"id"`
	Symbol       string    `json:"symbol"`
	BuyPrice     float64   `json:"buy_price"`
	SellPrice    float64   `json:"sell_price"`
	Quantity     float64   `json:"quantity"`
	BuyTime      time.Time `json:"buy_time"`
	SellTime     time.Time `json:"sell_time"`
	ProfitLoss   float64   `json:"profit_loss"`
	ProfitPercent float64  `json:"profit_percent"`
}

// TradeStatusResponse represents the status of a trade
type TradeStatusResponse struct {
	ID              uint      `json:"id"`
	Symbol          string    `json:"symbol"`
	Status          string    `json:"status"`
	PurchasePrice   float64   `json:"purchase_price"`
	CurrentPrice    float64   `json:"current_price"`
	Quantity        float64   `json:"quantity"`
	PurchaseTime    time.Time `json:"purchase_time"`
	ProfitPercent   float64   `json:"profit_percent"`
	CurrentValue    float64   `json:"current_value"`
	StopLossPrice   float64   `json:"stop_loss_price"`
	TakeProfitLevels []TakeProfitLevelResponse `json:"take_profit_levels"`
}

// TradeExecutionResponse represents the result of a trade execution
type TradeExecutionResponse struct {
	ID            uint      `json:"id"`
	Symbol        string    `json:"symbol"`
	Price         float64   `json:"price"`
	Quantity      float64   `json:"quantity"`
	Total         float64   `json:"total"`
	ExecutionTime time.Time `json:"execution_time"`
	Status        string    `json:"status"`
}

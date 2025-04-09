package request

// TradeRequest represents a request to execute a trade
type TradeRequest struct {
	Symbol       string  `json:"symbol" binding:"required"`
	Side         string  `json:"side" binding:"required,oneof=buy sell"`
	Amount       float64 `json:"amount,omitempty"`
	Price        float64 `json:"price,omitempty"`
	OrderType    string  `json:"order_type,omitempty" binding:"omitempty,oneof=market limit"`
	StopLoss     float64 `json:"stop_loss,omitempty"`
	TakeProfit   float64 `json:"take_profit,omitempty"`
	TrailingStop bool    `json:"trailing_stop,omitempty"`
}

// SellRequest represents a request to sell a coin
type SellRequest struct {
	CoinID      uint    `json:"coin_id" binding:"required"`
	Amount      float64 `json:"amount,omitempty"`
	All         bool    `json:"all,omitempty"`
	Price       float64 `json:"price,omitempty"`
	OrderType   string  `json:"order_type,omitempty" binding:"omitempty,oneof=market limit"`
	MarketOrder bool    `json:"market_order,omitempty"`
}

package market

// Symbol represents a trading pair's information
type Symbol struct {
	// Exchange is the name of the exchange (e.g., "binance")
	Exchange string `json:"exchange"`

	// Symbol is the trading pair (e.g., "BTCUSDT")
	Symbol string `json:"symbol"`

	// BaseAsset is the base asset (e.g., "BTC")
	BaseAsset string `json:"baseAsset"`

	// QuoteAsset is the quote asset (e.g., "USDT")
	QuoteAsset string `json:"quoteAsset"`

	// Status indicates if trading is enabled
	Status string `json:"status"`

	// MinPrice is the minimum price allowed
	MinPrice float64 `json:"minPrice"`

	// MaxPrice is the maximum price allowed
	MaxPrice float64 `json:"maxPrice"`

	// TickSize is the minimum price movement
	TickSize float64 `json:"tickSize"`

	// MinQty is the minimum order quantity
	MinQty float64 `json:"minQty"`

	// MaxQty is the maximum order quantity
	MaxQty float64 `json:"maxQty"`

	// StepSize is the minimum quantity movement
	StepSize float64 `json:"stepSize"`

	// MinNotional is the minimum order value
	MinNotional float64 `json:"minNotional"`
}

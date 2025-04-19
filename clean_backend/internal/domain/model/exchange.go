package model

// ExchangeInfo represents general information about the exchange
type ExchangeInfo struct {
	Symbols []SymbolInfo `json:"symbols"`
	// Add other relevant exchange-wide info if needed
}

// SymbolInfo represents detailed information about a specific trading symbol
// This might duplicate some fields from NewCoin or Symbol, needs consolidation later.
type SymbolInfo struct {
	Symbol               string   `json:"symbol"`
	Status               string   `json:"status"` // e.g., TRADING, BREAK, HALT
	BaseAsset            string   `json:"baseAsset"`
	BaseAssetPrecision   int      `json:"baseAssetPrecision"`
	QuoteAsset           string   `json:"quoteAsset"`
	QuoteAssetPrecision  int      `json:"quotePrecision"` // Note: API might use quotePrecision
	OrderTypes           []string `json:"orderTypes"`     // e.g., ["LIMIT", "MARKET", "STOP_LOSS_LIMIT"]
	IsSpotTradingAllowed bool     `json:"isSpotTradingAllowed"`
	Permissions          []string `json:"permissions"` // e.g., ["SPOT", "MARGIN"]

	// Filters define trading rules (can be complex, simplified here)
	MinNotional string `json:"minNotional,omitempty"` // Minimum order value (price * quantity)
	MinLotSize  string `json:"minLotSize,omitempty"`  // Minimum order quantity
	MaxLotSize  string `json:"maxLotSize,omitempty"`  // Maximum order quantity
	StepSize    string `json:"stepSize,omitempty"`    // Allowed quantity increments
	TickSize    string `json:"tickSize,omitempty"`    // Allowed price increments
	// Add other filters as needed (e.g., PRICE_FILTER, LOT_SIZE, MARKET_LOT_SIZE)

	// Additional precision fields needed for sync_symbols.go
	PricePrecision    int `json:"pricePrecision,omitempty"`    // Number of decimal places in price
	QuantityPrecision int `json:"quantityPrecision,omitempty"` // Number of decimal places in quantity
}

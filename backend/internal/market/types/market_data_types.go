package types

// MarketDataType represents different types of market data
type MarketDataType int

const (
	// MarketDataTypeCandle represents candlestick data
	MarketDataTypeCandle MarketDataType = iota
	// MarketDataTypeTrade represents trade data
	MarketDataTypeTrade
	// MarketDataTypeOrderBook represents order book data
	MarketDataTypeOrderBook
	// MarketDataTypeTicker represents ticker data
	MarketDataTypeTicker
)

// String returns the string representation of MarketDataType
func (t MarketDataType) String() string {
	switch t {
	case MarketDataTypeCandle:
		return "candle"
	case MarketDataTypeTrade:
		return "trade"
	case MarketDataTypeOrderBook:
		return "orderbook"
	case MarketDataTypeTicker:
		return "ticker"
	default:
		return "unknown"
	}
}

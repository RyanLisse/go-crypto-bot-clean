package types

import "time"

// MarketDataFilter represents filtering criteria for market data queries
type MarketDataFilter struct {
	Symbol    string    // Trading pair symbol (e.g., "BTCUSDT")
	StartTime time.Time // Start time for the data range
	EndTime   time.Time // End time for the data range
	Limit     int       // Maximum number of records to return (0 means no limit)
}

package backtest

import (
	"context"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// ConcreteDataProvider implements the DataProvider interface
// This is just an example, replace with your actual implementation
type ConcreteDataProvider struct {
	// Add necessary fields, e.g., API client
}

// GetHistoricalData fetches historical data
func (dp *ConcreteDataProvider) GetHistoricalData(ctx context.Context, symbol string, interval string, startTime time.Time, endTime time.Time) ([]*models.Kline, error) {
	// Implementation to fetch data (e.g., from an exchange API)
	fmt.Printf("Fetching historical data for %s (%s) from %s to %s\n", symbol, interval, startTime, endTime)
	// Return mock data or implement actual fetching logic
	return []*models.Kline{ /* ... mock klines ... */ }, nil
}

// Ensure ConcreteDataProvider implements DataProvider
var _ DataProvider = (*ConcreteDataProvider)(nil)

// SQLiteDataProvider implements the DataProvider interface using a SQLite database
type SQLiteDataProvider struct {
	dbPath string
}

// NewSQLiteDataProvider creates a new SQLiteDataProvider
func NewSQLiteDataProvider(dbPath string) *SQLiteDataProvider {
	return &SQLiteDataProvider{
		dbPath: dbPath,
	}
}

// GetKlines retrieves historical candlestick data from the SQLite database
func (p *SQLiteDataProvider) GetKlines(ctx context.Context, symbol string, interval string, startTime, endTime time.Time) ([]*models.Kline, error) {
	// TODO: Implement SQLite query to retrieve klines
	return nil, nil
}

// GetTickers retrieves historical ticker data from the SQLite database
func (p *SQLiteDataProvider) GetTickers(ctx context.Context, symbol string, startTime, endTime time.Time) ([]*models.Ticker, error) {
	// TODO: Implement SQLite query to retrieve tickers
	return nil, nil
}

// GetOrderBook retrieves historical order book snapshots from the SQLite database
func (p *SQLiteDataProvider) GetOrderBook(ctx context.Context, symbol string, timestamp time.Time) (*models.OrderBookUpdate, error) {
	// TODO: Implement SQLite query to retrieve order book
	return nil, nil
}

// CSVDataProvider is now implemented in csv_data_provider.go

// InMemoryDataProvider implements the DataProvider interface using in-memory data
// This is primarily used for testing
type InMemoryDataProvider struct {
	klines     map[string][]*models.Kline
	tickers    map[string][]*models.Ticker
	orderBooks map[string][]*models.OrderBookUpdate
}

// NewInMemoryDataProvider creates a new InMemoryDataProvider
func NewInMemoryDataProvider() *InMemoryDataProvider {
	return &InMemoryDataProvider{
		klines:     make(map[string][]*models.Kline),
		tickers:    make(map[string][]*models.Ticker),
		orderBooks: make(map[string][]*models.OrderBookUpdate),
	}
}

// AddKlines adds klines to the in-memory data store
func (p *InMemoryDataProvider) AddKlines(symbol string, interval string, klines []*models.Kline) {
	key := symbol + "_" + interval
	p.klines[key] = klines
}

// AddTickers adds tickers to the in-memory data store
func (p *InMemoryDataProvider) AddTickers(symbol string, tickers []*models.Ticker) {
	p.tickers[symbol] = tickers
}

// AddOrderBooks adds order books to the in-memory data store
func (p *InMemoryDataProvider) AddOrderBooks(symbol string, orderBooks []*models.OrderBookUpdate) {
	p.orderBooks[symbol] = orderBooks
}

// GetHistoricalData retrieves historical candlestick data from the in-memory data store
// (Renamed from GetKlines to match DataProvider interface)
func (p *InMemoryDataProvider) GetHistoricalData(ctx context.Context, symbol string, interval string, startTime, endTime time.Time) ([]*models.Kline, error) {
	key := symbol + "_" + interval
	klines, ok := p.klines[key]
	if !ok {
		return nil, nil // Or return an error if data for the key must exist
	}

	var result []*models.Kline
	for _, kline := range klines {
		// Ensure times are valid before comparison
		if kline.OpenTime.IsZero() || startTime.IsZero() || endTime.IsZero() {
			continue // Skip klines with zero times or invalid range
		}
		// Filter klines within the requested time range [startTime, endTime]
		if !kline.OpenTime.Before(startTime) && !kline.OpenTime.After(endTime) {
			result = append(result, kline)
		}
	}

	return result, nil
}

// // GetTickers retrieves historical ticker data from the in-memory data store
// // Commenting out as it's not part of the current DataProvider interface
// func (p *InMemoryDataProvider) GetTickers(ctx context.Context, symbol string, startTime, endTime time.Time) ([]*models.Ticker, error) {
// 	tickers, ok := p.tickers[symbol]
// 	if !ok {
// 		return nil, nil
// 	}

// 	var result []*models.Ticker
// 	for _, ticker := range tickers {
// 		tickerTime := ticker.Timestamp
// 		if tickerTime.After(startTime) && tickerTime.Before(endTime) {
// 			result = append(result, ticker)
// 		}
// 	}

// 	return result, nil
// }

// // GetOrderBook retrieves historical order book snapshots from the in-memory data store
// // Commenting out as it's not part of the current DataProvider interface
// func (p *InMemoryDataProvider) GetOrderBook(ctx context.Context, symbol string, timestamp time.Time) (*models.OrderBookUpdate, error) {
// 	orderBooks, ok := p.orderBooks[symbol]
// 	if !ok {
// 		return nil, nil
// 	}

// 	var closestOrderBook *models.OrderBookUpdate
// 	var minDiff int64 = math.MaxInt64

// 	for _, orderBook := range orderBooks {
// 		diff := abs(orderBook.Timestamp.UnixMilli() - timestamp.UnixMilli())
// 		if diff < minDiff {
// 			minDiff = diff
// 			closestOrderBook = orderBook
// 		}
// 	}

// 	return closestOrderBook, nil
// }

// abs function is now in csv_data_provider.go

package backtest

import (
	"context"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// DataProvider defines the interface for retrieving historical market data
type DataProvider interface {
	// GetKlines retrieves historical candlestick data for a symbol within a time range
	GetKlines(ctx context.Context, symbol string, interval string, startTime, endTime time.Time) ([]*models.Kline, error)

	// GetTickers retrieves historical ticker data for a symbol within a time range
	GetTickers(ctx context.Context, symbol string, startTime, endTime time.Time) ([]*models.Ticker, error)

	// GetOrderBook retrieves historical order book snapshots for a symbol at a specific time
	GetOrderBook(ctx context.Context, symbol string, timestamp time.Time) (*models.OrderBookUpdate, error)
}

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

// CSVDataProvider implements the DataProvider interface using CSV files
type CSVDataProvider struct {
	dataDir string
}

// NewCSVDataProvider creates a new CSVDataProvider
func NewCSVDataProvider(dataDir string) *CSVDataProvider {
	return &CSVDataProvider{
		dataDir: dataDir,
	}
}

// GetKlines retrieves historical candlestick data from CSV files
func (p *CSVDataProvider) GetKlines(ctx context.Context, symbol string, interval string, startTime, endTime time.Time) ([]*models.Kline, error) {
	// TODO: Implement CSV file reading to retrieve klines
	return nil, nil
}

// GetTickers retrieves historical ticker data from CSV files
func (p *CSVDataProvider) GetTickers(ctx context.Context, symbol string, startTime, endTime time.Time) ([]*models.Ticker, error) {
	// TODO: Implement CSV file reading to retrieve tickers
	return nil, nil
}

// GetOrderBook retrieves historical order book snapshots from CSV files
func (p *CSVDataProvider) GetOrderBook(ctx context.Context, symbol string, timestamp time.Time) (*models.OrderBookUpdate, error) {
	// TODO: Implement CSV file reading to retrieve order book
	return nil, nil
}

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

// GetKlines retrieves historical candlestick data from the in-memory data store
func (p *InMemoryDataProvider) GetKlines(ctx context.Context, symbol string, interval string, startTime, endTime time.Time) ([]*models.Kline, error) {
	key := symbol + "_" + interval
	klines, ok := p.klines[key]
	if !ok {
		return nil, nil
	}

	var result []*models.Kline
	for _, kline := range klines {
		// Convert time.Time to int64 for comparison
		startMillis := startTime.UnixMilli()
		endMillis := endTime.UnixMilli()

		// Get Unix milliseconds from kline times
		openTimeMillis := kline.OpenTime.UnixMilli()
		closeTimeMillis := kline.CloseTime.UnixMilli()

		if (openTimeMillis >= startMillis && openTimeMillis <= endMillis) ||
			(closeTimeMillis >= startMillis && closeTimeMillis <= endMillis) {
			result = append(result, kline)
		}
	}

	return result, nil
}

// GetTickers retrieves historical ticker data from the in-memory data store
func (p *InMemoryDataProvider) GetTickers(ctx context.Context, symbol string, startTime, endTime time.Time) ([]*models.Ticker, error) {
	tickers, ok := p.tickers[symbol]
	if !ok {
		return nil, nil
	}

	var result []*models.Ticker
	for _, ticker := range tickers {
		// Use the ticker's timestamp directly since it's already a time.Time
		tickerTime := ticker.Timestamp

		if tickerTime.After(startTime) && tickerTime.Before(endTime) {
			result = append(result, ticker)
		}
	}

	return result, nil
}

// GetOrderBook retrieves historical order book snapshots from the in-memory data store
func (p *InMemoryDataProvider) GetOrderBook(ctx context.Context, symbol string, timestamp time.Time) (*models.OrderBookUpdate, error) {
	orderBooks, ok := p.orderBooks[symbol]
	if !ok {
		return nil, nil
	}

	// Find the closest order book to the requested timestamp
	var closestOrderBook *models.OrderBookUpdate
	var minDiff int64 = 9223372036854775807 // Max int64

	for _, orderBook := range orderBooks {
		diff := abs(orderBook.Timestamp.UnixMilli() - timestamp.UnixMilli())
		if diff < minDiff {
			minDiff = diff
			closestOrderBook = orderBook
		}
	}

	return closestOrderBook, nil
}

// abs returns the absolute value of an int64
func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

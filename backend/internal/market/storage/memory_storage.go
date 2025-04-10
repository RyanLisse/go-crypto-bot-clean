package storage

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// MemoryStorage implements MarketDataStorage interface for in-memory storage
type MemoryStorage struct {
	mu      sync.RWMutex
	maxSize int
	candles map[string][]*models.Candle
	trades  map[string][]*models.MarketTrade
	books   map[string]*models.OrderBook
	tickers map[string]*models.Ticker
}

// NewMemoryStorage creates a new MemoryStorage instance
func NewMemoryStorage(maxSize int) *MemoryStorage {
	return &MemoryStorage{
		maxSize: maxSize,
		candles: make(map[string][]*models.Candle),
		trades:  make(map[string][]*models.MarketTrade),
		books:   make(map[string]*models.OrderBook),
		tickers: make(map[string]*models.Ticker),
	}
}

// StoreCandles implements MarketDataStorage.StoreCandles
func (m *MemoryStorage) StoreCandles(ctx context.Context, symbol string, candles []*models.Candle) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.candles[symbol] == nil {
		m.candles[symbol] = make([]*models.Candle, 0)
	}

	// Sort new candles by OpenTime
	sort.Slice(candles, func(i, j int) bool {
		return candles[i].OpenTime.Before(candles[j].OpenTime)
	})

	m.candles[symbol] = append(m.candles[symbol], candles...)
	m.enforceMaxSize(symbol)
	return nil
}

// StoreTrades implements MarketDataStorage.StoreTrades
func (m *MemoryStorage) StoreTrades(ctx context.Context, symbol string, trades []*models.MarketTrade) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.trades[symbol] == nil {
		m.trades[symbol] = make([]*models.MarketTrade, 0)
	}

	m.trades[symbol] = append(m.trades[symbol], trades...)
	m.enforceMaxSize(symbol)
	return nil
}

// StoreOrderBook implements MarketDataStorage.StoreOrderBook
func (m *MemoryStorage) StoreOrderBook(ctx context.Context, symbol string, orderBook *models.OrderBook) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.books[symbol] = orderBook
	return nil
}

// StoreTicker implements MarketDataStorage.StoreTicker
func (m *MemoryStorage) StoreTicker(ctx context.Context, symbol string, ticker *models.Ticker) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tickers[symbol] = ticker
	return nil
}

// GetCandles implements MarketDataStorage.GetCandles
func (m *MemoryStorage) GetCandles(ctx context.Context, symbol string, start, end time.Time) ([]*models.Candle, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	candles := m.candles[symbol]
	if len(candles) == 0 {
		return nil, nil
	}

	// Filter candles by time range
	var filtered []*models.Candle
	for _, c := range candles {
		if (c.OpenTime.Equal(start) || c.OpenTime.After(start)) &&
			(c.OpenTime.Equal(end) || c.OpenTime.Before(end)) {
			filtered = append(filtered, c)
		}
	}

	return filtered, nil
}

// GetTrades implements MarketDataStorage.GetTrades
func (m *MemoryStorage) GetTrades(ctx context.Context, filter MarketDataFilter) ([]*models.MarketTrade, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	trades, ok := m.trades[filter.Symbol]
	if !ok {
		return nil, nil
	}

	filtered := make([]*models.MarketTrade, 0)
	for _, trade := range trades {
		if (filter.StartTime.IsZero() || !trade.Timestamp.Before(filter.StartTime)) &&
			(filter.EndTime.IsZero() || !trade.Timestamp.After(filter.EndTime)) {
			filtered = append(filtered, trade)
		}
	}

	// Apply limit if specified
	if filter.Limit > 0 && len(filtered) > filter.Limit {
		filtered = filtered[len(filtered)-filter.Limit:]
	}

	return filtered, nil
}

// GetOrderBook implements MarketDataStorage.GetOrderBook
func (m *MemoryStorage) GetOrderBook(ctx context.Context, symbol string, timestamp time.Time) (*models.OrderBook, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	book := m.books[symbol]
	if book == nil {
		return nil, fmt.Errorf("no order book found for symbol %s", symbol)
	}

	return book, nil
}

// GetTicker implements MarketDataStorage.GetTicker
func (m *MemoryStorage) GetTicker(ctx context.Context, symbol string) (*models.Ticker, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ticker := m.tickers[symbol]
	if ticker == nil {
		return nil, fmt.Errorf("no ticker found for symbol %s", symbol)
	}

	return ticker, nil
}

// GetVWAP implements MarketDataStorage.GetVWAP
func (m *MemoryStorage) GetVWAP(ctx context.Context, symbol string, startTime, endTime time.Time) (float64, error) {
	trades, err := m.GetTrades(ctx, MarketDataFilter{
		Symbol:    symbol,
		StartTime: startTime,
		EndTime:   endTime,
	})
	if err != nil {
		return 0, err
	}

	return m.calculateVWAP(trades), nil
}

// GetOHLCV implements MarketDataStorage.GetOHLCV
func (m *MemoryStorage) GetOHLCV(ctx context.Context, symbol string, interval string, startTime, endTime time.Time) ([]*models.Candle, error) {
	candles, err := m.GetCandles(ctx, symbol, startTime, endTime)
	if err != nil {
		return nil, err
	}

	if len(candles) == 0 {
		return nil, nil
	}

	open, high, low, close, volume := m.calculateOHLCV(candles)

	// Return a single consolidated candle for the period
	return []*models.Candle{
		{
			Symbol:     symbol,
			OpenTime:   startTime,
			CloseTime:  endTime,
			OpenPrice:  open,
			HighPrice:  high,
			LowPrice:   low,
			ClosePrice: close,
			Volume:     volume,
		},
	}, nil
}

// GetVolume implements MarketDataStorage.GetVolume
func (m *MemoryStorage) GetVolume(ctx context.Context, symbol string, startTime, endTime time.Time) (float64, error) {
	trades, err := m.GetTrades(ctx, MarketDataFilter{
		Symbol:    symbol,
		StartTime: startTime,
		EndTime:   endTime,
	})
	if err != nil {
		return 0, err
	}

	return m.calculateVolume(trades), nil
}

// Cleanup implements MarketDataStorage.Cleanup
func (m *MemoryStorage) Cleanup(ctx context.Context, olderThan time.Time, dataType MarketDataType) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch dataType {
	case MarketDataTypeCandle:
		for symbol, candles := range m.candles {
			filtered := make([]*models.Candle, 0)
			for _, c := range candles {
				if c.OpenTime.After(olderThan) {
					filtered = append(filtered, c)
				}
			}
			m.candles[symbol] = filtered
		}
	case MarketDataTypeTrade:
		for symbol, trades := range m.trades {
			filtered := make([]*models.MarketTrade, 0)
			for _, t := range trades {
				if t.Timestamp.After(olderThan) {
					filtered = append(filtered, t)
				}
			}
			m.trades[symbol] = filtered
		}
	}

	return nil
}

// Migrate implements MarketDataStorage.Migrate
func (m *MemoryStorage) Migrate(ctx context.Context, fromLevel, toLevel StorageLevel, olderThan time.Time) error {
	// No-op for memory storage as it only implements hot storage
	return nil
}

// Vacuum implements MarketDataStorage.Vacuum
func (m *MemoryStorage) Vacuum(ctx context.Context) error {
	// No-op for memory storage
	return nil
}

// GetDataTypes implements MarketDataStorage.GetDataTypes
func (m *MemoryStorage) GetDataTypes(ctx context.Context, symbol string) ([]MarketDataType, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	types := make([]MarketDataType, 0)
	if _, ok := m.candles[symbol]; ok {
		types = append(types, MarketDataTypeCandle)
	}
	if _, ok := m.trades[symbol]; ok {
		types = append(types, MarketDataTypeTrade)
	}
	if _, ok := m.books[symbol]; ok {
		types = append(types, MarketDataTypeOrderBook)
	}
	if _, ok := m.tickers[symbol]; ok {
		types = append(types, MarketDataTypeTicker)
	}

	return types, nil
}

// GetTimeRange implements MarketDataStorage.GetTimeRange
func (m *MemoryStorage) GetTimeRange(ctx context.Context, symbol string, dataType MarketDataType) (time.Time, time.Time, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	switch dataType {
	case MarketDataTypeCandle:
		candles, ok := m.candles[symbol]
		if !ok || len(candles) == 0 {
			return time.Time{}, time.Time{}, nil
		}
		// Sort candles by OpenTime to ensure correct range
		sort.Slice(candles, func(i, j int) bool {
			return candles[i].OpenTime.Before(candles[j].OpenTime)
		})
		return candles[0].OpenTime, candles[len(candles)-1].OpenTime, nil
	case MarketDataTypeTrade:
		trades := m.trades[symbol]
		if len(trades) == 0 {
			return time.Time{}, time.Time{}, nil
		}
		firstTrade := trades[0]
		lastTrade := trades[len(trades)-1]
		return firstTrade.Timestamp, lastTrade.Timestamp, nil
	default:
		return time.Time{}, time.Time{}, fmt.Errorf("unsupported data type: %v", dataType)
	}
}

// GetSymbols implements MarketDataStorage.GetSymbols
func (m *MemoryStorage) GetSymbols(ctx context.Context) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	symbols := make(map[string]struct{})
	for symbol := range m.candles {
		symbols[symbol] = struct{}{}
	}
	for symbol := range m.trades {
		symbols[symbol] = struct{}{}
	}
	for symbol := range m.books {
		symbols[symbol] = struct{}{}
	}
	for symbol := range m.tickers {
		symbols[symbol] = struct{}{}
	}

	result := make([]string, 0, len(symbols))
	for symbol := range symbols {
		result = append(result, symbol)
	}
	return result, nil
}

// enforceMaxSize ensures the number of records for a symbol doesn't exceed maxSize
func (m *MemoryStorage) enforceMaxSize(symbol string) {
	// Enforce size limit for candles
	if candles, exists := m.candles[symbol]; exists && len(candles) > m.maxSize {
		// Sort by OpenTime in ascending order
		sort.Slice(candles, func(i, j int) bool {
			return candles[i].OpenTime.Before(candles[j].OpenTime)
		})
		// Keep only the most recent maxSize candles
		m.candles[symbol] = candles[len(candles)-m.maxSize:]
	}

	// Enforce size limit for trades
	if trades, exists := m.trades[symbol]; exists && len(trades) > m.maxSize {
		// Sort by Timestamp in ascending order
		sort.Slice(trades, func(i, j int) bool {
			return trades[i].Timestamp.Before(trades[j].Timestamp)
		})
		// Keep only the most recent maxSize trades
		m.trades[symbol] = trades[len(trades)-m.maxSize:]
	}
}

// calculateVWAP calculates the Volume Weighted Average Price
func (m *MemoryStorage) calculateVWAP(trades []*models.MarketTrade) float64 {
	if len(trades) == 0 {
		return 0
	}

	var totalVolume, weightedSum float64
	for _, trade := range trades {
		volume := trade.Quantity
		totalVolume += volume
		weightedSum += trade.Price * volume
	}

	if totalVolume == 0 {
		return 0
	}

	return weightedSum / totalVolume
}

// calculateOHLCV calculates OHLCV data from candles
func (m *MemoryStorage) calculateOHLCV(candles []*models.Candle) (float64, float64, float64, float64, float64) {
	if len(candles) == 0 {
		return 0, 0, 0, 0, 0
	}

	// Sort candles by OpenTime to ensure correct OHLCV calculation
	sort.Slice(candles, func(i, j int) bool {
		return candles[i].OpenTime.Before(candles[j].OpenTime)
	})

	open := candles[0].OpenPrice
	high := candles[0].HighPrice
	low := candles[0].LowPrice
	close := candles[len(candles)-1].ClosePrice
	volume := 0.0

	for _, c := range candles {
		if c.HighPrice > high {
			high = c.HighPrice
		}
		if c.LowPrice < low {
			low = c.LowPrice
		}
		volume += c.Volume
	}

	return open, high, low, close, volume
}

// calculateVolume calculates total volume from trades
func (m *MemoryStorage) calculateVolume(trades []*models.MarketTrade) float64 {
	volume := 0.0
	for _, trade := range trades {
		volume += trade.Quantity
	}
	return volume
}

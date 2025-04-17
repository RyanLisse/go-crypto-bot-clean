package gorm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/compat"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// Ensure MarketRepository implements the proper interfaces
var _ port.MarketRepository = (*MarketRepository)(nil)
var _ port.SymbolRepository = (*MarketRepository)(nil)

// Ticker entity is defined in entity.go

// CandleEntity is the GORM model for candlestick data
type CandleEntity struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Symbol      string    `gorm:"index:idx_candle_symbol"`
	Exchange    string    `gorm:"index:idx_candle_exchange"`
	Interval    string    `gorm:"index:idx_candle_interval"`
	OpenTime    time.Time `gorm:"index:idx_candle_opentime"`
	CloseTime   time.Time
	Open        float64
	High        float64
	Low         float64
	Close       float64
	Volume      float64
	QuoteVolume float64
	TradeCount  int64
	Complete    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName sets the table name for CandleEntity
func (CandleEntity) TableName() string {
	return "candles"
}

// OrderBookEntryEntity is the GORM model for order book entries
type OrderBookEntryEntity struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	OrderBookID uint   `gorm:"index:idx_orderbook_entry"`
	Type        string `gorm:"index:idx_entry_type"` // "bid" or "ask"
	Price       float64
	Quantity    float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName sets the table name for OrderBookEntryEntity
func (OrderBookEntryEntity) TableName() string {
	return "orderbook_entries"
}

// OrderBookEntity is the GORM model for order book data
type OrderBookEntity struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Symbol       string    `gorm:"index:idx_orderbook_symbol"`
	Exchange     string    `gorm:"index:idx_orderbook_exchange"`
	LastUpdated  time.Time `gorm:"index:idx_orderbook_updated"`
	SequenceNum  int64
	LastUpdateID int64
	Entries      []OrderBookEntryEntity `gorm:"foreignKey:OrderBookID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TableName sets the table name for OrderBookEntity
func (OrderBookEntity) TableName() string {
	return "orderbooks"
}

// Symbol entity is defined in entity.go

// MarketRepository implements the port.MarketRepository interface using GORM
type MarketRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewMarketRepository creates a new MarketRepository
func NewMarketRepository(db *gorm.DB, logger *zerolog.Logger) *MarketRepository {
	return &MarketRepository{
		db:     db,
		logger: logger,
	}
}

// SaveTicker stores a ticker in the database
func (r *MarketRepository) SaveTicker(ctx context.Context, ticker *model.Ticker) error {
	entity := r.tickerToEntity(ticker)

	result := r.db.WithContext(ctx).Save(&entity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", ticker.Symbol).Msg("Failed to save ticker")
		return fmt.Errorf("failed to save ticker: %w", result.Error)
	}

	r.logger.Info().Str("symbol", ticker.Symbol).Str("exchange", ticker.Exchange).Msg("Ticker saved successfully")
	return nil
}

// GetTicker retrieves the latest ticker for a symbol from a specific exchange
func (r *MarketRepository) GetTicker(ctx context.Context, symbol, exchange string) (*model.Ticker, error) {
	var entity TickerEntity

	result := r.db.WithContext(ctx).
		Where("symbol = ? AND exchange = ?", symbol, exchange).
		Order("last_updated DESC").
		First(&entity)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			r.logger.Info().Str("symbol", symbol).Str("exchange", exchange).Msg("Ticker not found")
			return nil, apperror.ErrNotFound
		}
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get ticker")
		return nil, fmt.Errorf("failed to get ticker: %w", result.Error)
	}

	return r.tickerToDomain(&entity), nil
}

// GetAllTickers retrieves all latest tickers from a specific exchange
func (r *MarketRepository) GetAllTickers(ctx context.Context, exchange string) ([]*model.Ticker, error) {
	var entities []TickerEntity

	// Using a subquery to get the latest ticker for each symbol
	subQuery := r.db.Model(&TickerEntity{}).
		Select("symbol, MAX(last_updated) as max_updated").
		Where("exchange = ?", exchange).
		Group("symbol")

	result := r.db.WithContext(ctx).
		Joins("JOIN (?) as sub ON tickers.symbol = sub.symbol AND tickers.last_updated = sub.max_updated", subQuery).
		Where("tickers.exchange = ?", exchange).
		Find(&entities)

	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("exchange", exchange).Msg("Failed to get all tickers")
		return nil, fmt.Errorf("failed to get all tickers: %w", result.Error)
	}

	tickers := make([]*model.Ticker, len(entities))
	for i, entity := range entities {
		tickers[i] = r.tickerToDomain(&entity)
	}

	return tickers, nil
}

// GetTickerHistory retrieves ticker history for a symbol within a time range
func (r *MarketRepository) GetTickerHistory(ctx context.Context, symbol, exchange string, start, end time.Time) ([]*model.Ticker, error) {
	var entities []TickerEntity

	result := r.db.WithContext(ctx).
		Where("symbol = ? AND exchange = ? AND last_updated BETWEEN ? AND ?",
			symbol, exchange, start, end).
		Order("last_updated ASC").
		Find(&entities)

	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get ticker history")
		return nil, fmt.Errorf("failed to get ticker history: %w", result.Error)
	}

	tickers := make([]*model.Ticker, len(entities))
	for i, entity := range entities {
		tickers[i] = r.tickerToDomain(&entity)
	}

	return tickers, nil
}

// SaveKline stores a kline/candle in the database
func (r *MarketRepository) SaveKline(ctx context.Context, kline *model.Kline) error {
	return r.SaveCandle(ctx, kline)
}

// SaveCandle stores a candle in the database
func (r *MarketRepository) SaveCandle(ctx context.Context, kline *model.Kline) error {
	entity := r.candleToEntity(kline)

	// Try to find an existing candle with the same symbol, exchange, interval, and open time
	var existing CandleEntity
	result := r.db.WithContext(ctx).
		Where("symbol = ? AND exchange = ? AND interval = ? AND open_time = ?",
			kline.Symbol, kline.Exchange, string(kline.Interval), kline.OpenTime).
		First(&existing)

	// If the candle exists, update it; otherwise, create a new one
	if result.Error == nil {
		entity.ID = existing.ID
	}

	result = r.db.WithContext(ctx).Save(&entity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", kline.Symbol).Msg("Failed to save candle")
		return fmt.Errorf("failed to save candle: %w", result.Error)
	}

	r.logger.Info().Str("symbol", kline.Symbol).Str("interval", string(kline.Interval)).Msg("Candle saved successfully")
	return nil
}

// SaveKlines stores multiple klines/candles in the database
func (r *MarketRepository) SaveKlines(ctx context.Context, klines []*model.Kline) error {
	return r.SaveCandles(ctx, klines)
}

// SaveCandles stores multiple candles in the database
func (r *MarketRepository) SaveCandles(ctx context.Context, klines []*model.Kline) error {
	if len(klines) == 0 {
		return nil
	}

	entities := make([]CandleEntity, len(klines))
	for i, kline := range klines {
		entities[i] = r.candleToEntity(kline)
	}

	// Use a transaction to save all candles
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		r.logger.Error().Err(tx.Error).Msg("Failed to begin transaction")
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Create or update each candle
	for i, entity := range entities {
		var existing CandleEntity
		result := tx.Where("symbol = ? AND exchange = ? AND interval = ? AND open_time = ?",
			entity.Symbol, entity.Exchange, entity.Interval, entity.OpenTime).
			First(&existing)

		if result.Error == nil {
			entity.ID = existing.ID
		}

		result = tx.Save(&entity)
		if result.Error != nil {
			tx.Rollback()
			r.logger.Error().Err(result.Error).Str("symbol", klines[i].Symbol).Msg("Failed to save candle in batch")
			return fmt.Errorf("failed to save candle in batch: %w", result.Error)
		}
	}

	if err := tx.Commit().Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to commit transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Info().Int("count", len(klines)).Msg("Successfully saved batch of candles")
	return nil
}

// GetKline retrieves a specific kline/candle for a symbol, interval, and time
func (r *MarketRepository) GetKline(ctx context.Context, symbol, exchange string, interval model.KlineInterval, openTime time.Time) (*model.Kline, error) {
	return r.GetCandle(ctx, symbol, exchange, interval, openTime)
}

// GetCandle retrieves a specific candle for a symbol, interval, and time
func (r *MarketRepository) GetCandle(ctx context.Context, symbol, exchange string, interval model.KlineInterval, openTime time.Time) (*model.Kline, error) {
	var entity CandleEntity

	result := r.db.WithContext(ctx).
		Where("symbol = ? AND exchange = ? AND interval = ? AND open_time = ?",
			symbol, exchange, string(interval), openTime).
		First(&entity)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			r.logger.Info().Str("symbol", symbol).Str("interval", string(interval)).Str("openTime", openTime.Format(time.RFC3339)).Msg("Candle not found")
			return nil, apperror.ErrNotFound
		}
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get candle")
		return nil, fmt.Errorf("failed to get candle: %w", result.Error)
	}

	return r.candleToDomain(&entity), nil
}

// GetKlines retrieves klines/candles for a symbol within a time range
func (r *MarketRepository) GetKlines(ctx context.Context, symbol, exchange string, interval model.KlineInterval, start, end time.Time, limit int) ([]*model.Kline, error) {
	return r.GetCandles(ctx, symbol, exchange, interval, start, end, limit)
}

// GetCandles retrieves candles for a symbol within a time range
func (r *MarketRepository) GetCandles(ctx context.Context, symbol, exchange string, interval model.KlineInterval, start, end time.Time, limit int) ([]*model.Kline, error) {
	var entities []CandleEntity

	query := r.db.WithContext(ctx).
		Where("symbol = ? AND exchange = ? AND interval = ? AND open_time BETWEEN ? AND ?",
			symbol, exchange, string(interval), start, end).
		Order("open_time ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	result := query.Find(&entities)

	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get candles")
		return nil, fmt.Errorf("failed to get candles: %w", result.Error)
	}

	klines := make([]*model.Kline, len(entities))
	for i, entity := range entities {
		klines[i] = r.candleToDomain(&entity)
	}

	return klines, nil
}

// GetLatestKline retrieves the most recent kline/candle for a symbol and interval
func (r *MarketRepository) GetLatestKline(ctx context.Context, symbol, exchange string, interval model.KlineInterval) (*model.Kline, error) {
	return r.GetLatestCandle(ctx, symbol, exchange, interval)
}

// GetLatestCandle retrieves the most recent candle for a symbol and interval
func (r *MarketRepository) GetLatestCandle(ctx context.Context, symbol, exchange string, interval model.KlineInterval) (*model.Kline, error) {
	var entity CandleEntity

	result := r.db.WithContext(ctx).
		Where("symbol = ? AND exchange = ? AND interval = ?",
			symbol, exchange, string(interval)).
		Order("open_time DESC").
		First(&entity)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			r.logger.Info().Str("symbol", symbol).Str("interval", string(interval)).Msg("Latest candle not found")
			return nil, apperror.ErrNotFound
		}
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get latest candle")
		return nil, fmt.Errorf("failed to get latest candle: %w", result.Error)
	}

	return r.candleToDomain(&entity), nil
}

// PurgeOldData removes market data older than the specified retention period
func (r *MarketRepository) PurgeOldData(ctx context.Context, olderThan time.Time) error {
	// Delete old ticker data
	if err := r.db.WithContext(ctx).Where("last_updated < ?", olderThan).Delete(&TickerEntity{}).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to purge old ticker data")
		return fmt.Errorf("failed to purge old ticker data: %w", err)
	}

	// Delete old candle data
	if err := r.db.WithContext(ctx).Where("open_time < ?", olderThan).Delete(&CandleEntity{}).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to purge old candle data")
		return fmt.Errorf("failed to purge old candle data: %w", err)
	}

	// Delete old orderbook data
	if err := r.db.WithContext(ctx).Where("last_updated < ?", olderThan).Delete(&OrderBookEntity{}).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to purge old orderbook data")
		return fmt.Errorf("failed to purge old orderbook data: %w", err)
	}

	r.logger.Info().Time("older_than", olderThan).Msg("Successfully purged old market data")
	return nil
}

// GetLatestTickers retrieves the latest tickers for all symbols
func (r *MarketRepository) GetLatestTickers(ctx context.Context, limit int) ([]*model.Ticker, error) {
	var entities []TickerEntity

	// Using a common table expression (CTE) to get the latest ticker for each symbol
	query := r.db.WithContext(ctx).
		Raw(`WITH latest_tickers AS (
			SELECT symbol, exchange, MAX(last_updated) as max_updated
			FROM tickers
			GROUP BY symbol, exchange
		)
		SELECT t.*
		FROM tickers t
		JOIN latest_tickers lt ON t.symbol = lt.symbol AND t.exchange = lt.exchange AND t.last_updated = lt.max_updated
		ORDER BY t.symbol`)

	if limit > 0 {
		query = query.Limit(limit)
	}

	result := query.Find(&entities)

	if result.Error != nil {
		r.logger.Error().Err(result.Error).Msg("Failed to get latest tickers")
		return nil, fmt.Errorf("failed to get latest tickers: %w", result.Error)
	}

	tickers := make([]*model.Ticker, len(entities))
	for i, entity := range entities {
		tickers[i] = r.tickerToDomain(&entity)
	}

	return tickers, nil
}

// GetTickersBySymbol retrieves tickers for a specific symbol with optional time range
func (r *MarketRepository) GetTickersBySymbol(ctx context.Context, symbol string, limit int) ([]*model.Ticker, error) {
	var entities []TickerEntity

	query := r.db.WithContext(ctx).
		Where("symbol = ?", symbol).
		Order("last_updated DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	result := query.Find(&entities)

	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get tickers by symbol")
		return nil, fmt.Errorf("failed to get tickers by symbol: %w", result.Error)
	}

	tickers := make([]*model.Ticker, len(entities))
	for i, entity := range entities {
		tickers[i] = r.tickerToDomain(&entity)
	}

	return tickers, nil
}

// Symbol Repository implementation

// Create stores a new Symbol
func (r *MarketRepository) Create(ctx context.Context, symbol *model.Symbol) error {
	entity := r.symbolToEntity(symbol)

	result := r.db.WithContext(ctx).Create(&entity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", symbol.Symbol).Msg("Failed to create symbol")
		return fmt.Errorf("failed to create symbol: %w", result.Error)
	}

	r.logger.Info().Str("symbol", symbol.Symbol).Str("exchange", symbol.Exchange).Msg("Symbol created successfully")
	return nil
}

// GetBySymbol returns a Symbol by its symbol string (e.g., "BTCUSDT")
func (r *MarketRepository) GetBySymbol(ctx context.Context, symbol string) (*model.Symbol, error) {
	var entity SymbolEntity

	result := r.db.WithContext(ctx).Where("symbol = ?", symbol).First(&entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			r.logger.Info().Str("symbol", symbol).Msg("Symbol not found")
			return nil, apperror.ErrNotFound
		}
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get symbol")
		return nil, fmt.Errorf("failed to get symbol: %w", result.Error)
	}

	return r.symbolToDomain(&entity), nil
}

// GetByExchange returns all Symbols from a specific exchange
func (r *MarketRepository) GetByExchange(ctx context.Context, exchange string) ([]*model.Symbol, error) {
	var entities []SymbolEntity

	result := r.db.WithContext(ctx).Where("exchange = ?", exchange).Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("exchange", exchange).Msg("Failed to get symbols by exchange")
		return nil, fmt.Errorf("failed to get symbols by exchange: %w", result.Error)
	}

	symbols := make([]*model.Symbol, len(entities))
	for i, entity := range entities {
		symbols[i] = r.symbolToDomain(&entity)
	}

	r.logger.Info().Str("exchange", exchange).Int("count", len(symbols)).Msg("Retrieved symbols by exchange")
	return symbols, nil
}

// GetAll returns all available Symbols
func (r *MarketRepository) GetAll(ctx context.Context) ([]*model.Symbol, error) {
	var entities []SymbolEntity

	result := r.db.WithContext(ctx).Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Msg("Failed to get all symbols")
		return nil, fmt.Errorf("failed to get all symbols: %w", result.Error)
	}

	symbols := make([]*model.Symbol, len(entities))
	for i, entity := range entities {
		symbols[i] = r.symbolToDomain(&entity)
	}

	r.logger.Info().Int("count", len(symbols)).Msg("Retrieved all symbols")
	return symbols, nil
}

// Update updates an existing Symbol
func (r *MarketRepository) Update(ctx context.Context, symbol *model.Symbol) error {
	entity := r.symbolToEntity(symbol)

	result := r.db.WithContext(ctx).Where("symbol = ?", symbol.Symbol).Updates(&entity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", symbol.Symbol).Msg("Failed to update symbol")
		return fmt.Errorf("failed to update symbol: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		r.logger.Warn().Str("symbol", symbol.Symbol).Msg("No symbol found to update")
		return apperror.ErrNotFound
	}

	r.logger.Info().Str("symbol", symbol.Symbol).Msg("Symbol updated successfully")
	return nil
}

// GetSymbolsByStatus returns symbols by status with pagination
func (r *MarketRepository) GetSymbolsByStatus(ctx context.Context, status string, limit int, offset int) ([]*model.Symbol, error) {
	var entities []SymbolEntity
	result := r.db.WithContext(ctx).Where("status = ?", status).Limit(limit).Offset(offset).Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("status", status).Msg("Failed to get symbols by status")
		return nil, fmt.Errorf("failed to get symbols by status: %w", result.Error)
	}
	symbols := make([]*model.Symbol, 0, len(entities))
	for _, entity := range entities {
		symbols = append(symbols, r.symbolToDomain(&entity))
	}
	return symbols, nil
}

// Delete removes a Symbol
func (r *MarketRepository) Delete(ctx context.Context, symbol string) error {
	result := r.db.WithContext(ctx).Where("symbol = ?", symbol).Delete(&SymbolEntity{})
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to delete symbol")
		return fmt.Errorf("failed to delete symbol: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		r.logger.Warn().Str("symbol", symbol).Msg("No symbol found to delete")
		return apperror.ErrNotFound
	}

	r.logger.Info().Str("symbol", symbol).Msg("Symbol deleted successfully")
	return nil
}

// Helper methods to convert between domain model and database entity

func (r *MarketRepository) tickerToEntity(ticker *model.Ticker) TickerEntity {
	return TickerEntity{
		ID:            ticker.ID,
		Symbol:        ticker.Symbol,
		Price:         ticker.LastPrice,
		Volume:        ticker.Volume,
		High24h:       ticker.HighPrice,
		Low24h:        ticker.LowPrice,
		PriceChange:   ticker.PriceChange,
		PercentChange: ticker.PriceChangePercent,
		LastUpdated:   ticker.Timestamp,
		Exchange:      ticker.Exchange,
	}
}

func (r *MarketRepository) tickerToDomain(entity *TickerEntity) *model.Ticker {
	return &model.Ticker{
		ID:                 entity.ID,
		Symbol:             entity.Symbol,
		LastPrice:          entity.Price,
		Volume:             entity.Volume,
		HighPrice:          entity.High24h,
		LowPrice:           entity.Low24h,
		PriceChange:        entity.PriceChange,
		PriceChangePercent: entity.PercentChange,
		Timestamp:          entity.LastUpdated,
		Exchange:           entity.Exchange,
	}
}

func (r *MarketRepository) candleToEntity(kline *model.Kline) CandleEntity {
	return CandleEntity{
		Symbol:      kline.Symbol,
		Exchange:    kline.Exchange,
		Interval:    string(kline.Interval),
		OpenTime:    kline.OpenTime,
		CloseTime:   kline.CloseTime,
		Open:        kline.Open,
		High:        kline.High,
		Low:         kline.Low,
		Close:       kline.Close,
		Volume:      kline.Volume,
		QuoteVolume: kline.QuoteVolume,
		TradeCount:  kline.TradeCount,
		Complete:    kline.Complete,
	}
}

func (r *MarketRepository) candleToDomain(entity *CandleEntity) *model.Kline {
	return &model.Kline{
		Symbol:      entity.Symbol,
		Exchange:    entity.Exchange,
		Interval:    model.KlineInterval(entity.Interval),
		OpenTime:    entity.OpenTime,
		CloseTime:   entity.CloseTime,
		Open:        entity.Open,
		High:        entity.High,
		Low:         entity.Low,
		Close:       entity.Close,
		Volume:      entity.Volume,
		QuoteVolume: entity.QuoteVolume,
		TradeCount:  entity.TradeCount,
		Complete:    entity.Complete,
	}
}

func (r *MarketRepository) orderBookToEntity(orderbook *model.OrderBook) (OrderBookEntity, []OrderBookEntryEntity) {
	entity := OrderBookEntity{
		Symbol:       orderbook.Symbol,
		Exchange:     "MEXC", // Default to MEXC since model.OrderBook doesn't have Exchange field
		LastUpdated:  orderbook.Timestamp,
		SequenceNum:  0, // Not available in model.OrderBook
		LastUpdateID: orderbook.LastUpdateID,
	}

	entries := make([]OrderBookEntryEntity, 0, len(orderbook.Bids)+len(orderbook.Asks))

	// Add bid entries
	for _, bid := range orderbook.Bids {
		entries = append(entries, OrderBookEntryEntity{
			Type:     "bid",
			Price:    bid.Price,
			Quantity: bid.Quantity,
		})
	}

	// Add ask entries
	for _, ask := range orderbook.Asks {
		entries = append(entries, OrderBookEntryEntity{
			Type:     "ask",
			Price:    ask.Price,
			Quantity: ask.Quantity,
		})
	}

	return entity, entries
}

func (r *MarketRepository) orderBookToDomain(entity *OrderBookEntity, entries []OrderBookEntryEntity) *model.OrderBook {
	orderbook := &model.OrderBook{
		Symbol:       entity.Symbol,
		LastUpdateID: entity.LastUpdateID,
		Bids:         make([]model.OrderBookEntry, 0),
		Asks:         make([]model.OrderBookEntry, 0),
		Timestamp:    entity.LastUpdated,
	}

	for _, entry := range entries {
		bookEntry := model.OrderBookEntry{
			Price:    entry.Price,
			Quantity: entry.Quantity,
		}

		if entry.Type == "bid" {
			orderbook.Bids = append(orderbook.Bids, bookEntry)
		} else if entry.Type == "ask" {
			orderbook.Asks = append(orderbook.Asks, bookEntry)
		}
	}

	return orderbook
}

func (r *MarketRepository) symbolToEntity(symbol *model.Symbol) SymbolEntity {
	return SymbolEntity{
		Symbol:            symbol.Symbol,
		BaseAsset:         symbol.BaseAsset,
		QuoteAsset:        symbol.QuoteAsset,
		Exchange:          symbol.Exchange,
		Status:            string(symbol.Status),
		MinPrice:          symbol.MinPrice,
		MaxPrice:          symbol.MaxPrice,
		PricePrecision:    symbol.PricePrecision,
		MinQty:            symbol.MinQuantity,
		MaxQty:            symbol.MaxQuantity,
		QtyPrecision:      symbol.QuantityPrecision,
		AllowedOrderTypes: strings.Join(symbol.AllowedOrderTypes, ","),
	}
}

func (r *MarketRepository) symbolToDomain(entity *SymbolEntity) *model.Symbol {
	var allowedOrderTypes []string
	if entity.AllowedOrderTypes != "" {
		allowedOrderTypes = strings.Split(entity.AllowedOrderTypes, ",")
	}

	return &model.Symbol{
		Symbol:            entity.Symbol,
		BaseAsset:         entity.BaseAsset,
		QuoteAsset:        entity.QuoteAsset,
		Exchange:          entity.Exchange,
		Status:            model.SymbolStatus(entity.Status),
		MinPrice:          entity.MinPrice,
		MaxPrice:          entity.MaxPrice,
		PricePrecision:    entity.PricePrecision,
		MinQuantity:       entity.MinQty,
		MaxQuantity:       entity.MaxQty,
		QuantityPrecision: entity.QtyPrecision,
		AllowedOrderTypes: allowedOrderTypes,
		CreatedAt:         entity.CreatedAt,
		UpdatedAt:         entity.UpdatedAt,
	}
}

// GetOrderBook retrieves the latest order book for a symbol from a specific exchange
func (r *MarketRepository) GetOrderBook(ctx context.Context, symbol, exchange string, depth int) (*model.OrderBook, error) {
	var entity OrderBookEntity

	result := r.db.WithContext(ctx).
		Where("symbol = ? AND exchange = ?", symbol, exchange).
		Order("last_updated DESC").
		First(&entity)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			r.logger.Info().Str("symbol", symbol).Str("exchange", exchange).Msg("Order book not found")
			return nil, apperror.ErrNotFound
		}
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get order book")
		return nil, fmt.Errorf("failed to get order book: %w", result.Error)
	}

	// Get order book entries
	var entries []OrderBookEntryEntity

	// Apply depth limit if provided (greater than 0)
	query := r.db.WithContext(ctx).Where("order_book_id = ?", entity.ID)

	if depth > 0 {
		// Get top "depth" bids ordered by price descending (highest first)
		var bidEntries []OrderBookEntryEntity
		bidQuery := query.Where("type = ?", "bid").Order("price DESC")
		if depth > 0 {
			bidQuery = bidQuery.Limit(depth)
		}
		bidQuery.Find(&bidEntries)

		// Get top "depth" asks ordered by price ascending (lowest first)
		var askEntries []OrderBookEntryEntity
		askQuery := query.Where("type = ?", "ask").Order("price ASC")
		if depth > 0 {
			askQuery = askQuery.Limit(depth)
		}
		askQuery.Find(&askEntries)

		// Combine bid and ask entries
		entries = append(bidEntries, askEntries...)
	} else {
		// If depth is 0 or negative, get all entries
		result = query.Find(&entries)
		if result.Error != nil {
			r.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get order book entries")
			return nil, fmt.Errorf("failed to get order book entries: %w", result.Error)
		}
	}

	return r.orderBookToDomain(&entity, entries), nil
}

// Legacy methods for backward compatibility

// SaveTickerLegacy stores a ticker in the database using the legacy model
func (r *MarketRepository) SaveTickerLegacy(ctx context.Context, ticker *market.Ticker) error {
	// Convert legacy model to canonical model
	canonicalTicker := compat.ConvertMarketTickerToTicker(ticker)
	// Use the canonical implementation
	return r.SaveTicker(ctx, canonicalTicker)
}

// GetTickerLegacy retrieves the latest ticker for a symbol from a specific exchange using the legacy model
func (r *MarketRepository) GetTickerLegacy(ctx context.Context, symbol, exchange string) (*market.Ticker, error) {
	// Use the canonical implementation and convert the result
	canonicalTicker, err := r.GetTicker(ctx, symbol, exchange)
	if err != nil {
		return nil, err
	}
	return compat.ConvertTickerToMarketTicker(canonicalTicker), nil
}

// GetAllTickersLegacy retrieves all latest tickers from a specific exchange using the legacy model
func (r *MarketRepository) GetAllTickersLegacy(ctx context.Context, exchange string) ([]*market.Ticker, error) {
	// Use the canonical implementation and convert the results
	canonicalTickers, err := r.GetAllTickers(ctx, exchange)
	if err != nil {
		return nil, err
	}

	legacyTickers := make([]*market.Ticker, len(canonicalTickers))
	for i, ticker := range canonicalTickers {
		legacyTickers[i] = compat.ConvertTickerToMarketTicker(ticker)
	}

	return legacyTickers, nil
}

// GetTickerHistoryLegacy retrieves ticker history for a symbol within a time range using the legacy model
func (r *MarketRepository) GetTickerHistoryLegacy(ctx context.Context, symbol, exchange string, start, end time.Time) ([]*market.Ticker, error) {
	// Use the canonical implementation and convert the results
	canonicalTickers, err := r.GetTickerHistory(ctx, symbol, exchange, start, end)
	if err != nil {
		return nil, err
	}

	legacyTickers := make([]*market.Ticker, len(canonicalTickers))
	for i, ticker := range canonicalTickers {
		legacyTickers[i] = compat.ConvertTickerToMarketTicker(ticker)
	}

	return legacyTickers, nil
}

// SaveCandleLegacy stores a candle in the database using the legacy model
func (r *MarketRepository) SaveCandleLegacy(ctx context.Context, candle *market.Candle) error {
	// Convert legacy model to canonical model
	canonicalKline := &model.Kline{
		Symbol:      candle.Symbol,
		Exchange:    candle.Exchange,
		Interval:    model.KlineInterval(candle.Interval),
		OpenTime:    candle.OpenTime,
		CloseTime:   candle.CloseTime,
		Open:        candle.Open,
		High:        candle.High,
		Low:         candle.Low,
		Close:       candle.Close,
		Volume:      candle.Volume,
		QuoteVolume: candle.QuoteVolume,
		TradeCount:  candle.TradeCount,
		Complete:    candle.Complete,
	}

	// Use the canonical implementation
	return r.SaveCandle(ctx, canonicalKline)
}

// SaveCandlesLegacy stores multiple candles in the database using the legacy model
func (r *MarketRepository) SaveCandlesLegacy(ctx context.Context, candles []*market.Candle) error {
	// Convert legacy models to canonical models
	canonicalKlines := make([]*model.Kline, len(candles))
	for i, candle := range candles {
		canonicalKlines[i] = &model.Kline{
			Symbol:      candle.Symbol,
			Exchange:    candle.Exchange,
			Interval:    model.KlineInterval(candle.Interval),
			OpenTime:    candle.OpenTime,
			CloseTime:   candle.CloseTime,
			Open:        candle.Open,
			High:        candle.High,
			Low:         candle.Low,
			Close:       candle.Close,
			Volume:      candle.Volume,
			QuoteVolume: candle.QuoteVolume,
			TradeCount:  candle.TradeCount,
			Complete:    candle.Complete,
		}
	}

	// Use the canonical implementation
	return r.SaveCandles(ctx, canonicalKlines)
}

// GetCandleLegacy retrieves a specific candle for a symbol, interval, and time using the legacy model
func (r *MarketRepository) GetCandleLegacy(ctx context.Context, symbol, exchange string, interval market.Interval, openTime time.Time) (*market.Candle, error) {
	// Use the canonical implementation and convert the result
	canonicalKline, err := r.GetCandle(ctx, symbol, exchange, model.KlineInterval(interval), openTime)
	if err != nil {
		return nil, err
	}

	return &market.Candle{
		Symbol:      canonicalKline.Symbol,
		Exchange:    canonicalKline.Exchange,
		Interval:    market.Interval(canonicalKline.Interval),
		OpenTime:    canonicalKline.OpenTime,
		CloseTime:   canonicalKline.CloseTime,
		Open:        canonicalKline.Open,
		High:        canonicalKline.High,
		Low:         canonicalKline.Low,
		Close:       canonicalKline.Close,
		Volume:      canonicalKline.Volume,
		QuoteVolume: canonicalKline.QuoteVolume,
		TradeCount:  canonicalKline.TradeCount,
		Complete:    canonicalKline.Complete,
	}, nil
}

// GetCandlesLegacy retrieves candles for a symbol within a time range using the legacy model
func (r *MarketRepository) GetCandlesLegacy(ctx context.Context, symbol, exchange string, interval market.Interval, start, end time.Time, limit int) ([]*market.Candle, error) {
	// Use the canonical implementation and convert the results
	canonicalKlines, err := r.GetCandles(ctx, symbol, exchange, model.KlineInterval(interval), start, end, limit)
	if err != nil {
		return nil, err
	}

	legacyCandles := make([]*market.Candle, len(canonicalKlines))
	for i, kline := range canonicalKlines {
		legacyCandles[i] = &market.Candle{
			Symbol:      kline.Symbol,
			Exchange:    kline.Exchange,
			Interval:    market.Interval(kline.Interval),
			OpenTime:    kline.OpenTime,
			CloseTime:   kline.CloseTime,
			Open:        kline.Open,
			High:        kline.High,
			Low:         kline.Low,
			Close:       kline.Close,
			Volume:      kline.Volume,
			QuoteVolume: kline.QuoteVolume,
			TradeCount:  kline.TradeCount,
			Complete:    kline.Complete,
		}
	}

	return legacyCandles, nil
}

// GetLatestCandleLegacy retrieves the most recent candle for a symbol and interval using the legacy model
func (r *MarketRepository) GetLatestCandleLegacy(ctx context.Context, symbol, exchange string, interval market.Interval) (*market.Candle, error) {
	// Use the canonical implementation and convert the result
	canonicalKline, err := r.GetLatestCandle(ctx, symbol, exchange, model.KlineInterval(interval))
	if err != nil {
		return nil, err
	}

	return &market.Candle{
		Symbol:      canonicalKline.Symbol,
		Exchange:    canonicalKline.Exchange,
		Interval:    market.Interval(canonicalKline.Interval),
		OpenTime:    canonicalKline.OpenTime,
		CloseTime:   canonicalKline.CloseTime,
		Open:        canonicalKline.Open,
		High:        canonicalKline.High,
		Low:         canonicalKline.Low,
		Close:       canonicalKline.Close,
		Volume:      canonicalKline.Volume,
		QuoteVolume: canonicalKline.QuoteVolume,
		TradeCount:  canonicalKline.TradeCount,
		Complete:    canonicalKline.Complete,
	}, nil
}

// GetLatestTickersLegacy retrieves the latest tickers for all symbols using the legacy model
func (r *MarketRepository) GetLatestTickersLegacy(ctx context.Context, limit int) ([]*market.Ticker, error) {
	// Use the canonical implementation and convert the results
	canonicalTickers, err := r.GetLatestTickers(ctx, limit)
	if err != nil {
		return nil, err
	}

	legacyTickers := make([]*market.Ticker, len(canonicalTickers))
	for i, ticker := range canonicalTickers {
		legacyTickers[i] = compat.ConvertTickerToMarketTicker(ticker)
	}

	return legacyTickers, nil
}

// GetTickersBySymbolLegacy retrieves tickers for a specific symbol with optional time range using the legacy model
func (r *MarketRepository) GetTickersBySymbolLegacy(ctx context.Context, symbol string, limit int) ([]*market.Ticker, error) {
	// Use the canonical implementation and convert the results
	canonicalTickers, err := r.GetTickersBySymbol(ctx, symbol, limit)
	if err != nil {
		return nil, err
	}

	legacyTickers := make([]*market.Ticker, len(canonicalTickers))
	for i, ticker := range canonicalTickers {
		legacyTickers[i] = compat.ConvertTickerToMarketTicker(ticker)
	}

	return legacyTickers, nil
}

// GetOrderBookLegacy retrieves the order book for a symbol using the legacy model
func (r *MarketRepository) GetOrderBookLegacy(ctx context.Context, symbol, exchange string, depth int) (*market.OrderBook, error) {
	// Use the canonical implementation and convert the result
	canonicalOrderBook, err := r.GetOrderBook(ctx, symbol, exchange, depth)
	if err != nil {
		return nil, err
	}

	// Convert to legacy model
	legacyOrderBook := compat.ConvertOrderBookToMarketOrderBook(canonicalOrderBook)
	// Set the exchange since it's not part of the model.OrderBook
	legacyOrderBook.Exchange = exchange

	return legacyOrderBook, nil
}

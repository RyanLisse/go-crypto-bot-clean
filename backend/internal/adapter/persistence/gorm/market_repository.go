package gorm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
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
func (r *MarketRepository) SaveTicker(ctx context.Context, ticker *market.Ticker) error {
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
func (r *MarketRepository) GetTicker(ctx context.Context, symbol, exchange string) (*market.Ticker, error) {
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
func (r *MarketRepository) GetAllTickers(ctx context.Context, exchange string) ([]*market.Ticker, error) {
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

	tickers := make([]*market.Ticker, len(entities))
	for i, entity := range entities {
		tickers[i] = r.tickerToDomain(&entity)
	}

	return tickers, nil
}

// GetTickerHistory retrieves ticker history for a symbol within a time range
func (r *MarketRepository) GetTickerHistory(ctx context.Context, symbol, exchange string, start, end time.Time) ([]*market.Ticker, error) {
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

	tickers := make([]*market.Ticker, len(entities))
	for i, entity := range entities {
		tickers[i] = r.tickerToDomain(&entity)
	}

	return tickers, nil
}

// SaveCandle stores a candle in the database
func (r *MarketRepository) SaveCandle(ctx context.Context, candle *market.Candle) error {
	entity := r.candleToEntity(candle)

	// Try to find an existing candle with the same symbol, exchange, interval, and open time
	var existing CandleEntity
	result := r.db.WithContext(ctx).
		Where("symbol = ? AND exchange = ? AND interval = ? AND open_time = ?",
			candle.Symbol, candle.Exchange, candle.Interval, candle.OpenTime).
		First(&existing)

	// If the candle exists, update it; otherwise, create a new one
	if result.Error == nil {
		entity.ID = existing.ID
	}

	result = r.db.WithContext(ctx).Save(&entity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", candle.Symbol).Msg("Failed to save candle")
		return fmt.Errorf("failed to save candle: %w", result.Error)
	}

	r.logger.Info().Str("symbol", candle.Symbol).Str("interval", string(candle.Interval)).Msg("Candle saved successfully")
	return nil
}

// SaveCandles stores multiple candles in the database
func (r *MarketRepository) SaveCandles(ctx context.Context, candles []*market.Candle) error {
	if len(candles) == 0 {
		return nil
	}

	entities := make([]CandleEntity, len(candles))
	for i, candle := range candles {
		entities[i] = r.candleToEntity(candle)
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
			r.logger.Error().Err(result.Error).Str("symbol", candles[i].Symbol).Msg("Failed to save candle in batch")
			return fmt.Errorf("failed to save candle in batch: %w", result.Error)
		}
	}

	if err := tx.Commit().Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to commit transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Info().Int("count", len(candles)).Msg("Successfully saved batch of candles")
	return nil
}

// GetCandle retrieves a specific candle for a symbol, interval, and time
func (r *MarketRepository) GetCandle(ctx context.Context, symbol, exchange string, interval market.Interval, openTime time.Time) (*market.Candle, error) {
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

// GetCandles retrieves candles for a symbol within a time range
func (r *MarketRepository) GetCandles(ctx context.Context, symbol, exchange string, interval market.Interval, start, end time.Time, limit int) ([]*market.Candle, error) {
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

	candles := make([]*market.Candle, len(entities))
	for i, entity := range entities {
		candles[i] = r.candleToDomain(&entity)
	}

	return candles, nil
}

// GetLatestCandle retrieves the most recent candle for a symbol and interval
func (r *MarketRepository) GetLatestCandle(ctx context.Context, symbol, exchange string, interval market.Interval) (*market.Candle, error) {
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
func (r *MarketRepository) GetLatestTickers(ctx context.Context, limit int) ([]*market.Ticker, error) {
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

	tickers := make([]*market.Ticker, len(entities))
	for i, entity := range entities {
		tickers[i] = r.tickerToDomain(&entity)
	}

	return tickers, nil
}

// GetTickersBySymbol retrieves tickers for a specific symbol with optional time range
func (r *MarketRepository) GetTickersBySymbol(ctx context.Context, symbol string, limit int) ([]*market.Ticker, error) {
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

	tickers := make([]*market.Ticker, len(entities))
	for i, entity := range entities {
		tickers[i] = r.tickerToDomain(&entity)
	}

	return tickers, nil
}

// Symbol Repository implementation

// Create stores a new Symbol
func (r *MarketRepository) Create(ctx context.Context, symbol *market.Symbol) error {
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
func (r *MarketRepository) GetBySymbol(ctx context.Context, symbol string) (*market.Symbol, error) {
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
func (r *MarketRepository) GetByExchange(ctx context.Context, exchange string) ([]*market.Symbol, error) {
	var entities []SymbolEntity

	result := r.db.WithContext(ctx).Where("exchange = ?", exchange).Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("exchange", exchange).Msg("Failed to get symbols by exchange")
		return nil, fmt.Errorf("failed to get symbols by exchange: %w", result.Error)
	}

	symbols := make([]*market.Symbol, len(entities))
	for i, entity := range entities {
		symbols[i] = r.symbolToDomain(&entity)
	}

	r.logger.Info().Str("exchange", exchange).Int("count", len(symbols)).Msg("Retrieved symbols by exchange")
	return symbols, nil
}

// GetAll returns all available Symbols
func (r *MarketRepository) GetAll(ctx context.Context) ([]*market.Symbol, error) {
	var entities []SymbolEntity

	result := r.db.WithContext(ctx).Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Msg("Failed to get all symbols")
		return nil, fmt.Errorf("failed to get all symbols: %w", result.Error)
	}

	symbols := make([]*market.Symbol, len(entities))
	for i, entity := range entities {
		symbols[i] = r.symbolToDomain(&entity)
	}

	r.logger.Info().Int("count", len(symbols)).Msg("Retrieved all symbols")
	return symbols, nil
}

// Update updates an existing Symbol
func (r *MarketRepository) Update(ctx context.Context, symbol *market.Symbol) error {
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

func (r *MarketRepository) tickerToEntity(ticker *market.Ticker) TickerEntity {
	return TickerEntity{
		ID:            ticker.ID,
		Symbol:        ticker.Symbol,
		Price:         ticker.Price,
		Volume:        ticker.Volume,
		High24h:       ticker.High24h,
		Low24h:        ticker.Low24h,
		PriceChange:   ticker.PriceChange,
		PercentChange: ticker.PercentChange,
		LastUpdated:   ticker.LastUpdated,
		Exchange:      ticker.Exchange,
	}
}

func (r *MarketRepository) tickerToDomain(entity *TickerEntity) *market.Ticker {
	return &market.Ticker{
		ID:            entity.ID,
		Symbol:        entity.Symbol,
		Price:         entity.Price,
		Volume:        entity.Volume,
		High24h:       entity.High24h,
		Low24h:        entity.Low24h,
		PriceChange:   entity.PriceChange,
		PercentChange: entity.PercentChange,
		LastUpdated:   entity.LastUpdated,
		Exchange:      entity.Exchange,
	}
}

func (r *MarketRepository) candleToEntity(candle *market.Candle) CandleEntity {
	return CandleEntity{
		Symbol:      candle.Symbol,
		Exchange:    candle.Exchange,
		Interval:    string(candle.Interval),
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

func (r *MarketRepository) candleToDomain(entity *CandleEntity) *market.Candle {
	return &market.Candle{
		Symbol:      entity.Symbol,
		Exchange:    entity.Exchange,
		Interval:    market.Interval(entity.Interval),
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

func (r *MarketRepository) orderBookToEntity(orderbook *market.OrderBook) (OrderBookEntity, []OrderBookEntryEntity) {
	entity := OrderBookEntity{
		Symbol:       orderbook.Symbol,
		Exchange:     orderbook.Exchange,
		LastUpdated:  orderbook.LastUpdated,
		SequenceNum:  orderbook.SequenceNum,
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

func (r *MarketRepository) orderBookToDomain(entity *OrderBookEntity, entries []OrderBookEntryEntity) *market.OrderBook {
	orderbook := &market.OrderBook{
		Symbol:       entity.Symbol,
		Exchange:     entity.Exchange,
		LastUpdated:  entity.LastUpdated,
		SequenceNum:  entity.SequenceNum,
		LastUpdateID: entity.LastUpdateID,
		Bids:         make([]market.OrderBookEntry, 0),
		Asks:         make([]market.OrderBookEntry, 0),
	}

	for _, entry := range entries {
		bookEntry := market.OrderBookEntry{
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

func (r *MarketRepository) symbolToEntity(symbol *market.Symbol) SymbolEntity {
	return SymbolEntity{
		Symbol:            symbol.Symbol,
		BaseAsset:         symbol.BaseAsset,
		QuoteAsset:        symbol.QuoteAsset,
		Exchange:          symbol.Exchange,
		Status:            symbol.Status,
		MinPrice:          symbol.MinPrice,
		MaxPrice:          symbol.MaxPrice,
		PricePrecision:    symbol.PricePrecision,
		MinQty:            symbol.MinQty,
		MaxQty:            symbol.MaxQty,
		QtyPrecision:      symbol.QtyPrecision,
		AllowedOrderTypes: strings.Join(symbol.AllowedOrderTypes, ","),
	}
}

func (r *MarketRepository) symbolToDomain(entity *SymbolEntity) *market.Symbol {
	var allowedOrderTypes []string
	if entity.AllowedOrderTypes != "" {
		allowedOrderTypes = strings.Split(entity.AllowedOrderTypes, ",")
	}

	return &market.Symbol{
		Symbol:            entity.Symbol,
		BaseAsset:         entity.BaseAsset,
		QuoteAsset:        entity.QuoteAsset,
		Exchange:          entity.Exchange,
		Status:            entity.Status,
		MinPrice:          entity.MinPrice,
		MaxPrice:          entity.MaxPrice,
		PricePrecision:    entity.PricePrecision,
		MinQty:            entity.MinQty,
		MaxQty:            entity.MaxQty,
		QtyPrecision:      entity.QtyPrecision,
		AllowedOrderTypes: allowedOrderTypes,
		CreatedAt:         entity.CreatedAt,
		UpdatedAt:         entity.UpdatedAt,
	}
}

// GetOrderBook retrieves the latest order book for a symbol from a specific exchange
func (r *MarketRepository) GetOrderBook(ctx context.Context, symbol, exchange string, depth int) (*market.OrderBook, error) {
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

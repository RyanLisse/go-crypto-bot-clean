package gorm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/compat"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// Ensure MarketRepositoryDirect implements the proper interfaces
var _ port.MarketRepository = (*MarketRepositoryDirect)(nil)
var _ port.SymbolRepository = (*MarketRepositoryDirect)(nil)

// MarketRepositoryDirect implements the port.MarketRepository interface using GORM
// with direct use of canonical model types (model.Ticker, model.OrderBook)
type MarketRepositoryDirect struct {
	db     *gorm.DB
	logger *zerolog.Logger
	legacy *MarketRepository // For legacy method implementations
}

// NewMarketRepositoryDirect creates a new MarketRepositoryDirect
func NewMarketRepositoryDirect(db *gorm.DB, logger *zerolog.Logger) *MarketRepositoryDirect {
	return &MarketRepositoryDirect{
		db:     db,
		logger: logger,
		legacy: NewMarketRepository(db, logger),
	}
}

// tickerToEntity converts a model.Ticker to a database entity
func (r *MarketRepositoryDirect) tickerToEntity(ticker *model.Ticker) *entity.Ticker {
	return &entity.Ticker{
		Symbol:        ticker.Symbol,
		Exchange:      ticker.Exchange,
		Price:         ticker.LastPrice,
		PriceChange:   ticker.PriceChange,
		PercentChange: ticker.PriceChangePercent,
		High24h:       ticker.HighPrice,
		Low24h:        ticker.LowPrice,
		Volume:        ticker.Volume,
		LastUpdated:   ticker.Timestamp,
	}
}

// entityToTicker converts a database entity to a model.Ticker
func (r *MarketRepositoryDirect) entityToTicker(entity *entity.Ticker) *model.Ticker {
	return &model.Ticker{
		Symbol:             entity.Symbol,
		Exchange:           entity.Exchange,
		LastPrice:          entity.Price,
		PriceChange:        entity.PriceChange,
		PriceChangePercent: entity.PercentChange,
		HighPrice:          entity.High24h,
		LowPrice:           entity.Low24h,
		Volume:             entity.Volume,
		Timestamp:          entity.LastUpdated,
	}
}

// orderBookToEntity converts a model.OrderBook to a database entity
func (r *MarketRepositoryDirect) orderBookToEntity(orderBook *model.OrderBook, exchange string) (*entity.OrderBook, error) {
	// Convert bids to JSON
	bidsJSON, err := json.Marshal(orderBook.Bids)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal bids: %w", err)
	}

	// Convert asks to JSON
	asksJSON, err := json.Marshal(orderBook.Asks)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal asks: %w", err)
	}

	return &entity.OrderBook{
		Symbol:      orderBook.Symbol,
		Exchange:    exchange,
		BidsJSON:    bidsJSON,
		AsksJSON:    asksJSON,
		LastUpdated: orderBook.Timestamp,
	}, nil
}

// entityToOrderBook converts a database entity to a model.OrderBook
func (r *MarketRepositoryDirect) entityToOrderBook(entity *entity.OrderBook) (*model.OrderBook, error) {
	// Parse bids from JSON
	var bids []model.OrderBookEntry
	if err := json.Unmarshal(entity.BidsJSON, &bids); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bids: %w", err)
	}

	// Parse asks from JSON
	var asks []model.OrderBookEntry
	if err := json.Unmarshal(entity.AsksJSON, &asks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal asks: %w", err)
	}

	return &model.OrderBook{
		Symbol:    entity.Symbol,
		Bids:      bids,
		Asks:      asks,
		Timestamp: entity.LastUpdated,
	}, nil
}

// SaveTicker stores a ticker in the database
func (r *MarketRepositoryDirect) SaveTicker(ctx context.Context, ticker *model.Ticker) error {
	entity := r.tickerToEntity(ticker)

	result := r.db.WithContext(ctx).Save(entity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", ticker.Symbol).Msg("Failed to save ticker")
		return fmt.Errorf("failed to save ticker: %w", result.Error)
	}

	r.logger.Info().Str("symbol", ticker.Symbol).Str("exchange", ticker.Exchange).Msg("Ticker saved successfully")
	return nil
}

// GetTicker retrieves the latest ticker for a symbol from a specific exchange
func (r *MarketRepositoryDirect) GetTicker(ctx context.Context, symbol, exchange string) (*model.Ticker, error) {
	var entity entity.Ticker

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

	return r.entityToTicker(&entity), nil
}

// GetAllTickers retrieves all latest tickers from a specific exchange
func (r *MarketRepositoryDirect) GetAllTickers(ctx context.Context, exchange string) ([]*model.Ticker, error) {
	var entities []entity.Ticker

	// Using a subquery to get the latest ticker for each symbol
	subQuery := r.db.Model(&entity.Ticker{}).
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
		tickers[i] = r.entityToTicker(&entity)
	}

	return tickers, nil
}

// GetTickerHistory retrieves ticker history for a symbol within a time range
func (r *MarketRepositoryDirect) GetTickerHistory(ctx context.Context, symbol, exchange string, start, end time.Time) ([]*model.Ticker, error) {
	var entities []entity.Ticker

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
		tickers[i] = r.entityToTicker(&entity)
	}

	return tickers, nil
}

// klineToEntity converts a model.Kline to a database entity
func (r *MarketRepositoryDirect) klineToEntity(kline *model.Kline) CandleEntity {
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

// entityToKline converts a database entity to a model.Kline
func (r *MarketRepositoryDirect) entityToKline(entity *CandleEntity) *model.Kline {
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

// SaveKline stores a kline/candle in the database
func (r *MarketRepositoryDirect) SaveKline(ctx context.Context, kline *model.Kline) error {
	entity := r.klineToEntity(kline)

	// Try to find an existing kline with the same symbol, interval, and open time
	var existing CandleEntity
	result := r.db.WithContext(ctx).
		Where("symbol = ? AND interval = ? AND open_time = ?",
			kline.Symbol, string(kline.Interval), kline.OpenTime).
		First(&existing)

	// If the kline exists, update it; otherwise, create a new one
	if result.Error == nil {
		entity.ID = existing.ID
	}

	result = r.db.WithContext(ctx).Save(&entity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", kline.Symbol).Msg("Failed to save kline")
		return fmt.Errorf("failed to save kline: %w", result.Error)
	}

	r.logger.Info().Str("symbol", kline.Symbol).Str("interval", string(kline.Interval)).Msg("Kline saved successfully")
	return nil
}

// SaveKlines stores multiple klines/candles in the database
func (r *MarketRepositoryDirect) SaveKlines(ctx context.Context, klines []*model.Kline) error {
	if len(klines) == 0 {
		return nil
	}

	entities := make([]CandleEntity, len(klines))
	for i, kline := range klines {
		entities[i] = r.klineToEntity(kline)
	}

	// Use a transaction to save all klines
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		r.logger.Error().Err(tx.Error).Msg("Failed to begin transaction")
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Create or update each kline
	for i, entity := range entities {
		var existing CandleEntity
		result := tx.Where("symbol = ? AND interval = ? AND open_time = ?",
			entity.Symbol, entity.Interval, entity.OpenTime).
			First(&existing)

		if result.Error == nil {
			entity.ID = existing.ID
		}

		result = tx.Save(&entity)
		if result.Error != nil {
			tx.Rollback()
			r.logger.Error().Err(result.Error).Str("symbol", klines[i].Symbol).Msg("Failed to save kline in batch")
			return fmt.Errorf("failed to save kline in batch: %w", result.Error)
		}
	}

	if err := tx.Commit().Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to commit transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Info().Int("count", len(klines)).Msg("Successfully saved batch of klines")
	return nil
}

// GetKline retrieves a specific kline/candle for a symbol, interval, and time
func (r *MarketRepositoryDirect) GetKline(ctx context.Context, symbol, exchange string, interval model.KlineInterval, openTime time.Time) (*model.Kline, error) {
	var entity CandleEntity

	result := r.db.WithContext(ctx).
		Where("symbol = ? AND interval = ? AND open_time = ?", symbol, string(interval), openTime).
		First(&entity)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			r.logger.Info().Str("symbol", symbol).Str("interval", string(interval)).Time("openTime", openTime).Msg("Kline not found")
			return nil, apperror.ErrNotFound
		}
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get kline")
		return nil, fmt.Errorf("failed to get kline: %w", result.Error)
	}

	return r.entityToKline(&entity), nil
}

// GetKlines retrieves klines/candles for a symbol within a time range
func (r *MarketRepositoryDirect) GetKlines(ctx context.Context, symbol, exchange string, interval model.KlineInterval, start, end time.Time, limit int) ([]*model.Kline, error) {
	var entities []CandleEntity

	query := r.db.WithContext(ctx).
		Where("symbol = ? AND interval = ? AND open_time BETWEEN ? AND ?",
			symbol, string(interval), start, end).
		Order("open_time ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	result := query.Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get klines")
		return nil, fmt.Errorf("failed to get klines: %w", result.Error)
	}

	klines := make([]*model.Kline, len(entities))
	for i, entity := range entities {
		klines[i] = r.entityToKline(&entity)
	}

	return klines, nil
}

// GetLatestKline retrieves the most recent kline/candle for a symbol and interval
func (r *MarketRepositoryDirect) GetLatestKline(ctx context.Context, symbol, exchange string, interval model.KlineInterval) (*model.Kline, error) {
	var entity CandleEntity

	result := r.db.WithContext(ctx).
		Where("symbol = ? AND interval = ?", symbol, string(interval)).
		Order("open_time DESC").
		First(&entity)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			r.logger.Info().Str("symbol", symbol).Str("interval", string(interval)).Msg("Latest kline not found")
			return nil, apperror.ErrNotFound
		}
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get latest kline")
		return nil, fmt.Errorf("failed to get latest kline: %w", result.Error)
	}

	return r.entityToKline(&entity), nil
}

// PurgeOldData removes market data older than the specified retention period
func (r *MarketRepositoryDirect) PurgeOldData(ctx context.Context, olderThan time.Time) error {
	// Delete old tickers
	if err := r.db.WithContext(ctx).
		Where("last_updated < ?", olderThan).
		Delete(&entity.Ticker{}).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to purge old tickers")
		return fmt.Errorf("failed to purge old tickers: %w", err)
	}

	// Delete old order books
	if err := r.db.WithContext(ctx).
		Where("last_updated < ?", olderThan).
		Delete(&entity.OrderBook{}).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to purge old order books")
		return fmt.Errorf("failed to purge old order books: %w", err)
	}

	// Delete old candles
	if err := r.db.WithContext(ctx).
		Where("open_time < ?", olderThan).
		Delete(&CandleEntity{}).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to purge old candles")
		return fmt.Errorf("failed to purge old candles: %w", err)
	}

	r.logger.Info().Time("older_than", olderThan).Msg("Successfully purged old market data")
	return nil
}

// GetLatestTickers retrieves the latest tickers for all symbols
func (r *MarketRepositoryDirect) GetLatestTickers(ctx context.Context, limit int) ([]*model.Ticker, error) {
	var entities []entity.Ticker

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
		tickers[i] = r.entityToTicker(&entity)
	}

	return tickers, nil
}

// GetTickersBySymbol retrieves tickers for a specific symbol with optional time range
func (r *MarketRepositoryDirect) GetTickersBySymbol(ctx context.Context, symbol string, limit int) ([]*model.Ticker, error) {
	var entities []entity.Ticker

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
		tickers[i] = r.entityToTicker(&entity)
	}

	return tickers, nil
}

// GetOrderBook retrieves the order book for a symbol
func (r *MarketRepositoryDirect) GetOrderBook(ctx context.Context, symbol, exchange string, depth int) (*model.OrderBook, error) {
	var entity entity.OrderBook

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

	orderBook, err := r.entityToOrderBook(&entity)
	if err != nil {
		r.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to convert order book entity to domain model")
		return nil, fmt.Errorf("failed to convert order book entity to domain model: %w", err)
	}

	// Apply depth limit if specified
	if depth > 0 {
		if len(orderBook.Bids) > depth {
			orderBook.Bids = orderBook.Bids[:depth]
		}
		if len(orderBook.Asks) > depth {
			orderBook.Asks = orderBook.Asks[:depth]
		}
	}

	return orderBook, nil
}

// Symbol Repository implementation

// Create stores a new Symbol
func (r *MarketRepositoryDirect) Create(ctx context.Context, symbol *model.Symbol) error {
	// Create a new SymbolEntity
	entity := SymbolEntity{
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

	result := r.db.WithContext(ctx).Create(&entity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", symbol.Symbol).Msg("Failed to create symbol")
		return fmt.Errorf("failed to create symbol: %w", result.Error)
	}

	r.logger.Info().Str("symbol", symbol.Symbol).Msg("Symbol created successfully")
	return nil
}

// GetBySymbol returns a Symbol by its symbol string (e.g., "BTCUSDT")
func (r *MarketRepositoryDirect) GetBySymbol(ctx context.Context, symbol string) (*model.Symbol, error) {
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

	// Convert entity to domain model
	allowedOrderTypes := []string{}
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
	}, nil
}

// GetByExchange returns all Symbols from a specific exchange
func (r *MarketRepositoryDirect) GetByExchange(ctx context.Context, exchange string) ([]*model.Symbol, error) {
	var entities []SymbolEntity

	result := r.db.WithContext(ctx).Where("exchange = ?", exchange).Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("exchange", exchange).Msg("Failed to get symbols by exchange")
		return nil, fmt.Errorf("failed to get symbols by exchange: %w", result.Error)
	}

	symbols := make([]*model.Symbol, len(entities))
	for i, entity := range entities {
		allowedOrderTypes := []string{}
		if entity.AllowedOrderTypes != "" {
			allowedOrderTypes = strings.Split(entity.AllowedOrderTypes, ",")
		}

		symbols[i] = &model.Symbol{
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

	return symbols, nil
}

// GetAll returns all available Symbols
func (r *MarketRepositoryDirect) GetAll(ctx context.Context) ([]*model.Symbol, error) {
	var entities []SymbolEntity

	result := r.db.WithContext(ctx).Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Msg("Failed to get all symbols")
		return nil, fmt.Errorf("failed to get all symbols: %w", result.Error)
	}

	symbols := make([]*model.Symbol, len(entities))
	for i, entity := range entities {
		allowedOrderTypes := []string{}
		if entity.AllowedOrderTypes != "" {
			allowedOrderTypes = strings.Split(entity.AllowedOrderTypes, ",")
		}

		symbols[i] = &model.Symbol{
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

	return symbols, nil
}

// Update updates an existing Symbol
func (r *MarketRepositoryDirect) Update(ctx context.Context, symbol *model.Symbol) error {
	// Create entity from domain model
	entity := SymbolEntity{
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

	// Check if the symbol exists
	var existing SymbolEntity
	result := r.db.WithContext(ctx).Where("symbol = ?", symbol.Symbol).First(&existing)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			r.logger.Warn().Str("symbol", symbol.Symbol).Msg("Symbol not found for update")
			return apperror.ErrNotFound
		}
		r.logger.Error().Err(result.Error).Str("symbol", symbol.Symbol).Msg("Failed to check if symbol exists")
		return fmt.Errorf("failed to check if symbol exists: %w", result.Error)
	}

	// Update the symbol
	result = r.db.WithContext(ctx).Model(&SymbolEntity{}).Where("symbol = ?", symbol.Symbol).Updates(entity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", symbol.Symbol).Msg("Failed to update symbol")
		return fmt.Errorf("failed to update symbol: %w", result.Error)
	}

	r.logger.Info().Str("symbol", symbol.Symbol).Msg("Symbol updated successfully")
	return nil
}

// Delete removes a Symbol
func (r *MarketRepositoryDirect) Delete(ctx context.Context, symbol string) error {
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

// GetSymbolsByStatus returns symbols by status with pagination
func (r *MarketRepositoryDirect) GetSymbolsByStatus(ctx context.Context, status string, limit int, offset int) ([]*model.Symbol, error) {
	var entities []SymbolEntity

	query := r.db.WithContext(ctx).Where("status = ?", status)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	result := query.Find(&entities)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("status", status).Msg("Failed to get symbols by status")
		return nil, fmt.Errorf("failed to get symbols by status: %w", result.Error)
	}

	symbols := make([]*model.Symbol, len(entities))
	for i, entity := range entities {
		allowedOrderTypes := []string{}
		if entity.AllowedOrderTypes != "" {
			allowedOrderTypes = strings.Split(entity.AllowedOrderTypes, ",")
		}

		symbols[i] = &model.Symbol{
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

	return symbols, nil
}

// Legacy methods for backward compatibility
// These methods are required to implement the port.MarketRepository interface

// SaveTickerLegacy stores a ticker in the database using the legacy model
func (r *MarketRepositoryDirect) SaveTickerLegacy(ctx context.Context, ticker *market.Ticker) error {
	// Convert legacy model to canonical model
	canonicalTicker := compat.ConvertMarketTickerToTicker(ticker)
	// Use the canonical implementation
	return r.SaveTicker(ctx, canonicalTicker)
}

// GetTickerLegacy retrieves the latest ticker for a symbol from a specific exchange using the legacy model
func (r *MarketRepositoryDirect) GetTickerLegacy(ctx context.Context, symbol, exchange string) (*market.Ticker, error) {
	// Use the canonical implementation and convert the result
	canonicalTicker, err := r.GetTicker(ctx, symbol, exchange)
	if err != nil {
		return nil, err
	}
	return compat.ConvertTickerToMarketTicker(canonicalTicker), nil
}

// GetAllTickersLegacy retrieves all latest tickers from a specific exchange using the legacy model
func (r *MarketRepositoryDirect) GetAllTickersLegacy(ctx context.Context, exchange string) ([]*market.Ticker, error) {
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
func (r *MarketRepositoryDirect) GetTickerHistoryLegacy(ctx context.Context, symbol, exchange string, start, end time.Time) ([]*market.Ticker, error) {
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
func (r *MarketRepositoryDirect) SaveCandleLegacy(ctx context.Context, candle *market.Candle) error {
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
	return r.SaveKline(ctx, canonicalKline)
}

// SaveCandlesLegacy stores multiple candles in the database using the legacy model
func (r *MarketRepositoryDirect) SaveCandlesLegacy(ctx context.Context, candles []*market.Candle) error {
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
	return r.SaveKlines(ctx, canonicalKlines)
}

// GetCandleLegacy retrieves a specific candle for a symbol, interval, and time using the legacy model
func (r *MarketRepositoryDirect) GetCandleLegacy(ctx context.Context, symbol, exchange string, interval market.Interval, openTime time.Time) (*market.Candle, error) {
	// Use the canonical implementation and convert the result
	canonicalKline, err := r.GetKline(ctx, symbol, exchange, model.KlineInterval(interval), openTime)
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
func (r *MarketRepositoryDirect) GetCandlesLegacy(ctx context.Context, symbol, exchange string, interval market.Interval, start, end time.Time, limit int) ([]*market.Candle, error) {
	// Use the canonical implementation and convert the results
	canonicalKlines, err := r.GetKlines(ctx, symbol, exchange, model.KlineInterval(interval), start, end, limit)
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
func (r *MarketRepositoryDirect) GetLatestCandleLegacy(ctx context.Context, symbol, exchange string, interval market.Interval) (*market.Candle, error) {
	// Use the canonical implementation and convert the result
	canonicalKline, err := r.GetLatestKline(ctx, symbol, exchange, model.KlineInterval(interval))
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
func (r *MarketRepositoryDirect) GetLatestTickersLegacy(ctx context.Context, limit int) ([]*market.Ticker, error) {
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
func (r *MarketRepositoryDirect) GetTickersBySymbolLegacy(ctx context.Context, symbol string, limit int) ([]*market.Ticker, error) {
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
func (r *MarketRepositoryDirect) GetOrderBookLegacy(ctx context.Context, symbol, exchange string, depth int) (*market.OrderBook, error) {
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

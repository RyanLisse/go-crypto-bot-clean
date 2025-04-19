package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/infrastructure/persistence/gorm/entity"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// MarketRepository implements port.MarketRepository using GORM
type MarketRepository struct {
	BaseRepository
}

// NewMarketRepository creates a new MarketRepository
func NewMarketRepository(db *gorm.DB, logger *zerolog.Logger) *MarketRepository {
	return &MarketRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// SaveTicker stores a ticker in the database
func (r *MarketRepository) SaveTicker(ctx context.Context, ticker *model.Ticker) error {
	entity := &entity.MexcTickerEntity{
		Symbol:             ticker.Symbol,
		Exchange:           ticker.Exchange,
		LastPrice:          ticker.LastPrice,
		Volume:             ticker.Volume,
		HighPrice:          ticker.HighPrice,
		LowPrice:           ticker.LowPrice,
		PriceChange:        ticker.PriceChange,
		PriceChangePercent: ticker.PriceChangePercent,
		Timestamp:          time.Now(),
	}

	result := r.GetDB(ctx).Save(entity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", ticker.Symbol).Msg("Failed to save ticker")
		return fmt.Errorf("failed to save ticker: %w", result.Error)
	}

	r.logger.Info().Str("symbol", ticker.Symbol).Str("exchange", ticker.Exchange).Msg("Ticker saved successfully")
	return nil
}

// GetTicker retrieves the latest ticker for a symbol from a specific exchange
func (r *MarketRepository) GetTicker(ctx context.Context, symbol, exchange string) (*model.Ticker, error) {
	var entity entity.MexcTickerEntity

	result := r.GetDB(ctx).
		Where("symbol = ? AND exchange = ?", symbol, exchange).
		Order("timestamp DESC").
		First(&entity)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			r.logger.Info().Str("symbol", symbol).Str("exchange", exchange).Msg("Ticker not found")
			return nil, apperror.ErrNotFound
		}
		r.logger.Error().Err(result.Error).Str("symbol", symbol).Msg("Failed to get ticker")
		return nil, fmt.Errorf("failed to get ticker: %w", result.Error)
	}

	return &model.Ticker{
		Symbol:             entity.Symbol,
		Exchange:           entity.Exchange,
		LastPrice:          entity.LastPrice,
		Volume:             entity.Volume,
		HighPrice:          entity.HighPrice,
		LowPrice:           entity.LowPrice,
		PriceChange:        entity.PriceChange,
		PriceChangePercent: entity.PriceChangePercent,
	}, nil
}

// GetAllTickers retrieves all latest tickers from a specific exchange
func (r *MarketRepository) GetAllTickers(ctx context.Context, exchange string) ([]*model.Ticker, error) {
	var entities []entity.MexcTickerEntity

	// Using a subquery to get the latest ticker for each symbol
	subQuery := r.GetDB(ctx).Model(&entity.MexcTickerEntity{}).
		Select("symbol, MAX(timestamp) as max_timestamp").
		Where("exchange = ?", exchange).
		Group("symbol")

	result := r.GetDB(ctx).
		Joins("JOIN (?) as sub ON mexc_tickers.symbol = sub.symbol AND mexc_tickers.timestamp = sub.max_timestamp", subQuery).
		Where("mexc_tickers.exchange = ?", exchange).
		Find(&entities)

	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("exchange", exchange).Msg("Failed to get all tickers")
		return nil, fmt.Errorf("failed to get all tickers: %w", result.Error)
	}

	tickers := make([]*model.Ticker, len(entities))
	for i, entity := range entities {
		tickers[i] = &model.Ticker{
			Symbol:             entity.Symbol,
			Exchange:           entity.Exchange,
			LastPrice:          entity.LastPrice,
			Volume:             entity.Volume,
			HighPrice:          entity.HighPrice,
			LowPrice:           entity.LowPrice,
			PriceChange:        entity.PriceChange,
			PriceChangePercent: entity.PriceChangePercent,
		}
	}

	return tickers, nil
}

// SaveKline stores a kline/candle in the database
func (r *MarketRepository) SaveKline(ctx context.Context, kline *model.Kline) error {
	// Generate a unique ID for the candle
	id := fmt.Sprintf("%s_%s_%s_%d", kline.Exchange, kline.Symbol, string(kline.Interval), kline.OpenTime.Unix())

	entity := &entity.MexcCandleEntity{
		ID:          id,
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

	// Try to find an existing candle with the same symbol, exchange, interval, and open time
	var existingID string
	result := r.GetDB(ctx).
		Table("mexc_candles").
		Select("id").
		Where("symbol = ? AND exchange = ? AND interval = ? AND open_time = ?",
			kline.Symbol, kline.Exchange, string(kline.Interval), kline.OpenTime).
		First(&existingID)

	// If the candle exists, update it; otherwise, create a new one
	if result.Error == nil {
		entity.ID = existingID
	}

	result = r.GetDB(ctx).Save(entity)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Str("symbol", kline.Symbol).Msg("Failed to save candle")
		return fmt.Errorf("failed to save candle: %w", result.Error)
	}

	r.logger.Info().Str("symbol", kline.Symbol).Str("interval", string(kline.Interval)).Msg("Candle saved successfully")
	return nil
}

// SaveKlines stores multiple klines/candles in the database
func (r *MarketRepository) SaveKlines(ctx context.Context, klines []*model.Kline) error {
	if len(klines) == 0 {
		return nil
	}

	// Use a transaction to save all candles
	tx := r.GetDB(ctx).Begin()
	if tx.Error != nil {
		r.logger.Error().Err(tx.Error).Msg("Failed to begin transaction")
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Create or update each candle
	for _, kline := range klines {
		// Generate a unique ID for the candle
		id := fmt.Sprintf("%s_%s_%s_%d", kline.Exchange, kline.Symbol, string(kline.Interval), kline.OpenTime.Unix())

		entity := entity.MexcCandleEntity{
			ID:          id,
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

		// Check if the candle already exists
		var existingID string
		result := tx.Table("mexc_candles").Select("id").Where(
			"symbol = ? AND exchange = ? AND interval = ? AND open_time = ?",
			kline.Symbol, kline.Exchange, string(kline.Interval), kline.OpenTime,
		).First(&existingID)

		if result.Error == nil {
			entity.ID = existingID
		}

		result = tx.Save(&entity)
		if result.Error != nil {
			tx.Rollback()
			r.logger.Error().Err(result.Error).Str("symbol", kline.Symbol).Msg("Failed to save candle in batch")
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
	var entity entity.MexcCandleEntity

	result := r.GetDB(ctx).
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
	}, nil
}

// GetKlines retrieves klines/candles for a symbol within a time range
func (r *MarketRepository) GetKlines(ctx context.Context, symbol, exchange string, interval model.KlineInterval, start, end time.Time, limit int) ([]*model.Kline, error) {
	var entities []entity.MexcCandleEntity

	query := r.GetDB(ctx).
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
		klines[i] = &model.Kline{
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

	return klines, nil
}

// GetLatestKline retrieves the most recent kline/candle for a symbol and interval
func (r *MarketRepository) GetLatestKline(ctx context.Context, symbol, exchange string, interval model.KlineInterval) (*model.Kline, error) {
	var entity entity.MexcCandleEntity

	result := r.GetDB(ctx).
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
	}, nil
}

// PurgeOldData removes market data older than the specified retention period
func (r *MarketRepository) PurgeOldData(ctx context.Context, olderThan time.Time) error {
	// Delete old ticker data
	if err := r.GetDB(ctx).Where("timestamp < ?", olderThan).Delete(&entity.MexcTickerEntity{}).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to purge old ticker data")
		return fmt.Errorf("failed to purge old ticker data: %w", err)
	}

	// Delete old candle data
	if err := r.GetDB(ctx).Where("open_time < ?", olderThan).Delete(&entity.MexcCandleEntity{}).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to purge old candle data")
		return fmt.Errorf("failed to purge old candle data: %w", err)
	}

	r.logger.Info().Time("olderThan", olderThan).Msg("Successfully purged old market data")
	return nil
}

// GetLatestTickers retrieves the latest tickers for all symbols
func (r *MarketRepository) GetLatestTickers(ctx context.Context, limit int) ([]*model.Ticker, error) {
	var entities []entity.MexcTickerEntity

	// Using a common table expression (CTE) to get the latest ticker for each symbol
	query := r.GetDB(ctx).
		Raw(`WITH latest_tickers AS (
			SELECT symbol, exchange, MAX(timestamp) as max_timestamp
			FROM mexc_tickers
			GROUP BY symbol, exchange
		)
		SELECT t.*
		FROM mexc_tickers t
		JOIN latest_tickers lt ON t.symbol = lt.symbol AND t.exchange = lt.exchange AND t.timestamp = lt.max_timestamp
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
		tickers[i] = &model.Ticker{
			Symbol:             entity.Symbol,
			Exchange:           entity.Exchange,
			LastPrice:          entity.LastPrice,
			Volume:             entity.Volume,
			HighPrice:          entity.HighPrice,
			LowPrice:           entity.LowPrice,
			PriceChange:        entity.PriceChange,
			PriceChangePercent: entity.PriceChangePercent,
		}
	}

	return tickers, nil
}

// Ensure MarketRepository implements port.MarketRepository
var _ port.MarketRepository = (*MarketRepository)(nil)

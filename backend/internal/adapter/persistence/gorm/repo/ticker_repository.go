package repo

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// TickerEntity represents a ticker in the database
type TickerEntity struct {
	ID                 string `gorm:"primaryKey;type:varchar(50)"`
	Symbol             string `gorm:"index;type:varchar(20)"`
	Exchange           string `gorm:"index;type:varchar(20)"`
	LastPrice          float64
	Volume             float64
	HighPrice          float64
	LowPrice           float64
	PriceChange        float64
	PriceChangePercent float64
	Timestamp          time.Time `gorm:"index"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
}

// KlineEntity represents a kline/candlestick in the database
type KlineEntity struct {
	ID        string    `gorm:"primaryKey;type:varchar(50)"`
	Symbol    string    `gorm:"index;type:varchar(20)"`
	Interval  string    `gorm:"index;type:varchar(10)"`
	OpenTime  time.Time `gorm:"index"`
	CloseTime time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// GormTickerRepository implements port.TickerRepository using GORM
type GormTickerRepository struct {
	BaseRepository
}

// NewGormTickerRepository creates a new GormTickerRepository
func NewGormTickerRepository(db *gorm.DB, logger *zerolog.Logger) *GormTickerRepository {
	return &GormTickerRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// Save saves a ticker to the database
func (r *GormTickerRepository) Save(ctx context.Context, ticker *model.Ticker) error {
	entity := r.toEntity(ticker)
	return r.Upsert(ctx, entity, []string{"symbol", "exchange"}, []string{
		"last_price", "volume", "high_price", "low_price",
		"price_change", "price_change_percent", "timestamp",
	})
}

// GetBySymbol retrieves a ticker by symbol
func (r *GormTickerRepository) GetBySymbol(ctx context.Context, symbol string) (*model.Ticker, error) {
	var entity TickerEntity
	err := r.FindOne(ctx, &entity, "symbol = ?", symbol)
	if err != nil {
		return nil, err
	}

	if entity.ID == "" {
		return nil, nil // Not found
	}

	return r.toDomain(&entity), nil
}

// GetAll retrieves all tickers
func (r *GormTickerRepository) GetAll(ctx context.Context) ([]*model.Ticker, error) {
	var entities []TickerEntity
	err := r.GetDB(ctx).
		Order("symbol ASC").
		Find(&entities).Error
	if err != nil {
		return nil, err
	}

	return r.toDomainSlice(entities), nil
}

// GetRecent retrieves the most recent tickers
func (r *GormTickerRepository) GetRecent(ctx context.Context, limit int) ([]*model.Ticker, error) {
	var entities []TickerEntity
	err := r.GetDB(ctx).
		Order("timestamp DESC").
		Limit(limit).
		Find(&entities).Error
	if err != nil {
		return nil, err
	}

	return r.toDomainSlice(entities), nil
}

// SaveKline saves a kline/candlestick to the database
func (r *GormTickerRepository) SaveKline(ctx context.Context, kline *model.Kline) error {
	entity := r.klineToEntity(kline)
	return r.Upsert(ctx, entity, []string{"symbol", "interval", "open_time"}, []string{
		"close_time", "open", "high", "low", "close", "volume",
	})
}

// GetKlines retrieves klines/candlesticks for a symbol and interval within a time range
func (r *GormTickerRepository) GetKlines(ctx context.Context, symbol string, interval model.KlineInterval, from, to time.Time, limit int) ([]*model.Kline, error) {
	var entities []KlineEntity

	query := r.GetDB(ctx).
		Where("symbol = ? AND interval = ?", symbol, string(interval))

	// Add time range conditions
	if !from.IsZero() {
		query = query.Where("open_time >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("close_time <= ?", to)
	}

	// Execute query
	err := query.
		Order("open_time ASC").
		Limit(limit).
		Find(&entities).Error
	if err != nil {
		return nil, err
	}

	return r.klinesToDomain(entities), nil
}

// Helper methods for entity conversion

// toEntity converts a domain ticker to a database entity
func (r *GormTickerRepository) toEntity(ticker *model.Ticker) *TickerEntity {
	if ticker == nil {
		return nil
	}

	return &TickerEntity{
		ID:                 ticker.Symbol + "_" + ticker.Exchange,
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
}

// toDomain converts a database entity to a domain ticker
func (r *GormTickerRepository) toDomain(entity *TickerEntity) *model.Ticker {
	if entity == nil {
		return nil
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
	}
}

// toDomainSlice converts a slice of database entities to domain tickers
func (r *GormTickerRepository) toDomainSlice(entities []TickerEntity) []*model.Ticker {
	tickers := make([]*model.Ticker, len(entities))
	for i, entity := range entities {
		tickers[i] = r.toDomain(&entity)
	}
	return tickers
}

// klineToEntity converts a domain kline to a database entity
func (r *GormTickerRepository) klineToEntity(kline *model.Kline) *KlineEntity {
	if kline == nil {
		return nil
	}

	return &KlineEntity{
		ID:        kline.Symbol + "_" + string(kline.Interval) + "_" + kline.OpenTime.Format(time.RFC3339),
		Symbol:    kline.Symbol,
		Interval:  string(kline.Interval),
		OpenTime:  kline.OpenTime,
		CloseTime: kline.CloseTime,
		Open:      kline.Open,
		High:      kline.High,
		Low:       kline.Low,
		Close:     kline.Close,
		Volume:    kline.Volume,
	}
}

// klineToDomain converts a database entity to a domain kline
func (r *GormTickerRepository) klineToDomain(entity *KlineEntity) *model.Kline {
	if entity == nil {
		return nil
	}

	return &model.Kline{
		Symbol:    entity.Symbol,
		Interval:  model.KlineInterval(entity.Interval),
		OpenTime:  entity.OpenTime,
		CloseTime: entity.CloseTime,
		Open:      entity.Open,
		High:      entity.High,
		Low:       entity.Low,
		Close:     entity.Close,
		Volume:    entity.Volume,
	}
}

// klinesToDomain converts a slice of database entities to domain klines
func (r *GormTickerRepository) klinesToDomain(entities []KlineEntity) []*model.Kline {
	klines := make([]*model.Kline, len(entities))
	for i, entity := range entities {
		klines[i] = r.klineToDomain(&entity)
	}
	return klines
}

// Ensure GormTickerRepository implements port.TickerRepository
var _ port.TickerRepository = (*GormTickerRepository)(nil)

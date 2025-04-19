package entity

import (
	"time"

	domainmodel "github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"gorm.io/gorm"
)

// MexcTickerEntity represents a ticker from MEXC exchange in the database
type MexcTickerEntity struct {
	ID            string `gorm:"primaryKey"`
	Symbol        string `gorm:"index"`
	Exchange      string `gorm:"index"`
	Price         float64
	Volume        float64
	High24h       float64
	Low24h        float64
	PriceChange   float64
	PercentChange float64
	LastUpdated   time.Time
}

// TableName specifies the table name for MexcTickerEntity
func (MexcTickerEntity) TableName() string {
	return "mexc_tickers"
}

// --- Mapping Functions: Entity <-> Domain Models ---

// MexcTickerEntity <-> model.Ticker
func (e *MexcTickerEntity) ToDomain() *domainmodel.Ticker {
	return &domainmodel.Ticker{
		ID:                 e.ID,
		Symbol:             e.Symbol,
		Exchange:           e.Exchange,
		LastPrice:          e.Price,
		Volume:             e.Volume,
		HighPrice:          e.High24h,
		LowPrice:           e.Low24h,
		PriceChange:        e.PriceChange,
		PriceChangePercent: e.PercentChange,
		Timestamp:          e.LastUpdated,
	}
}

func TickerEntityFromDomain(ticker *domainmodel.Ticker) *MexcTickerEntity {
	return &MexcTickerEntity{
		ID:            ticker.ID,
		Symbol:        ticker.Symbol,
		Exchange:      ticker.Exchange,
		Price:         ticker.LastPrice,
		Volume:        ticker.Volume,
		High24h:       ticker.HighPrice,
		Low24h:        ticker.LowPrice,
		PriceChange:   ticker.PriceChange,
		PercentChange: ticker.PriceChangePercent,
		LastUpdated:   ticker.Timestamp,
	}
}

// MexcCandleEntity <-> model.Kline
func (e *MexcCandleEntity) ToDomain() *domainmodel.Kline {
	return &domainmodel.Kline{
		Symbol:      e.Symbol,
		Exchange:    "MEXC",
		Interval:    domainmodel.KlineInterval(e.Interval),
		OpenTime:    e.OpenTime,
		CloseTime:   e.CloseTime,
		Open:        e.Open,
		High:        e.High,
		Low:         e.Low,
		Close:       e.Close,
		Volume:      e.Volume,
		QuoteVolume: e.QuoteVolume,
		TradeCount:  e.TradeCount,
		Complete:    true,
		IsClosed:    false,
	}
}

func CandleEntityFromDomain(c *domainmodel.Kline) *MexcCandleEntity {
	return &MexcCandleEntity{
		Symbol:      c.Symbol,
		Interval:    string(c.Interval),
		OpenTime:    c.OpenTime,
		CloseTime:   c.CloseTime,
		Open:        c.Open,
		High:        c.High,
		Low:         c.Low,
		Close:       c.Close,
		Volume:      c.Volume,
		QuoteVolume: c.QuoteVolume,
		TradeCount:  c.TradeCount,
	}
}

// MexcOrderBookEntity <-> model.OrderBook
func (e *MexcOrderBookEntity) ToDomain(entries []MexcOrderBookEntryEntity) *domainmodel.OrderBook {
	bids := []domainmodel.OrderBookEntry{}
	asks := []domainmodel.OrderBookEntry{}
	for _, entry := range entries {
		obEntry := domainmodel.OrderBookEntry{
			Price:    entry.Price,
			Quantity: entry.Quantity,
		}
		if entry.IsBid {
			bids = append(bids, obEntry)
		} else {
			asks = append(asks, obEntry)
		}
	}
	return &domainmodel.OrderBook{
		Symbol:       e.Symbol,
		Exchange:     "MEXC",
		LastUpdateID: e.LastUpdateID,
		SequenceNum:  0,
		Bids:         bids,
		Asks:         asks,
		Timestamp:    e.Timestamp,
	}
}

func OrderBookEntityFromDomain(ob *domainmodel.OrderBook) (*MexcOrderBookEntity, []MexcOrderBookEntryEntity) {
	entity := &MexcOrderBookEntity{
		Symbol:       ob.Symbol,
		LastUpdateID: ob.LastUpdateID,
		Timestamp:    ob.Timestamp,
	}
	entries := []MexcOrderBookEntryEntity{}
	for _, bid := range ob.Bids {
		entries = append(entries, MexcOrderBookEntryEntity{
			Price:    bid.Price,
			Quantity: bid.Quantity,
			IsBid:    true,
		})
	}
	for _, ask := range ob.Asks {
		entries = append(entries, MexcOrderBookEntryEntity{
			Price:    ask.Price,
			Quantity: ask.Quantity,
			IsBid:    false,
		})
	}
	return entity, entries
}

// MexcSymbolEntity <-> model.Symbol
func (e *MexcSymbolEntity) ToDomain() *domainmodel.Symbol {
	return &domainmodel.Symbol{
		Symbol:              e.Symbol,
		BaseAsset:           e.BaseAsset,
		QuoteAsset:          e.QuoteAsset,
		Exchange:            "MEXC",
		Status:              domainmodel.SymbolStatus(e.Status),
		MinPrice:            0,
		MaxPrice:            0,
		PricePrecision:      e.PricePrecision,
		MinQuantity:         e.MinQuantity,
		MaxQuantity:         e.MaxQuantity,
		QuantityPrecision:   e.QuantityPrecision,
		BaseAssetPrecision:  0,
		QuoteAssetPrecision: 0,
		MinNotional:         e.MinNotional,
		StepSize:            e.StepSize,
		TickSize:            e.TickSize,
		CreatedAt:           e.CreatedAt,
		UpdatedAt:           e.UpdatedAt,
	}
}

func SymbolEntityFromDomain(s *domainmodel.Symbol) *MexcSymbolEntity {
	return &MexcSymbolEntity{
		Symbol:            s.Symbol,
		BaseAsset:         s.BaseAsset,
		QuoteAsset:        s.QuoteAsset,
		Status:            string(s.Status),
		PricePrecision:    s.PricePrecision,
		QuantityPrecision: s.QuantityPrecision,
		MinNotional:       s.MinNotional,
		MinQuantity:       s.MinQuantity,
		MaxQuantity:       s.MaxQuantity,
		StepSize:          s.StepSize,
		TickSize:          s.TickSize,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}
}

// MexcCandleEntity represents candle (kline) data stored in the database
type MexcCandleEntity struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Symbol      string    `gorm:"not null;index:idx_mexc_candles_symbol_interval_time"`
	Interval    string    `gorm:"not null;index:idx_mexc_candles_symbol_interval_time,priority:2"`
	OpenTime    time.Time `gorm:"not null;index:idx_mexc_candles_symbol_interval_time,priority:3"`
	CloseTime   time.Time `gorm:"not null"`
	Open        float64   `gorm:"not null"`
	High        float64   `gorm:"not null"`
	Low         float64   `gorm:"not null"`
	Close       float64   `gorm:"not null"`
	Volume      float64   `gorm:"not null"`
	QuoteVolume float64   `gorm:"not null"`
	TradeCount  int64     `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"not null;autoUpdateTime"`
}

// TableName specifies the table name for MexcCandleEntity
func (MexcCandleEntity) TableName() string {
	return "mexc_candles"
}

// MexcOrderBookEntity represents order book data stored in the database
type MexcOrderBookEntity struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Symbol       string    `gorm:"not null;index:idx_mexc_orderbooks_symbol_time"`
	LastUpdateID int64     `gorm:"not null"`
	Timestamp    time.Time `gorm:"not null;index:idx_mexc_orderbooks_symbol_time,priority:2"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"not null;autoUpdateTime"`
}

// TableName specifies the table name for MexcOrderBookEntity
func (MexcOrderBookEntity) TableName() string {
	return "mexc_orderbooks"
}

// MexcOrderBookEntryEntity represents a single entry in the order book
type MexcOrderBookEntryEntity struct {
	ID          uint    `gorm:"primaryKey;autoIncrement"`
	OrderBookID uint    `gorm:"not null;index"`
	Price       float64 `gorm:"not null"`
	Quantity    float64 `gorm:"not null"`
	IsBid       bool    `gorm:"not null;index"` // true for bid, false for ask
}

// TableName specifies the table name for MexcOrderBookEntryEntity
func (MexcOrderBookEntryEntity) TableName() string {
	return "mexc_orderbook_entries"
}

// MexcSymbolEntity represents symbol information from MEXC
type MexcSymbolEntity struct {
	Symbol                 string  `gorm:"primaryKey"`
	BaseAsset              string  `gorm:"not null;index"`
	QuoteAsset             string  `gorm:"not null;index"`
	Status                 string  `gorm:"not null"` // e.g., "TRADING", "BREAK", etc.
	PricePrecision         int     `gorm:"not null"`
	QuantityPrecision      int     `gorm:"not null"`
	MinNotional            float64 `gorm:"not null"`
	MinQuantity            float64 `gorm:"not null"`
	MaxQuantity            float64 `gorm:"not null"`
	StepSize               float64 `gorm:"not null;default:0"`
	TickSize               float64 `gorm:"not null;default:0"`
	ListingDate            *time.Time
	TradingStartDate       *time.Time
	IsSpotTradingAllowed   bool      `gorm:"not null;default:true"`
	IsMarginTradingAllowed bool      `gorm:"not null;default:false"`
	CreatedAt              time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt              time.Time `gorm:"not null;autoUpdateTime"`
}

// TableName specifies the table name for MexcSymbolEntity
func (MexcSymbolEntity) TableName() string {
	return "mexc_symbols"
}

// MexcSyncStateEntity represents sync state information
type MexcSyncStateEntity struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement"`
	DataType           string    `gorm:"not null;uniqueIndex"` // "tickers", "candles", "orderbooks", "symbols"
	LastSyncTime       time.Time `gorm:"not null"`
	LastSuccessfulSync time.Time `gorm:"not null"`
	Status             string    `gorm:"not null;default:'idle'"` // "idle", "syncing", "failed"
	SyncInterval       int       `gorm:"not null"`                // in seconds
	ErrorMessage       string
	AdditionalInfo     string    // For storing info like which symbols/intervals were synced
	CreatedAt          time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt          time.Time `gorm:"not null;autoUpdateTime"`
}

// TableName specifies the table name for MexcSyncStateEntity
func (MexcSyncStateEntity) TableName() string {
	return "mexc_sync_states"
}

// BeforeCreate sets timestamps before creating a MexcSyncStateEntity
func (e *MexcSyncStateEntity) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now
	return nil
}

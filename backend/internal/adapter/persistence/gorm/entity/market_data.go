package entity

import (
	"time"

	domainMarket "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"gorm.io/gorm"
)

// MexcTickerEntity represents market ticker data stored in the database
type MexcTickerEntity struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement"`
	Symbol             string    `gorm:"not null;index:idx_mexc_tickers_symbol_time"`
	Price              float64   `gorm:"not null"`
	Volume             float64   `gorm:"not null"`
	QuoteVolume        float64   `gorm:"not null"`
	PriceChange        float64   `gorm:"not null"`
	PriceChangePercent float64   `gorm:"not null"`
	High               float64   `gorm:"not null"`
	Low                float64   `gorm:"not null"`
	OpenPrice          float64   `gorm:"not null"`
	ClosePrice         float64   `gorm:"not null"`
	Count              int64     `gorm:"not null"`
	Timestamp          time.Time `gorm:"not null;index:idx_mexc_tickers_symbol_time,priority:2"`
	IsFrozen           bool      `gorm:"not null;default:false"`
	CreatedAt          time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt          time.Time `gorm:"not null;autoUpdateTime"`
}

// TableName specifies the table name for MexcTickerEntity
func (MexcTickerEntity) TableName() string {
	return "mexc_tickers"
}

// --- Mapping Functions: Entity <-> Domain Models ---

// MexcTickerEntity <-> domainMarket.Ticker
func (e *MexcTickerEntity) ToDomain() *domainMarket.Ticker {
	return &domainMarket.Ticker{
		ID:            "", // Fill if you have an ID in domain
		Symbol:        e.Symbol,
		Price:         e.Price,
		Volume:        e.Volume,
		High24h:       e.High,
		Low24h:        e.Low,
		PriceChange:   e.PriceChange,
		PercentChange: e.PriceChangePercent,
		LastUpdated:   e.Timestamp,
		Exchange:      "MEXC", // Or map if you store exchange info
	}
}

func TickerEntityFromDomain(t *domainMarket.Ticker) *MexcTickerEntity {
	return &MexcTickerEntity{
		Symbol:             t.Symbol,
		Price:              t.Price,
		Volume:             t.Volume,
		QuoteVolume:        0, // Map if available
		PriceChange:        t.PriceChange,
		PriceChangePercent: t.PercentChange,
		High:               t.High24h,
		Low:                t.Low24h,
		OpenPrice:          0, // Map if available
		ClosePrice:         0, // Map if available
		Count:              0, // Map if available
		Timestamp:          t.LastUpdated,
		IsFrozen:           false, // Map if available
	}
}

// MexcCandleEntity <-> domainMarket.Candle
func (e *MexcCandleEntity) ToDomain() *domainMarket.Candle {
	return &domainMarket.Candle{
		Symbol:      e.Symbol,
		Exchange:    "MEXC",
		Interval:    domainMarket.Interval(e.Interval),
		OpenTime:    e.OpenTime,
		CloseTime:   e.CloseTime,
		Open:        e.Open,
		High:        e.High,
		Low:         e.Low,
		Close:       e.Close,
		Volume:      e.Volume,
		QuoteVolume: e.QuoteVolume,
		TradeCount:  e.TradeCount,
		Complete:    true, // Map if available
	}
}

func CandleEntityFromDomain(c *domainMarket.Candle) *MexcCandleEntity {
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

// MexcOrderBookEntity <-> domainMarket.OrderBook
func (e *MexcOrderBookEntity) ToDomain(entries []MexcOrderBookEntryEntity) *domainMarket.OrderBook {
	bids := []domainMarket.OrderBookEntry{}
	asks := []domainMarket.OrderBookEntry{}
	for _, entry := range entries {
		obEntry := domainMarket.OrderBookEntry{
			Price:    entry.Price,
			Quantity: entry.Quantity,
		}
		if entry.IsBid {
			bids = append(bids, obEntry)
		} else {
			asks = append(asks, obEntry)
		}
	}
	return &domainMarket.OrderBook{
		Symbol:       e.Symbol,
		LastUpdated:  e.Timestamp,
		Bids:         bids,
		Asks:         asks,
		Exchange:     "MEXC",
		SequenceNum:  0, // Map if available
		LastUpdateID: e.LastUpdateID,
	}
}

func OrderBookEntityFromDomain(ob *domainMarket.OrderBook) (*MexcOrderBookEntity, []MexcOrderBookEntryEntity) {
	entity := &MexcOrderBookEntity{
		Symbol:       ob.Symbol,
		LastUpdateID: ob.LastUpdateID,
		Timestamp:    ob.LastUpdated,
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

// MexcSymbolEntity <-> domainMarket.Symbol
func (e *MexcSymbolEntity) ToDomain() *domainMarket.Symbol {
	return &domainMarket.Symbol{
		Symbol:              e.Symbol,
		BaseAsset:           e.BaseAsset,
		QuoteAsset:          e.QuoteAsset,
		Exchange:            "MEXC",
		Status:              e.Status,
		MinPrice:            0, // Map if available
		MaxPrice:            0, // Map if available
		PricePrecision:      e.PricePrecision,
		MinQty:              e.MinQuantity,
		MaxQty:              e.MaxQuantity,
		QtyPrecision:        e.QuantityPrecision,
		BaseAssetPrecision:  0, // Map if available
		QuoteAssetPrecision: 0, // Map if available
		MinNotional:         e.MinNotional,
		MinLotSize:          0, // Map if available
		MaxLotSize:          0, // Map if available
	}
}

func SymbolEntityFromDomain(s *domainMarket.Symbol) *MexcSymbolEntity {
	return &MexcSymbolEntity{
		Symbol:            s.Symbol,
		BaseAsset:         s.BaseAsset,
		QuoteAsset:        s.QuoteAsset,
		Status:            s.Status,
		PricePrecision:    s.PricePrecision,
		QuantityPrecision: s.QtyPrecision,
		MinNotional:       s.MinNotional,
		MinQuantity:       s.MinQty,
		MaxQuantity:       s.MaxQty,
		StepSize:          s.StepSize,
		TickSize:          s.TickSize,
		// ListingDate:    s.ListingDate, // not present in domain model
		// TradingStartDate: s.TradingStartDate, // not present in domain model
		// IsSpotTradingAllowed: s.IsSpotTradingAllowed, // not present in domain model
		// IsMarginTradingAllowed: s.IsMarginTradingAllowed, // not present in domain model
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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

// MexcSyncStateEntity tracks the last successful sync with MEXC API
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

// BeforeCreate hook to initialize timestamps
func (e *MexcSyncStateEntity) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if e.LastSyncTime.IsZero() {
		e.LastSyncTime = now
	}
	if e.LastSuccessfulSync.IsZero() {
		e.LastSuccessfulSync = now
	}
	return nil
}

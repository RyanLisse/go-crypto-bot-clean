package compat

import (
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
)

// ConvertMarketTickerToTicker converts market.Ticker to model.Ticker
func ConvertMarketTickerToTicker(mt *market.Ticker) *model.Ticker {
	return &model.Ticker{
		ID:                 mt.ID,
		Symbol:             mt.Symbol,
		LastPrice:          mt.Price,
		Volume:             mt.Volume,
		HighPrice:          mt.High24h,
		LowPrice:           mt.Low24h,
		PriceChange:        mt.PriceChange,
		PriceChangePercent: mt.PercentChange,
		Timestamp:          mt.LastUpdated,
		Exchange:           mt.Exchange,
	}
}

// ConvertTickerToMarketTicker converts model.Ticker to market.Ticker
func ConvertTickerToMarketTicker(t *model.Ticker) *market.Ticker {
	return &market.Ticker{
		ID:            t.ID,
		Symbol:        t.Symbol,
		Price:         t.LastPrice,
		Volume:        t.Volume,
		High24h:       t.HighPrice,
		Low24h:        t.LowPrice,
		PriceChange:   t.PriceChange,
		PercentChange: t.PriceChangePercent,
		LastUpdated:   t.Timestamp,
		Exchange:      t.Exchange,
	}
}

// ConvertMarketSymbolToSymbol converts market.Symbol to model.Symbol
func ConvertMarketSymbolToSymbol(ms *market.Symbol) *model.Symbol {
	var status model.SymbolStatus = model.SymbolStatusHalt
	if ms.Status == string(model.SymbolStatusTrading) {
		status = model.SymbolStatusTrading
	} else if ms.Status == string(model.SymbolStatusBreak) {
		status = model.SymbolStatusBreak
	}

	return &model.Symbol{
		Symbol:              ms.Symbol,
		BaseAsset:           ms.BaseAsset,
		QuoteAsset:          ms.QuoteAsset,
		Exchange:            ms.Exchange,
		Status:              status,
		MinPrice:            ms.MinPrice,
		MaxPrice:            ms.MaxPrice,
		PricePrecision:      ms.PricePrecision,
		MinQuantity:         ms.MinQty,
		MaxQuantity:         ms.MaxQty,
		QuantityPrecision:   ms.QtyPrecision,
		BaseAssetPrecision:  ms.BaseAssetPrecision,
		QuoteAssetPrecision: ms.QuoteAssetPrecision,
		MinNotional:         ms.MinNotional,
		StepSize:            ms.StepSize,
		TickSize:            ms.TickSize,
		AllowedOrderTypes:   ms.AllowedOrderTypes,
		CreatedAt:           ms.CreatedAt,
		UpdatedAt:           ms.UpdatedAt,
	}
}

// ConvertSymbolToMarketSymbol converts model.Symbol to market.Symbol
func ConvertSymbolToMarketSymbol(s *model.Symbol) *market.Symbol {
	return &market.Symbol{
		Symbol:              s.Symbol,
		BaseAsset:           s.BaseAsset,
		QuoteAsset:          s.QuoteAsset,
		Exchange:            s.Exchange,
		Status:              string(s.Status),
		MinPrice:            s.MinPrice,
		MaxPrice:            s.MaxPrice,
		PricePrecision:      s.PricePrecision,
		MinQty:              s.MinQuantity,
		MaxQty:              s.MaxQuantity,
		QtyPrecision:        s.QuantityPrecision,
		BaseAssetPrecision:  s.BaseAssetPrecision,
		QuoteAssetPrecision: s.QuoteAssetPrecision,
		MinNotional:         s.MinNotional,
		StepSize:            s.StepSize,
		TickSize:            s.TickSize,
		AllowedOrderTypes:   s.AllowedOrderTypes,
		CreatedAt:           s.CreatedAt,
		UpdatedAt:           s.UpdatedAt,
	}
}

// NewSimpleTicker creates a new ticker with minimal information
func NewSimpleTicker(symbol string, price float64) *model.Ticker {
	return &model.Ticker{
		Symbol:    symbol,
		LastPrice: price,
		Timestamp: time.Now(),
	}
}

// NewMarketTicker creates a new market.Ticker with minimal information
func NewMarketTicker(symbol string, price float64) *market.Ticker {
	return &market.Ticker{
		Symbol:      symbol,
		Price:       price,
		LastUpdated: time.Now(),
	}
}

// ConvertMarketOrderBookToOrderBook converts market.OrderBook to model.OrderBook
func ConvertMarketOrderBookToOrderBook(mob *market.OrderBook) *model.OrderBook {
	// Convert bids
	bids := make([]model.OrderBookEntry, len(mob.Bids))
	for i, bid := range mob.Bids {
		bids[i] = model.OrderBookEntry{
			Price:    bid.Price,
			Quantity: bid.Quantity,
		}
	}

	// Convert asks
	asks := make([]model.OrderBookEntry, len(mob.Asks))
	for i, ask := range mob.Asks {
		asks[i] = model.OrderBookEntry{
			Price:    ask.Price,
			Quantity: ask.Quantity,
		}
	}

	return &model.OrderBook{
		Symbol:       mob.Symbol,
		LastUpdateID: mob.LastUpdateID,
		Bids:         bids,
		Asks:         asks,
		Timestamp:    mob.LastUpdated,
	}
}

// ConvertOrderBookToMarketOrderBook converts model.OrderBook to market.OrderBook
func ConvertOrderBookToMarketOrderBook(ob *model.OrderBook) *market.OrderBook {
	// Convert bids
	bids := make([]market.OrderBookEntry, len(ob.Bids))
	for i, bid := range ob.Bids {
		bids[i] = market.OrderBookEntry{
			Price:    bid.Price,
			Quantity: bid.Quantity,
		}
	}

	// Convert asks
	asks := make([]market.OrderBookEntry, len(ob.Asks))
	for i, ask := range ob.Asks {
		asks[i] = market.OrderBookEntry{
			Price:    ask.Price,
			Quantity: ask.Quantity,
		}
	}

	return &market.OrderBook{
		Symbol:       ob.Symbol,
		LastUpdateID: ob.LastUpdateID,
		Bids:         bids,
		Asks:         asks,
		LastUpdated:  ob.Timestamp,
	}
}

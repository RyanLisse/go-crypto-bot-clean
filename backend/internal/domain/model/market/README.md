# DEPRECATED: Market Models

The models in this directory are deprecated and will be removed in a future version.

## Migration Guide

Please use the canonical models in `internal/domain/model` instead:

| Deprecated Model | Canonical Model |
|------------------|----------------|
| `market.Ticker` | `model.Ticker` |
| `market.Symbol` | `model.Symbol` |
| `market.OrderBook` | `model.OrderBook` |
| `market.Candle` | `model.Kline` |

## Conversion Utilities

The `internal/domain/compat` package provides conversion functions to help with the migration:

```go
// Convert from market.Ticker to model.Ticker
ticker := compat.ConvertMarketTickerToTicker(marketTicker)

// Convert from model.Ticker to market.Ticker
marketTicker := compat.ConvertTickerToMarketTicker(ticker)

// Convert from market.Symbol to model.Symbol
symbol := compat.ConvertMarketSymbolToSymbol(marketSymbol)

// Convert from model.Symbol to market.Symbol
marketSymbol := compat.ConvertSymbolToMarketSymbol(symbol)

// Convert from market.OrderBook to model.OrderBook
orderBook := compat.ConvertMarketOrderBookToOrderBook(marketOrderBook)

// Convert from model.OrderBook to market.OrderBook
marketOrderBook := compat.ConvertOrderBookToMarketOrderBook(orderBook)
```

## Repository Interfaces

All repository interfaces have been updated to use the canonical models, with legacy methods provided for backward compatibility. These legacy methods will be removed in a future version.

## Timeline

- **Current**: Both models are supported, with the canonical models being the preferred option
- **Next Release**: Legacy models will be marked as deprecated
- **Future Release**: Legacy models will be removed entirely

Please update your code to use the canonical models as soon as possible.

package port

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
)

// OrderRepository defines the interface for order persistence operations
type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	GetByID(ctx context.Context, id string) (*model.Order, error)
	GetByClientOrderID(ctx context.Context, clientOrderID string) (*model.Order, error)
	Update(ctx context.Context, order *model.Order) error
	GetBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Order, error)
	GetByStatus(ctx context.Context, status model.OrderStatus, limit, offset int) ([]*model.Order, error)
	Count(ctx context.Context, filters map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) error
}

// WalletRepository defines the interface for wallet persistence operations
type WalletRepository interface {
	// Core wallet operations
	Save(ctx context.Context, wallet *model.Wallet) error
	GetByID(ctx context.Context, id string) (*model.Wallet, error)
	GetByUserID(ctx context.Context, userID string) (*model.Wallet, error)
	GetWalletsByUserID(ctx context.Context, userID string) ([]*model.Wallet, error)
	DeleteWallet(ctx context.Context, id string) error

	// Balance history operations
	SaveBalanceHistory(ctx context.Context, history *model.BalanceHistory) error
	GetBalanceHistory(ctx context.Context, userID string, asset model.Asset, from, to time.Time) ([]*model.BalanceHistory, error)
}

// NewCoinRepository defines the interface for new coin persistence operations
type NewCoinRepository interface {
	Save(ctx context.Context, newCoin *model.NewCoin) error
	GetByID(ctx context.Context, id string) (*model.NewCoin, error)
	GetBySymbol(ctx context.Context, symbol string) (*model.NewCoin, error)
	GetRecent(ctx context.Context, limit int) ([]*model.NewCoin, error)
	GetByStatus(ctx context.Context, status model.CoinStatus) ([]*model.NewCoin, error)
	Update(ctx context.Context, newCoin *model.NewCoin) error
	// FindRecentlyListed retrieves coins expected to list soon or recently became tradable.
	FindRecentlyListed(ctx context.Context, thresholdTime time.Time) ([]*model.NewCoin, error)
}

// TickerRepository defines the interface for ticker persistence operations
type TickerRepository interface {
	Save(ctx context.Context, ticker *model.Ticker) error
	GetBySymbol(ctx context.Context, symbol string) (*model.Ticker, error)
	GetAll(ctx context.Context) ([]*model.Ticker, error)
	GetRecent(ctx context.Context, limit int) ([]*model.Ticker, error)
	SaveKline(ctx context.Context, kline *model.Kline) error
	GetKlines(ctx context.Context, symbol string, interval model.KlineInterval, from, to time.Time, limit int) ([]*model.Kline, error)
}

// AIConversationRepository defines the interface for AI conversation persistence
type AIConversationRepository interface {
	SaveMessage(ctx context.Context, userID string, message map[string]interface{}) error
	GetConversation(ctx context.Context, userID string, limit int) ([]map[string]interface{}, error)
	ClearConversation(ctx context.Context, userID string) error
}

// StrategyRepository defines the interface for strategy persistence
type StrategyRepository interface {
	SaveConfig(ctx context.Context, strategyID string, config map[string]interface{}) error
	GetConfig(ctx context.Context, strategyID string) (map[string]interface{}, error)
	ListStrategies(ctx context.Context) ([]string, error)
	DeleteStrategy(ctx context.Context, strategyID string) error
}

// NotificationRepository defines the interface for notification persistence
type NotificationRepository interface {
	SavePreferences(ctx context.Context, userID string, preferences map[string]interface{}) error
	GetPreferences(ctx context.Context, userID string) (map[string]interface{}, error)
	SaveNotification(ctx context.Context, notification map[string]interface{}) error
	GetNotifications(ctx context.Context, userID string, limit, offset int) ([]map[string]interface{}, error)
}

// AnalyticsRepository defines the interface for analytics data persistence
type AnalyticsRepository interface {
	SaveMetrics(ctx context.Context, metrics map[string]interface{}) error
	GetMetrics(ctx context.Context, from, to time.Time) ([]map[string]interface{}, error)
	GetPerformanceByStrategy(ctx context.Context, strategyID string, from, to time.Time) (map[string]interface{}, error)
}

// MarketDataRepository defines the interface for market data persistence operations
type MarketDataRepository interface {
	// Ticker operations
	SaveTicker(ctx context.Context, ticker *model.Ticker) error
	GetTicker(ctx context.Context, exchange, symbol string) (*model.Ticker, error)
	GetAllTickers(ctx context.Context, exchange string) ([]*model.Ticker, error)

	// Kline operations
	SaveKline(ctx context.Context, kline *model.Kline) error
	GetKlines(ctx context.Context, exchange, symbol string, interval model.KlineInterval, from, to time.Time, limit int) ([]*model.Kline, error)

	// Order book operations
	SaveOrderBook(ctx context.Context, orderBook *model.OrderBook) error
	GetOrderBook(ctx context.Context, exchange, symbol string) (*model.OrderBook, error)

	// Symbol operations
	SaveSymbol(ctx context.Context, symbol *model.Symbol) error
	GetSymbol(ctx context.Context, exchange, symbol string) (*model.Symbol, error)
	GetAllSymbols(ctx context.Context, exchange string) ([]*model.Symbol, error)

	// Legacy methods for backward compatibility
	// These methods will be removed in a future version

	// Ticker operations (legacy)
	SaveTickerLegacy(ctx context.Context, ticker market.Ticker) error
	GetTickerLegacy(ctx context.Context, exchange, symbol string) (*market.Ticker, error)
	GetAllTickersLegacy(ctx context.Context, exchange string) ([]market.Ticker, error)

	// Candle operations (legacy)
	SaveCandleLegacy(ctx context.Context, candle market.Candle) error
	GetCandlesLegacy(ctx context.Context, exchange, symbol string, interval string, from, to time.Time, limit int) ([]market.Candle, error)

	// Order book operations (legacy)
	SaveOrderBookLegacy(ctx context.Context, orderBook market.OrderBook) error
	GetOrderBookLegacy(ctx context.Context, exchange, symbol string) (*market.OrderBook, error)

	// Symbol operations (legacy)
	SaveSymbolLegacy(ctx context.Context, symbol market.Symbol) error
	GetSymbolLegacy(ctx context.Context, exchange, symbol string) (*market.Symbol, error)
	GetAllSymbolsLegacy(ctx context.Context, exchange string) ([]market.Symbol, error)
}

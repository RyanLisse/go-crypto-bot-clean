package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// MEXCClient defines the interface for interacting with MEXC exchange API
type MEXCClient interface {
	// GetNewListings retrieves information about newly listed coins
	GetNewListings(ctx context.Context) ([]*model.NewCoin, error)

	// GetSymbolInfo retrieves detailed information about a trading symbol
	GetSymbolInfo(ctx context.Context, symbol string) (*model.SymbolInfo, error) // Changed return type

	// GetSymbolStatus checks if a symbol is currently tradeable
   GetSymbolStatus(ctx context.Context, symbol string) (model.CoinStatus, error)

	// GetTradingSchedule retrieves the listing and trading schedule for a symbol
	GetTradingSchedule(ctx context.Context, symbol string) (model.TradingSchedule, error)

	// GetSymbolConstraints retrieves trading constraints for a symbol
	GetSymbolConstraints(ctx context.Context, symbol string) (*model.SymbolConstraints, error)

	// GetExchangeInfo retrieves information about all symbols on the exchange
	GetExchangeInfo(ctx context.Context) (*model.ExchangeInfo, error) // Assuming model.ExchangeInfo exists or needs creation

	// GetMarketData retrieves ticker data (added based on market_data_service usage)
	GetMarketData(ctx context.Context, symbol string) (*model.Ticker, error)

	// GetKlines retrieves candle data (added based on market_data_service usage)
	GetKlines(ctx context.Context, symbol string, interval model.KlineInterval, limit int) ([]*model.Kline, error)

	// GetOrderBook retrieves order book data (added based on market_data_service usage)
	GetOrderBook(ctx context.Context, symbol string, depth int) (*model.OrderBook, error)

	// GetAccount retrieves account information from MEXC (added back from old MexcAPI)
	GetAccount(ctx context.Context) (*model.Wallet, error)
	// Trading Methods (merged from MexcAPI)
	PlaceOrder(ctx context.Context, symbol string, side model.OrderSide, orderType model.OrderType, quantity float64, price float64, timeInForce model.TimeInForce) (*model.Order, error)
	CancelOrder(ctx context.Context, symbol string, orderID string) error
	GetOrderStatus(ctx context.Context, symbol string, orderID string) (*model.Order, error)
	GetOpenOrders(ctx context.Context, symbol string) ([]*model.Order, error)
	GetOrderHistory(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error)
}

// Package huma provides OpenAPI documentation for the API.
package huma

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"

	"go-crypto-bot-clean/backend/internal/api/huma/auth"
	"go-crypto-bot-clean/backend/internal/api/huma/strategy"
	"go-crypto-bot-clean/backend/internal/api/huma/user"
	"go-crypto-bot-clean/backend/internal/api/service"
)

// Config represents the configuration for the Huma API documentation.
type Config struct {
	Title       string
	Description string
	Version     string
	BasePath    string
}

// DefaultConfig returns a default configuration for the Huma API documentation.
func DefaultConfig() Config {
	return Config{
		Title:       "Crypto Trading Bot API",
		Description: "API for the cryptocurrency trading bot",
		Version:     "1.0.0",
		BasePath:    "/api/v1",
	}
}

// SetupHuma sets up the Huma API documentation.
func SetupHuma(router chi.Router, config Config, services *service.Provider) huma.API {
	// Create a new Huma API
	api := humachi.New(router, huma.DefaultConfig(config.Title, config.Version))

	// Register endpoints
	registerBacktestEndpointsWithService(api, config.BasePath, services)
	strategy.RegisterStrategyEndpoints(api, config.BasePath, services.StrategyService)
	auth.RegisterAuthEndpoints(api, config.BasePath, services.AuthService)
	user.RegisterUserEndpoints(api, config.BasePath, services.UserService)
	// Register newly added endpoints (using basic http.HandlerFunc for now)
	// Portfolio
	huma.Register(api, huma.Operation{
		OperationID: "get-portfolio",
		Method:      http.MethodGet,
		Path:        config.BasePath + "/portfolio",
		Summary:     "Get user portfolio",
		Tags:        []string{"Portfolio"},
	}, PortfolioHandler) // Assumes auth middleware applied upstream
	huma.Register(api, huma.Operation{OperationID: "get-portfolio-performance", Method: http.MethodGet, Path: config.BasePath + "/portfolio/performance", Summary: "Get portfolio performance", Tags: []string{"Portfolio"}}, PortfolioPerformanceHandler)
	huma.Register(api, huma.Operation{OperationID: "get-portfolio-value", Method: http.MethodGet, Path: config.BasePath + "/portfolio/value", Summary: "Get portfolio total value", Tags: []string{"Portfolio"}}, PortfolioValueHandler)
	huma.Register(api, huma.Operation{OperationID: "get-portfolio-holdings-top", Method: http.MethodGet, Path: config.BasePath + "/portfolio/holdings/top", Summary: "Get top portfolio holdings", Tags: []string{"Portfolio"}}, TopHoldingsHandler)
	huma.Register(api, huma.Operation{OperationID: "validate-account-keys", Method: http.MethodPost, Path: config.BasePath + "/account/validate-keys", Summary: "Validate account API keys", Tags: []string{"Account"}}, ValidateKeysHandler)
	huma.Register(api, huma.Operation{OperationID: "get-status", Method: http.MethodGet, Path: config.BasePath + "/status", Summary: "Get API status", Tags: []string{"Status"}}, StatusHandler)
	huma.Register(api, huma.Operation{OperationID: "start-processes", Method: http.MethodPost, Path: config.BasePath + "/processes/start", Summary: "Start background processes", Tags: []string{"Processes"}}, StartProcessesHandler)
	huma.Register(api, huma.Operation{OperationID: "stop-processes", Method: http.MethodPost, Path: config.BasePath + "/processes/stop", Summary: "Stop background processes", Tags: []string{"Processes"}}, StopProcessesHandler)
	huma.Register(api, huma.Operation{OperationID: "get-upcoming-today-tomorrow", Method: http.MethodGet, Path: config.BasePath + "/newcoins/upcoming/today-and-tomorrow", Summary: "Get upcoming coins for today and tomorrow", Tags: []string{"New Coins"}}, GetUpcomingTodayTomorrowHandler)
	huma.Register(api, huma.Operation{OperationID: "get-newcoins-by-date", Method: http.MethodGet, Path: config.BasePath + "/newcoins/date/{date}", Summary: "Get new coins by date", Tags: []string{"New Coins"}}, GetNewCoinsByDateHandler)
	huma.Register(api, huma.Operation{OperationID: "get-newcoins-by-date-range", Method: http.MethodGet, Path: config.BasePath + "/newcoins/date-range", Summary: "Get new coins by date range", Tags: []string{"New Coins"}}, GetNewCoinsByDateRangeHandler)

	// Account
	huma.Register(api, huma.Operation{
		OperationID: "get-account-details",
		Method:      http.MethodGet,
		Path:        config.BasePath + "/account/details",
		Summary:     "Get user account details (wallet balances)",
		Tags:        []string{"Account"},
	}, AccountDetailsHandler) // Assumes auth middleware applied upstream

	// Trades
	huma.Register(api, huma.Operation{
		OperationID: "get-trades",
		Method:      http.MethodGet,
		Path:        config.BasePath + "/trades",
		Summary:     "Get recent trades",
		Tags:        []string{"Trades"},
	}, GetTradesHandler) // Assumes auth middleware applied upstream

	huma.Register(api, huma.Operation{
		OperationID: "create-trade",
		Method:      http.MethodPost,
		Path:        config.BasePath + "/trades",
		Summary:     "Execute a new trade",
		Tags:        []string{"Trades"},
	}, CreateTradeHandler) // Assumes auth middleware applied upstream

	huma.Register(api, huma.Operation{
		OperationID: "get-trade-status",
		Method:      http.MethodGet,
		Path:        config.BasePath + "/trades/{tradeId}",
		Summary:     "Get status of a specific trade",
		Tags:        []string{"Trades"},
	}, GetTradeStatusHandler) // Assumes auth middleware applied upstream

	// Config
	huma.Register(api, huma.Operation{
		OperationID: "get-config",
		Method:      http.MethodGet,
		Path:        config.BasePath + "/config",
		Summary:     "Get current bot configuration",
		Tags:        []string{"Configuration"},
	}, GetConfigHandler) // Assumes auth middleware applied upstream

	huma.Register(api, huma.Operation{
		OperationID: "update-config",
		Method:      http.MethodPut,
		Path:        config.BasePath + "/config",
		Summary:     "Update bot configuration",
		Tags:        []string{"Configuration"},
	}, UpdateConfigHandler) // Assumes auth middleware applied upstream

	huma.Register(api, huma.Operation{
		OperationID: "get-default-config",
		Method:      http.MethodGet,
		Path:        config.BasePath + "/config/default",
		Summary:     "Get default bot configuration",
		Tags:        []string{"Configuration"},
	}, GetDefaultConfigHandler) // Public?

	// Analytics
	huma.Register(api, huma.Operation{
		OperationID: "get-analytics",
		Method:      http.MethodGet,
		Path:        config.BasePath + "/analytics",
		Summary:     "Get general trading analytics",
		Tags:        []string{"Analytics"},
	}, GetAnalyticsHandler) // Assumes auth middleware applied upstream

	huma.Register(api, huma.Operation{
		OperationID: "get-win-rate",
		Method:      http.MethodGet,
		Path:        config.BasePath + "/analytics/win-rate",
		Summary:     "Get trading win rate",
		Tags:        []string{"Analytics"},
	}, GetWinRateHandler) // Assumes auth middleware applied upstream

	huma.Register(api, huma.Operation{
		OperationID: "get-balance-history",
		Method:      http.MethodGet,
		Path:        config.BasePath + "/analytics/balance-history",
		Summary:     "Get historical account balance",
		Tags:        []string{"Analytics"},
	}, GetBalanceHistoryHandler) // Assumes auth middleware applied upstream

	// New Coins
	huma.Register(api, huma.Operation{
		OperationID: "get-new-coins",
		Method:      http.MethodGet,
		Path:        config.BasePath + "/newcoins",
		Summary:     "Get recently listed coins",
		Tags:        []string{"New Coins"},
	}, GetNewCoinsHandler) // Public?

	huma.Register(api, huma.Operation{
		OperationID: "get-upcoming-coins",
		Method:      http.MethodGet,
		Path:        config.BasePath + "/newcoins/upcoming",
		Summary:     "Get upcoming coin listings",
		Tags:        []string{"New Coins"},
	}, GetUpcomingCoinsHandler) // Public?

	huma.Register(api, huma.Operation{
		OperationID: "process-new-coins",
		Method:      http.MethodPost,
		Path:        config.BasePath + "/newcoins/process",
		Summary:     "Trigger processing of new coin listings",
		Tags:        []string{"New Coins"},
	}, ProcessNewCoinsHandler) // Assumes auth middleware applied upstream

	return api
}

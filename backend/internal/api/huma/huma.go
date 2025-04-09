// Package huma provides OpenAPI documentation for the API.
package huma

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor" // Enable CBOR support
	"github.com/go-chi/chi/v5"
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

// PortfolioSummaryResponse represents the overall portfolio status
type PortfolioSummaryResponse struct {
	Body struct {
		TotalValue       float64             `json:"total_value" doc:"Total value of the portfolio in USDT" example:"1250.75"`
		ActiveTradeCount int                 `json:"active_trade_count" doc:"Number of active trades" example:"3"`
		ActiveTrades     []TradeResponse     `json:"active_trades" doc:"List of active trades"`
		Performance      PerformanceResponse `json:"performance" doc:"Performance metrics"`
		Timestamp        time.Time           `json:"timestamp" doc:"Timestamp of the response" example:"2023-01-01T12:00:00Z"`
	}
}

// TradeResponse represents a single trading position
type TradeResponse struct {
	ID               uint                      `json:"id" doc:"Unique identifier for the trade" example:"123"`
	Symbol           string                    `json:"symbol" doc:"Trading pair symbol" example:"BTC/USDT"`
	PurchasePrice    float64                   `json:"purchase_price" doc:"Price at which the coin was purchased" example:"50000"`
	CurrentPrice     float64                   `json:"current_price" doc:"Current price of the coin" example:"52000"`
	Quantity         float64                   `json:"quantity" doc:"Quantity of the coin" example:"0.02"`
	PurchaseTime     time.Time                 `json:"purchase_time" doc:"Time when the coin was purchased" example:"2023-01-01T10:00:00Z"`
	ProfitPercent    float64                   `json:"profit_percent" doc:"Current profit percentage" example:"4.0"`
	CurrentValue     float64                   `json:"current_value" doc:"Current value of the position in USDT" example:"1040"`
	StopLossPrice    float64                   `json:"stop_loss_price" doc:"Stop loss price" example:"48000"`
	TakeProfitLevels []TakeProfitLevelResponse `json:"take_profit_levels" doc:"Take profit levels"`
}

// TakeProfitLevelResponse represents a take profit level
type TakeProfitLevelResponse struct {
	Price    float64 `json:"price" doc:"Take profit price" example:"55000"`
	Percent  float64 `json:"percent" doc:"Percentage of position to sell at this level" example:"0.25"`
	Executed bool    `json:"executed" doc:"Whether this level has been executed" example:"false"`
}

// PerformanceResponse represents trading performance metrics
type PerformanceResponse struct {
	TotalTrades           int     `json:"total_trades" doc:"Total number of trades" example:"50"`
	WinningTrades         int     `json:"winning_trades" doc:"Number of winning trades" example:"30"`
	LosingTrades          int     `json:"losing_trades" doc:"Number of losing trades" example:"20"`
	WinRate               float64 `json:"win_rate" doc:"Win rate percentage" example:"60.0"`
	TotalProfitLoss       float64 `json:"total_profit_loss" doc:"Total profit/loss in USDT" example:"500.25"`
	AverageProfitPerTrade float64 `json:"average_profit_per_trade" doc:"Average profit per trade in USDT" example:"10.0"`
	LargestProfit         float64 `json:"largest_profit" doc:"Largest profit in USDT" example:"100.0"`
	LargestLoss           float64 `json:"largest_loss" doc:"Largest loss in USDT" example:"-50.0"`
	TimeRange             string  `json:"time_range" doc:"Time range for the metrics" example:"30d"`
}

// TradeRequest represents a request to execute a trade
type TradeRequest struct {
	Body struct {
		Symbol string  `json:"symbol" doc:"Trading pair symbol" example:"BTC/USDT" maxLength:"20" binding:"required"`
		Amount float64 `json:"amount" doc:"Amount to trade in USDT (optional)" example:"100.0" minimum:"0"`
	}
}

// TradeExecutionResponse represents the result of a trade execution
type TradeExecutionResponse struct {
	Body struct {
		ID            uint      `json:"id" doc:"Unique identifier for the trade" example:"123"`
		Symbol        string    `json:"symbol" doc:"Trading pair symbol" example:"BTC/USDT"`
		Price         float64   `json:"price" doc:"Execution price" example:"50000"`
		Quantity      float64   `json:"quantity" doc:"Quantity of the coin" example:"0.002"`
		Total         float64   `json:"total" doc:"Total cost in USDT" example:"100.0"`
		ExecutionTime time.Time `json:"execution_time" doc:"Time of execution" example:"2023-01-01T12:00:00Z"`
		Status        string    `json:"status" doc:"Status of the trade" example:"completed" enum:"pending,completed,failed"`
	}
}

// SellRequest represents a request to sell a coin
type SellRequest struct {
	Body struct {
		CoinID uint    `json:"coin_id" doc:"ID of the coin to sell" example:"123" binding:"required"`
		Amount float64 `json:"amount" doc:"Amount to sell (optional)" example:"0.001" minimum:"0"`
		All    bool    `json:"all" doc:"Whether to sell all" example:"false"`
	}
}

// NewCoinResponse represents a newly detected coin
type NewCoinResponse struct {
	Body struct {
		ID          uint      `json:"id" doc:"Unique identifier for the coin" example:"123"`
		Symbol      string    `json:"symbol" doc:"Trading pair symbol" example:"NEW/USDT"`
		Name        string    `json:"name,omitempty" doc:"Name of the coin" example:"New Coin"`
		FoundAt     time.Time `json:"found_at" doc:"Time when the coin was detected" example:"2023-01-01T12:00:00Z"`
		QuoteVolume float64   `json:"quote_volume" doc:"Quote volume in USDT" example:"1000000"`
		IsProcessed bool      `json:"is_processed" doc:"Whether the coin has been processed" example:"false"`
	}
}

// NewCoinsListResponse represents a list of newly detected coins
type NewCoinsListResponse struct {
	Body struct {
		Coins     []NewCoinResponse `json:"coins" doc:"List of newly detected coins"`
		Count     int               `json:"count" doc:"Number of coins" example:"5"`
		Timestamp time.Time         `json:"timestamp" doc:"Timestamp of the response" example:"2023-01-01T12:00:00Z"`
	}
}

// ProcessNewCoinsRequest represents a request to process new coins
type ProcessNewCoinsRequest struct {
	Body struct {
		CoinIDs []uint `json:"coin_ids" doc:"IDs of coins to process" example:"[123, 456]" binding:"required"`
	}
}

// ConfigResponse represents the bot configuration
type ConfigResponse struct {
	Body struct {
		USDTPerTrade     float64   `json:"usdt_per_trade" doc:"Amount of USDT per trade" example:"20.0"`
		StopLossPercent  float64   `json:"stop_loss_percent" doc:"Stop loss percentage" example:"10.0"`
		TakeProfitLevels []float64 `json:"take_profit_levels" doc:"Take profit levels" example:"[5.0, 10.0, 15.0, 20.0]"`
		SellPercentages  []float64 `json:"sell_percentages" doc:"Sell percentages for each take profit level" example:"[0.25, 0.25, 0.25, 0.25]"`
		UpdatedAt        time.Time `json:"updated_at" doc:"Last update time" example:"2023-01-01T12:00:00Z"`
	}
}

// ConfigUpdateRequest represents a request to update bot configuration
type ConfigUpdateRequest struct {
	Body struct {
		USDTPerTrade     *float64  `json:"usdt_per_trade,omitempty" doc:"Amount of USDT per trade" example:"20.0" minimum:"0"`
		StopLossPercent  *float64  `json:"stop_loss_percent,omitempty" doc:"Stop loss percentage" example:"10.0" minimum:"0" maximum:"100"`
		TakeProfitLevels []float64 `json:"take_profit_levels,omitempty" doc:"Take profit levels" example:"[5.0, 10.0, 15.0, 20.0]"`
		SellPercentages  []float64 `json:"sell_percentages,omitempty" doc:"Sell percentages for each take profit level" example:"[0.25, 0.25, 0.25, 0.25]"`
	}
}

// StatusResponse represents the system status
type StatusResponse struct {
	Body struct {
		Status           string    `json:"status" doc:"Overall system status" example:"running" enum:"running,stopped,error"`
		UptimeSeconds    int       `json:"uptime_seconds" doc:"Uptime in seconds" example:"3600"`
		ActiveProcesses  []string  `json:"active_processes" doc:"List of active processes" example:"[\"trade_watcher\", \"market_data\"]"`
		LastError        string    `json:"last_error,omitempty" doc:"Last error message" example:"Connection timeout"`
		LastErrorTime    time.Time `json:"last_error_time,omitempty" doc:"Time of the last error" example:"2023-01-01T12:00:00Z"`
		MemoryUsageMB    float64   `json:"memory_usage_mb" doc:"Memory usage in MB" example:"256.5"`
		CPUUsagePercent  float64   `json:"cpu_usage_percent" doc:"CPU usage percentage" example:"15.2"`
		DiskUsagePercent float64   `json:"disk_usage_percent" doc:"Disk usage percentage" example:"45.7"`
		Timestamp        time.Time `json:"timestamp" doc:"Timestamp of the status" example:"2023-01-01T12:00:00Z"`
	}
}

// ErrorResponse represents a standardized API error response
type ErrorResponse struct {
	Body struct {
		Code    string `json:"code" doc:"Error code identifier" example:"internal_error"`
		Message string `json:"message" doc:"Human-readable error message" example:"An internal error occurred"`
		Details string `json:"details,omitempty" doc:"Optional additional details" example:"Database connection failed"`
	}
}

// SetupHuma initializes the Huma API documentation.
func SetupHuma(router chi.Router, config Config) huma.API {
	// Create a new Huma API
	humaConfig := huma.DefaultConfig(config.Title, config.Version)
	humaConfig.Info.Description = config.Description

	api := humachi.New(router, humaConfig)

	// Register health check endpoint
	huma.Register(api, huma.Operation{
		OperationID: "health-check",
		Method:      http.MethodGet,
		Path:        "/health",
		Summary:     "Health check",
		Description: "Returns the health status of the API",
		Tags:        []string{"System"},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body struct {
			Status string `json:"status" doc:"Health status" example:"ok"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				Status string `json:"status" doc:"Health status" example:"ok"`
			}
		}{}
		resp.Body.Status = "ok"
		return resp, nil
	})

	// Register portfolio endpoints
	registerPortfolioEndpoints(api, config.BasePath)

	// Register trade endpoints
	registerTradeEndpoints(api, config.BasePath)

	// Register newcoin endpoints
	registerNewCoinEndpoints(api, config.BasePath)

	// Register config endpoints
	registerConfigEndpoints(api, config.BasePath)

	// Register status endpoints
	registerStatusEndpoints(api, config.BasePath)

	return api
}

// registerPortfolioEndpoints registers the portfolio endpoints.
func registerPortfolioEndpoints(api huma.API, basePath string) {
	// GET /portfolio
	huma.Register(api, huma.Operation{
		OperationID: "get-portfolio-summary",
		Method:      http.MethodGet,
		Path:        basePath + "/portfolio",
		Summary:     "Get portfolio summary",
		Description: "Returns a summary of the current portfolio including total value and active trades",
		Tags:        []string{"Portfolio"},
	}, func(ctx context.Context, input *struct{}) (*PortfolioSummaryResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})

	// GET /portfolio/active
	huma.Register(api, huma.Operation{
		OperationID: "get-active-trades",
		Method:      http.MethodGet,
		Path:        basePath + "/portfolio/active",
		Summary:     "Get active trades",
		Description: "Returns a list of active trades",
		Tags:        []string{"Portfolio"},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body struct {
			Trades    []TradeResponse `json:"trades" doc:"List of active trades"`
			Count     int             `json:"count" doc:"Number of trades" example:"3"`
			Timestamp time.Time       `json:"timestamp" doc:"Timestamp of the response" example:"2023-01-01T12:00:00Z"`
		}
	}, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})

	// GET /portfolio/performance
	huma.Register(api, huma.Operation{
		OperationID: "get-performance-metrics",
		Method:      http.MethodGet,
		Path:        basePath + "/portfolio/performance",
		Summary:     "Get performance metrics",
		Description: "Returns performance metrics for the portfolio",
		Tags:        []string{"Portfolio"},
	}, func(ctx context.Context, input *struct {
		TimeRange string `query:"time_range" doc:"Time range for metrics" example:"30d" enum:"7d,30d,90d,1y,all"`
	}) (*struct {
		Body PerformanceResponse
	}, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})

	// GET /portfolio/value
	huma.Register(api, huma.Operation{
		OperationID: "get-total-value",
		Method:      http.MethodGet,
		Path:        basePath + "/portfolio/value",
		Summary:     "Get total portfolio value",
		Description: "Returns the total value of the portfolio in USDT",
		Tags:        []string{"Portfolio"},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body struct {
			Value     float64   `json:"value" doc:"Total value in USDT" example:"1250.75"`
			Timestamp time.Time `json:"timestamp" doc:"Timestamp of the response" example:"2023-01-01T12:00:00Z"`
		}
	}, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})
}

// registerTradeEndpoints registers the trade endpoints.
func registerTradeEndpoints(api huma.API, basePath string) {
	// GET /trade/history
	huma.Register(api, huma.Operation{
		OperationID: "get-trade-history",
		Method:      http.MethodGet,
		Path:        basePath + "/trade/history",
		Summary:     "Get trade history",
		Description: "Returns the history of trades",
		Tags:        []string{"Trading"},
	}, func(ctx context.Context, input *struct {
		Limit  int    `query:"limit" doc:"Maximum number of trades to return" example:"10" default:"10" minimum:"1" maximum:"100"`
		Offset int    `query:"offset" doc:"Offset for pagination" example:"0" default:"0" minimum:"0"`
		Symbol string `query:"symbol" doc:"Filter by symbol" example:"BTC/USDT"`
	}) (*struct {
		Body struct {
			Trades []struct {
				ID            uint      `json:"id" doc:"Unique identifier for the trade" example:"123"`
				Symbol        string    `json:"symbol" doc:"Trading pair symbol" example:"BTC/USDT"`
				BuyPrice      float64   `json:"buy_price" doc:"Buy price" example:"50000"`
				SellPrice     float64   `json:"sell_price" doc:"Sell price" example:"52000"`
				Quantity      float64   `json:"quantity" doc:"Quantity of the coin" example:"0.02"`
				BuyTime       time.Time `json:"buy_time" doc:"Time of purchase" example:"2023-01-01T10:00:00Z"`
				SellTime      time.Time `json:"sell_time" doc:"Time of sale" example:"2023-01-01T12:00:00Z"`
				ProfitLoss    float64   `json:"profit_loss" doc:"Profit/loss in USDT" example:"40.0"`
				ProfitPercent float64   `json:"profit_percent" doc:"Profit/loss percentage" example:"4.0"`
			} `json:"trades" doc:"List of historical trades"`
			Count     int       `json:"count" doc:"Number of trades" example:"10"`
			Timestamp time.Time `json:"timestamp" doc:"Timestamp of the response" example:"2023-01-01T12:00:00Z"`
		}
	}, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})

	// POST /trade/buy
	huma.Register(api, huma.Operation{
		OperationID: "execute-trade",
		Method:      http.MethodPost,
		Path:        basePath + "/trade/buy",
		Summary:     "Execute a trade",
		Description: "Executes a buy trade for the specified symbol",
		Tags:        []string{"Trading"},
	}, func(ctx context.Context, input *TradeRequest) (*TradeExecutionResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})

	// POST /trade/sell
	huma.Register(api, huma.Operation{
		OperationID: "sell-coin",
		Method:      http.MethodPost,
		Path:        basePath + "/trade/sell",
		Summary:     "Sell a coin",
		Description: "Sells a previously bought coin",
		Tags:        []string{"Trading"},
	}, func(ctx context.Context, input *SellRequest) (*TradeExecutionResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})

	// GET /trade/status/{id}
	huma.Register(api, huma.Operation{
		OperationID: "get-trade-status",
		Method:      http.MethodGet,
		Path:        basePath + "/trade/status/{id}",
		Summary:     "Get trade status",
		Description: "Returns the status of a specific trade",
		Tags:        []string{"Trading"},
	}, func(ctx context.Context, input *struct {
		ID uint `path:"id" doc:"Trade ID" example:"123"`
	}) (*struct {
		Body struct {
			ID               uint                      `json:"id" doc:"Unique identifier for the trade" example:"123"`
			Symbol           string                    `json:"symbol" doc:"Trading pair symbol" example:"BTC/USDT"`
			Status           string                    `json:"status" doc:"Status of the trade" example:"active" enum:"active,completed,cancelled"`
			PurchasePrice    float64                   `json:"purchase_price" doc:"Price at which the coin was purchased" example:"50000"`
			CurrentPrice     float64                   `json:"current_price" doc:"Current price of the coin" example:"52000"`
			Quantity         float64                   `json:"quantity" doc:"Quantity of the coin" example:"0.02"`
			PurchaseTime     time.Time                 `json:"purchase_time" doc:"Time when the coin was purchased" example:"2023-01-01T10:00:00Z"`
			ProfitPercent    float64                   `json:"profit_percent" doc:"Current profit percentage" example:"4.0"`
			CurrentValue     float64                   `json:"current_value" doc:"Current value of the position in USDT" example:"1040"`
			StopLossPrice    float64                   `json:"stop_loss_price" doc:"Stop loss price" example:"48000"`
			TakeProfitLevels []TakeProfitLevelResponse `json:"take_profit_levels" doc:"Take profit levels"`
		}
	}, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})
}

// registerNewCoinEndpoints registers the newcoin endpoints.
func registerNewCoinEndpoints(api huma.API, basePath string) {
	// GET /newcoins
	huma.Register(api, huma.Operation{
		OperationID: "get-detected-coins",
		Method:      http.MethodGet,
		Path:        basePath + "/newcoins",
		Summary:     "Get detected coins",
		Description: "Returns a list of newly detected coins",
		Tags:        []string{"New Coins"},
	}, func(ctx context.Context, input *struct {
		Processed bool `query:"processed" doc:"Filter by processed status" example:"false"`
	}) (*NewCoinsListResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})

	// POST /newcoins/process
	huma.Register(api, huma.Operation{
		OperationID: "process-new-coins",
		Method:      http.MethodPost,
		Path:        basePath + "/newcoins/process",
		Summary:     "Process new coins",
		Description: "Processes newly detected coins",
		Tags:        []string{"New Coins"},
	}, func(ctx context.Context, input *ProcessNewCoinsRequest) (*struct {
		Body struct {
			ProcessedCoins []NewCoinResponse `json:"processed_coins" doc:"List of processed coins"`
			Count          int               `json:"count" doc:"Number of processed coins" example:"2"`
			Timestamp      time.Time         `json:"timestamp" doc:"Timestamp of the response" example:"2023-01-01T12:00:00Z"`
		}
	}, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})

	// POST /newcoins/detect
	huma.Register(api, huma.Operation{
		OperationID: "detect-new-coins",
		Method:      http.MethodPost,
		Path:        basePath + "/newcoins/detect",
		Summary:     "Detect new coins",
		Description: "Triggers detection of new coins",
		Tags:        []string{"New Coins"},
	}, func(ctx context.Context, input *struct{}) (*NewCoinsListResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})
}

// registerConfigEndpoints registers the config endpoints.
func registerConfigEndpoints(api huma.API, basePath string) {
	// GET /config
	huma.Register(api, huma.Operation{
		OperationID: "get-current-config",
		Method:      http.MethodGet,
		Path:        basePath + "/config",
		Summary:     "Get current configuration",
		Description: "Returns the current bot configuration",
		Tags:        []string{"Configuration"},
	}, func(ctx context.Context, input *struct{}) (*ConfigResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})

	// PUT /config
	huma.Register(api, huma.Operation{
		OperationID: "update-config",
		Method:      http.MethodPut,
		Path:        basePath + "/config",
		Summary:     "Update configuration",
		Description: "Updates the bot configuration",
		Tags:        []string{"Configuration"},
	}, func(ctx context.Context, input *ConfigUpdateRequest) (*ConfigResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})

	// GET /config/defaults
	huma.Register(api, huma.Operation{
		OperationID: "get-default-config",
		Method:      http.MethodGet,
		Path:        basePath + "/config/defaults",
		Summary:     "Get default configuration",
		Description: "Returns the default bot configuration",
		Tags:        []string{"Configuration"},
	}, func(ctx context.Context, input *struct{}) (*ConfigResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})
}

// registerStatusEndpoints registers the status endpoints.
func registerStatusEndpoints(api huma.API, basePath string) {
	// GET /status
	huma.Register(api, huma.Operation{
		OperationID: "get-status",
		Method:      http.MethodGet,
		Path:        basePath + "/status",
		Summary:     "Get system status",
		Description: "Returns the current system status",
		Tags:        []string{"System"},
	}, func(ctx context.Context, input *struct{}) (*StatusResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})

	// POST /status/start
	huma.Register(api, huma.Operation{
		OperationID: "start-processes",
		Method:      http.MethodPost,
		Path:        basePath + "/status/start",
		Summary:     "Start system processes",
		Description: "Starts all system processes",
		Tags:        []string{"System"},
	}, func(ctx context.Context, input *struct{}) (*StatusResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})

	// POST /status/stop
	huma.Register(api, huma.Operation{
		OperationID: "stop-processes",
		Method:      http.MethodPost,
		Path:        basePath + "/status/stop",
		Summary:     "Stop system processes",
		Description: "Stops all system processes",
		Tags:        []string{"System"},
	}, func(ctx context.Context, input *struct{}) (*StatusResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})
}

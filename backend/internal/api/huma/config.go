package huma

import (
	"context"
)

// --- GetConfig ---

// GetConfigInput defines input for getting config (empty).
type GetConfigInput struct{}

// GetConfigResponse defines the output structure for the bot's configuration.
type GetConfigResponse struct {
	Body ConfigResponseBody
}

// --- UpdateConfig ---

// UpdateConfigInput defines the input structure for updating the config.
type UpdateConfigInput struct {
	Body ConfigResponseBody // Expecting the full config structure for update
}

// UpdateConfigResponse defines the output structure after updating the config.
type UpdateConfigResponse struct {
	Body ConfigResponseBody
}

// --- GetDefaultConfig ---

// GetDefaultConfigInput defines input for getting default config (empty).
type GetDefaultConfigInput struct{}

// GetDefaultConfigResponse defines the output structure for the default config.
type GetDefaultConfigResponse struct {
	Body ConfigResponseBody
}

// --- Common Config Body ---

// ConfigResponseBody defines the structure for the bot's configuration used in responses/requests.
type ConfigResponseBody struct {
	Strategy            string   `json:"strategy"`
	RiskLevel           float64  `json:"risk_level"`
	MaxConcurrentTrades int      `json:"max_concurrent_trades"`
	TakeProfitPercent   float64  `json:"take_profit_percent"`
	StopLossPercent     float64  `json:"stop_loss_percent"`
	DailyTradeLimit     int      `json:"daily_trade_limit"`
	TradingPairs        []string `json:"trading_pairs"`
	TradingSchedule     struct {
		Days      []string `json:"days"`
		StartTime string   `json:"start_time"`
		EndTime   string   `json:"end_time"`
	} `json:"trading_schedule"`
}

// GetConfigHandler handles GET requests to /api/v1/config using Huma signature.
func GetConfigHandler(ctx context.Context, input *GetConfigInput) (*GetConfigResponse, error) {
	// Mock data.
	respBody := ConfigResponseBody{
		Strategy:            "EMA Crossover",
		RiskLevel:           0.5,
		MaxConcurrentTrades: 5,
		TakeProfitPercent:   2.0,
		StopLossPercent:     1.0,
		DailyTradeLimit:     20,
		TradingPairs:        []string{"BTCUSDT", "ETHUSDT", "SOLUSDT"},
		TradingSchedule: struct {
			Days      []string `json:"days"`
			StartTime string   `json:"start_time"`
			EndTime   string   `json:"end_time"`
		}{
			Days:      []string{"Mon", "Tue", "Wed", "Thu", "Fri"},
			StartTime: "09:00",
			EndTime:   "17:00",
		},
	}
	resp := &GetConfigResponse{Body: respBody}
	return resp, nil
}

// UpdateConfigHandler handles PUT requests to /api/v1/config using Huma signature.
func UpdateConfigHandler(ctx context.Context, input *UpdateConfigInput) (*UpdateConfigResponse, error) {
	// In real implementation, validate and save the updated config.
	// Huma handles validation based on struct tags.

	// Log the received config for now.
	// logger.Info("Received config update request", zap.Any("config", input.Body))

	// Return the updated config (received one as confirmation).
	resp := &UpdateConfigResponse{Body: input.Body}
	return resp, nil
}

// GetDefaultConfigHandler handles GET requests to /api/v1/config/default using Huma signature.
func GetDefaultConfigHandler(ctx context.Context, input *GetDefaultConfigInput) (*GetDefaultConfigResponse, error) {
	// Mock default data.
	respBody := ConfigResponseBody{
		Strategy:            "Default Strategy",
		RiskLevel:           0.3,
		MaxConcurrentTrades: 3,
		TakeProfitPercent:   1.5,
		StopLossPercent:     0.8,
		DailyTradeLimit:     15,
		TradingPairs:        []string{"BTCUSDT", "ETHUSDT"},
		TradingSchedule: struct {
			Days      []string `json:"days"`
			StartTime string   `json:"start_time"`
			EndTime   string   `json:"end_time"`
		}{
			Days:      []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
			StartTime: "00:00",
			EndTime:   "23:59",
		},
	}
	resp := &GetDefaultConfigResponse{Body: respBody}
	return resp, nil
}

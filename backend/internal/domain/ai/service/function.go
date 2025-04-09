package service

// FunctionDefinition defines a function that can be called by the AI
type FunctionDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Required    []string               `json:"required"`
}

// RegisterTradingFunctions registers available functions for the AI
func RegisterTradingFunctions() []FunctionDefinition {
	return []FunctionDefinition{
		{
			Name:        "get_market_data",
			Description: "Get current market data for a specific cryptocurrency",
			Parameters: map[string]interface{}{
				"symbol": map[string]interface{}{
					"type":        "string",
					"description": "The trading symbol (e.g., BTC, ETH)",
				},
				"timeframe": map[string]interface{}{
					"type":        "string",
					"description": "Timeframe for the data (e.g., 1h, 4h, 1d)",
					"enum":        []string{"1m", "5m", "15m", "1h", "4h", "1d", "1w"},
				},
			},
			Required: []string{"symbol"},
		},
		{
			Name:        "analyze_technical_indicators",
			Description: "Analyze technical indicators for a specific cryptocurrency",
			Parameters: map[string]interface{}{
				"symbol": map[string]interface{}{
					"type":        "string",
					"description": "The trading symbol (e.g., BTC, ETH)",
				},
				"indicators": map[string]interface{}{
					"type":        "array",
					"description": "List of indicators to analyze",
					"items": map[string]interface{}{
						"type": "string",
						"enum": []string{"rsi", "macd", "bollinger", "ema", "sma", "fibonacci"},
					},
				},
			},
			Required: []string{"symbol", "indicators"},
		},
		{
			Name:        "execute_trade",
			Description: "Execute a trade (requires confirmation)",
			Parameters: map[string]interface{}{
				"symbol": map[string]interface{}{
					"type":        "string",
					"description": "The trading symbol (e.g., BTC, ETH)",
				},
				"action": map[string]interface{}{
					"type":        "string",
					"description": "Buy or sell",
					"enum":        []string{"buy", "sell"},
				},
				"amount": map[string]interface{}{
					"type":        "number",
					"description": "Amount to trade (in USD or percentage of portfolio)",
				},
				"price_type": map[string]interface{}{
					"type":        "string",
					"description": "Market or limit order",
					"enum":        []string{"market", "limit"},
				},
				"limit_price": map[string]interface{}{
					"type":        "number",
					"description": "Price for limit orders",
				},
				"stop_loss": map[string]interface{}{
					"type":        "number",
					"description": "Stop loss price",
				},
				"take_profit": map[string]interface{}{
					"type":        "number",
					"description": "Take profit price",
				},
			},
			Required: []string{"symbol", "action", "amount", "price_type"},
		},
	}
}

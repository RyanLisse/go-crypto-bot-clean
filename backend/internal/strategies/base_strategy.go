package strategies

import (
	"context"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/indicators"
	"go-crypto-bot-clean/backend/internal/trading"

	"go.uber.org/zap"
)

// BaseStrategy provides a basic implementation of the trading.Strategy interface
type BaseStrategy struct {
	// Common fields for strategies
	Name             string                 // Strategy name
	Config           map[string]interface{} // Strategy configuration
	Positions        map[string]bool        // Track open positions by symbol
	PriceHistory     map[string][]float64   // Store price history for each symbol
	VolumeHistory    map[string][]float64   // Store volume history for each symbol
	MaxHistoryLength int                    // Maximum length of price history to keep
	Logger           *zap.Logger            // Logger for strategy events
	Metadata         map[string]interface{} // Additional strategy metadata
	IndicatorCache   map[string]interface{} // Cache for calculated indicators
}

// NewBaseStrategy creates a new BaseStrategy with the given name and configuration
func NewBaseStrategy(name string, config map[string]interface{}, logger *zap.Logger) *BaseStrategy {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}

	return &BaseStrategy{
		Name:             name,
		Config:           config,
		Positions:        make(map[string]bool),
		PriceHistory:     make(map[string][]float64),
		VolumeHistory:    make(map[string][]float64),
		MaxHistoryLength: 1000, // Default to 1000 data points
		Logger:           logger,
		Metadata:         make(map[string]interface{}),
		IndicatorCache:   make(map[string]interface{}),
	}
}

// Initialize sets up the strategy
func (s *BaseStrategy) Initialize(ctx context.Context, config interface{}) error {
	// If config is a map, update the strategy configuration
	if cfg, ok := config.(map[string]interface{}); ok {
		s.Config = cfg
	}

	// Initialize other fields if they're nil
	if s.Positions == nil {
		s.Positions = make(map[string]bool)
	}
	if s.PriceHistory == nil {
		s.PriceHistory = make(map[string][]float64)
	}
	if s.VolumeHistory == nil {
		s.VolumeHistory = make(map[string][]float64)
	}
	if s.Metadata == nil {
		s.Metadata = make(map[string]interface{})
	}
	if s.IndicatorCache == nil {
		s.IndicatorCache = make(map[string]interface{})
	}

	return nil
}

// OnTick processes a market tick
func (s *BaseStrategy) OnTick(ctx context.Context, symbol string, timestamp time.Time, data interface{}) ([]*trading.Signal, error) {
	// Update price and volume history if data is available
	price := s.extractPrice(data)
	if price > 0 {
		s.updatePriceHistory(symbol, price)
	}

	volume := s.extractVolume(data)
	if volume > 0 {
		s.updateVolumeHistory(symbol, volume)
	}

	// Base implementation returns no signals, must be overridden
	return nil, nil
}

// extractPrice attempts to extract a price from various data types
func (s *BaseStrategy) extractPrice(data interface{}) float64 {
	switch v := data.(type) {
	case *models.Kline:
		return v.Close
	case *models.Candle:
		return v.ClosePrice
	case *models.Ticker:
		return v.Price
	case *models.Trade:
		return v.Price
	case float64:
		return v
	case map[string]interface{}:
		if price, ok := v["price"].(float64); ok {
			return price
		}
		if price, ok := v["close"].(float64); ok {
			return price
		}
	}
	return 0
}

// extractVolume attempts to extract volume from various data types
func (s *BaseStrategy) extractVolume(data interface{}) float64 {
	switch v := data.(type) {
	case *models.Kline:
		return v.Volume
	case *models.Candle:
		return v.Volume
	case *models.Ticker:
		return v.Volume
	case *models.Trade:
		return v.Quantity
	case map[string]interface{}:
		if vol, ok := v["volume"].(float64); ok {
			return vol
		}
	}
	return 0
}

// updatePriceHistory adds a price to the history for a symbol
func (s *BaseStrategy) updatePriceHistory(symbol string, price float64) {
	// Initialize price history for this symbol if it doesn't exist
	if _, ok := s.PriceHistory[symbol]; !ok {
		s.PriceHistory[symbol] = make([]float64, 0, s.MaxHistoryLength)
	}

	// Add the price to the history
	s.PriceHistory[symbol] = append(s.PriceHistory[symbol], price)

	// Trim the history if it exceeds the maximum length
	if len(s.PriceHistory[symbol]) > s.MaxHistoryLength {
		s.PriceHistory[symbol] = s.PriceHistory[symbol][len(s.PriceHistory[symbol])-s.MaxHistoryLength:]
	}
}

// updateVolumeHistory adds a volume to the history for a symbol
func (s *BaseStrategy) updateVolumeHistory(symbol string, volume float64) {
	// Initialize volume history for this symbol if it doesn't exist
	if _, ok := s.VolumeHistory[symbol]; !ok {
		s.VolumeHistory[symbol] = make([]float64, 0, s.MaxHistoryLength)
	}

	// Add the volume to the history
	s.VolumeHistory[symbol] = append(s.VolumeHistory[symbol], volume)

	// Trim the history if it exceeds the maximum length
	if len(s.VolumeHistory[symbol]) > s.MaxHistoryLength {
		s.VolumeHistory[symbol] = s.VolumeHistory[symbol][len(s.VolumeHistory[symbol])-s.MaxHistoryLength:]
	}
}

// OnOrderFilled handles filled orders
func (s *BaseStrategy) OnOrderFilled(ctx context.Context, order *models.Order) error {
	// Update position tracking
	if order.Side == models.OrderSideBuy {
		s.Positions[order.Symbol] = true
		s.Logger.Info("Long position opened",
			zap.String("symbol", order.Symbol),
			zap.Float64("price", order.Price),
			zap.Float64("quantity", order.Quantity),
			zap.String("strategy", s.Name),
		)
	} else if order.Side == models.OrderSideSell {
		s.Positions[order.Symbol] = false
		s.Logger.Info("Position closed",
			zap.String("symbol", order.Symbol),
			zap.Float64("price", order.Price),
			zap.Float64("quantity", order.Quantity),
			zap.String("strategy", s.Name),
		)
	}
	return nil
}

// ClosePositions is called at the end of the backtest
func (s *BaseStrategy) ClosePositions(ctx context.Context) ([]*trading.Signal, error) {
	// Generate signals to close all open positions
	signals := make([]*trading.Signal, 0)

	for symbol, isOpen := range s.Positions {
		if isOpen {
			s.Logger.Info("Closing position at end of backtest",
				zap.String("symbol", symbol),
				zap.String("strategy", s.Name),
			)

			// Create a sell signal to close the position
			signal := &trading.Signal{
				Symbol:    symbol,
				Side:      models.OrderSideSell,
				Quantity:  0, // Will be filled in by the backtester based on position size
				Price:     0, // Will be filled in by the backtester based on current price
				Timestamp: time.Now(),
				Reason:    "End of backtest - closing position",
			}
			signals = append(signals, signal)
		}
	}

	return signals, nil
}

// CalculateIndicator calculates a technical indicator on the price history
func (s *BaseStrategy) CalculateIndicator(symbol string, indicator string, period int) ([]float64, error) {
	// Check if we have price history for this symbol
	prices, ok := s.PriceHistory[symbol]
	if !ok || len(prices) < period {
		return nil, fmt.Errorf("insufficient price history for %s", symbol)
	}

	// Check if the indicator is already cached
	cacheKey := fmt.Sprintf("%s_%s_%d", symbol, indicator, period)
	if cached, ok := s.IndicatorCache[cacheKey]; ok {
		if result, ok := cached.([]float64); ok {
			return result, nil
		}
	}

	// Calculate the requested indicator
	var result []float64
	var err error

	switch indicator {
	case "sma":
		result, err = indicators.SMA(prices, period)
	case "ema":
		result, err = indicators.EMA(prices, period)
	case "rsi":
		result, err = indicators.RSI(prices, period)
	case "bollinger_middle":
		_, middle, _, err := indicators.BollingerBands(prices, period, 2.0)
		if err != nil {
			return nil, err
		}
		result = middle
	case "bollinger_upper":
		upper, _, _, err := indicators.BollingerBands(prices, period, 2.0)
		if err != nil {
			return nil, err
		}
		result = upper
	case "bollinger_lower":
		_, _, lower, err := indicators.BollingerBands(prices, period, 2.0)
		if err != nil {
			return nil, err
		}
		result = lower
	case "macd":
		macd, _, _, err := indicators.MACD(prices, 12, 26, 9)
		if err != nil {
			return nil, err
		}
		result = macd
	case "macd_signal":
		_, signal, _, err := indicators.MACD(prices, 12, 26, 9)
		if err != nil {
			return nil, err
		}
		result = signal
	case "macd_histogram":
		_, _, histogram, err := indicators.MACD(prices, 12, 26, 9)
		if err != nil {
			return nil, err
		}
		result = histogram
	default:
		return nil, fmt.Errorf("unsupported indicator: %s", indicator)
	}

	if err != nil {
		return nil, err
	}

	// Cache the result
	s.IndicatorCache[cacheKey] = result

	return result, nil
}

// GetConfig returns the strategy configuration
func (s *BaseStrategy) GetConfig() map[string]interface{} {
	return s.Config
}

// SetConfig updates the strategy configuration
func (s *BaseStrategy) SetConfig(config map[string]interface{}) {
	s.Config = config
}

// GetConfigValue gets a configuration value by key
func (s *BaseStrategy) GetConfigValue(key string) (interface{}, bool) {
	value, ok := s.Config[key]
	return value, ok
}

// SetConfigValue sets a configuration value by key
func (s *BaseStrategy) SetConfigValue(key string, value interface{}) {
	s.Config[key] = value
}

// GetFloatConfigValue gets a float64 configuration value by key
func (s *BaseStrategy) GetFloatConfigValue(key string, defaultValue float64) float64 {
	if value, ok := s.Config[key]; ok {
		if floatValue, ok := value.(float64); ok {
			return floatValue
		}
	}
	return defaultValue
}

// GetIntConfigValue gets an int configuration value by key
func (s *BaseStrategy) GetIntConfigValue(key string, defaultValue int) int {
	if value, ok := s.Config[key]; ok {
		if intValue, ok := value.(int); ok {
			return intValue
		}
		if floatValue, ok := value.(float64); ok {
			return int(floatValue)
		}
	}
	return defaultValue
}

// GetStringConfigValue gets a string configuration value by key
func (s *BaseStrategy) GetStringConfigValue(key string, defaultValue string) string {
	if value, ok := s.Config[key]; ok {
		if stringValue, ok := value.(string); ok {
			return stringValue
		}
	}
	return defaultValue
}

// GetBoolConfigValue gets a bool configuration value by key
func (s *BaseStrategy) GetBoolConfigValue(key string, defaultValue bool) bool {
	if value, ok := s.Config[key]; ok {
		if boolValue, ok := value.(bool); ok {
			return boolValue
		}
	}
	return defaultValue
}

// ClearIndicatorCache clears the indicator cache
func (s *BaseStrategy) ClearIndicatorCache() {
	s.IndicatorCache = make(map[string]interface{})
}

// GetName returns the name of the strategy
func (s *BaseStrategy) GetName() string {
	return s.Name
}

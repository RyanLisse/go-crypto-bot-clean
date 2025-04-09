package advanced

import (
	"context"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/types"
)

// NewCoinConfig represents the configuration for the NewCoinStrategy
type NewCoinConfig struct {
	MinVolume         float64 `json:"minVolume"`         // Minimum volume threshold
	MinPrice          float64 `json:"minPrice"`          // Minimum price threshold
	MaxPrice          float64 `json:"maxPrice"`          // Maximum price threshold
	EntryDelay        int     `json:"entryDelay"`        // Delay in minutes before entry
	StopLossPercent   float64 `json:"stopLossPercent"`   // Stop loss percentage
	TakeProfitPercent float64 `json:"takeProfitPercent"` // Take profit percentage
	PositionSize      float64 `json:"positionSize"`      // Position size as percentage of capital
	MaxAge            int     `json:"maxAge"`            // Maximum age in minutes to track a coin
}

// NewCoinStrategy implements a strategy for trading newly listed coins
type NewCoinStrategy struct {
	config        *NewCoinConfig
	detectedCoins map[string]time.Time // Map of detected coins and their first seen timestamp
}

// NewNewCoinStrategy creates a new NewCoinStrategy instance
func NewNewCoinStrategy(config *NewCoinConfig) *NewCoinStrategy {
	return &NewCoinStrategy{
		config:        config,
		detectedCoins: make(map[string]time.Time),
	}
}

// GetName returns the name of the strategy
func (s *NewCoinStrategy) GetName() string {
	return "NewCoinStrategy"
}

// Initialize prepares the strategy with initial configuration
func (s *NewCoinStrategy) Initialize(ctx context.Context, config map[string]interface{}) error {
	// Convert config to NewCoinConfig
	if config == nil {
		return fmt.Errorf("config is required")
	}

	// Set default values if not provided
	s.config = &NewCoinConfig{
		MinVolume:         1000,    // Default minimum volume
		MinPrice:          0.00001, // Default minimum price
		MaxPrice:          1.0,     // Default maximum price
		EntryDelay:        5,       // Default 5 minutes delay
		StopLossPercent:   5.0,     // Default 5% stop loss
		TakeProfitPercent: 10.0,    // Default 10% take profit
		PositionSize:      1.0,     // Default 1% position size
		MaxAge:            60,      // Default 60 minutes tracking
	}

	// Override defaults with provided values
	if v, ok := config["minVolume"].(float64); ok {
		s.config.MinVolume = v
	}
	if v, ok := config["minPrice"].(float64); ok {
		s.config.MinPrice = v
	}
	if v, ok := config["maxPrice"].(float64); ok {
		s.config.MaxPrice = v
	}
	if v, ok := config["entryDelay"].(int); ok {
		s.config.EntryDelay = v
	}
	if v, ok := config["stopLossPercent"].(float64); ok {
		s.config.StopLossPercent = v
	}
	if v, ok := config["takeProfitPercent"].(float64); ok {
		s.config.TakeProfitPercent = v
	}
	if v, ok := config["positionSize"].(float64); ok {
		s.config.PositionSize = v
	}
	if v, ok := config["maxAge"].(int); ok {
		s.config.MaxAge = v
	}

	return nil
}

// UpdateParameters updates the strategy parameters
func (s *NewCoinStrategy) UpdateParameters(ctx context.Context, params map[string]interface{}) error {
	return s.Initialize(ctx, params)
}

// OnTickerUpdate processes a new ticker update and may generate a signal
func (s *NewCoinStrategy) OnTickerUpdate(ctx context.Context, ticker *models.Ticker) (*types.Signal, error) {
	if ticker == nil {
		return nil, fmt.Errorf("ticker is nil")
	}

	// Check if we've seen this coin before
	firstSeen, exists := s.detectedCoins[ticker.Symbol]
	now := time.Now()

	// If we haven't seen this coin before and it meets our criteria, track it
	if !exists {
		if ticker.Volume >= s.config.MinVolume &&
			ticker.Price >= s.config.MinPrice &&
			ticker.Price <= s.config.MaxPrice {
			s.detectedCoins[ticker.Symbol] = now
			return nil, nil // No signal yet, just started tracking
		}
		return nil, nil // Not interested in this coin
	}

	// Calculate how long we've been tracking this coin
	trackingDuration := now.Sub(firstSeen)

	// Remove old coins from tracking
	if trackingDuration.Minutes() > float64(s.config.MaxAge) {
		delete(s.detectedCoins, ticker.Symbol)
		return nil, nil
	}

	// Check if it's time to enter a position
	if trackingDuration.Minutes() >= float64(s.config.EntryDelay) {
		// Generate buy signal
		stopLoss := ticker.Price * (1 - s.config.StopLossPercent/100)
		takeProfit := ticker.Price * (1 + s.config.TakeProfitPercent/100)

		signal := &types.Signal{
			Symbol:          ticker.Symbol,
			Type:            types.BUY,
			Confidence:      0.8, // High confidence for new listings
			Price:           ticker.Price,
			StopLoss:        stopLoss,
			TakeProfit:      takeProfit,
			Timestamp:       now,
			ExpirationTime:  now.Add(time.Hour), // Signal valid for 1 hour
			RecommendedSize: s.config.PositionSize,
			Metadata: map[string]interface{}{
				"strategy":     "NewCoin",
				"trackingTime": trackingDuration.Minutes(),
				"volume":       ticker.Volume,
			},
		}

		// Remove from tracking after generating signal
		delete(s.detectedCoins, ticker.Symbol)
		return signal, nil
	}

	return nil, nil
}

// OnCandleUpdate processes a new candle and may generate a signal
func (s *NewCoinStrategy) OnCandleUpdate(ctx context.Context, candle *models.Candle) (*types.Signal, error) {
	// Not used for this strategy
	return nil, nil
}

// OnTradeUpdate processes a new trade and may generate a signal
func (s *NewCoinStrategy) OnTradeUpdate(ctx context.Context, trade *models.Trade) (*types.Signal, error) {
	// Not used for this strategy
	return nil, nil
}

// OnMarketDepthUpdate processes a new market depth update and may generate a signal
func (s *NewCoinStrategy) OnMarketDepthUpdate(ctx context.Context, depth *models.OrderBook) (*types.Signal, error) {
	// Not used for this strategy
	return nil, nil
}

// OnTimerEvent processes a scheduled timer event
func (s *NewCoinStrategy) OnTimerEvent(ctx context.Context, eventType string) (*types.Signal, error) {
	// Could add periodic cleanup of old detected coins
	return nil, nil
}

// GetTimeframes returns the timeframes this strategy requires
func (s *NewCoinStrategy) GetTimeframes() []string {
	return []string{"1m"} // We mainly work with real-time ticker updates
}

// GetRequiredDataTypes returns the types of data this strategy needs
func (s *NewCoinStrategy) GetRequiredDataTypes() []string {
	return []string{"ticker"}
}

// PerformBacktest runs a backtest of the strategy
func (s *NewCoinStrategy) PerformBacktest(ctx context.Context, historicalData []*models.Candle, params map[string]interface{}) ([]*types.Signal, interface{}, error) {
	return nil, nil, fmt.Errorf("backtesting not implemented for NewCoinStrategy")
}

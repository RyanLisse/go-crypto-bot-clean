package advanced

import (
	"context"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/types"
)

// NewCoinStrategy configuration parameters
type NewCoinConfig struct {
	MinVolume         float64 `json:"min_volume"`          // Minimum volume threshold
	MinPrice          float64 `json:"min_price"`           // Minimum price threshold
	MaxPrice          float64 `json:"max_price"`           // Maximum price threshold
	EntryDelay        int     `json:"entry_delay"`         // Delay in minutes before entering after detection
	StopLossPercent   float64 `json:"stop_loss_percent"`   // Stop loss percentage
	TakeProfitPercent float64 `json:"take_profit_percent"` // Take profit percentage
	PositionSize      float64 `json:"position_size"`       // Position size as percentage of portfolio
	MaxAge            int     `json:"max_age"`             // Maximum age in hours to consider a coin "new"
}

// NewCoinStrategy implements a strategy for trading newly listed coins
type NewCoinStrategy struct {
	config        NewCoinConfig
	detectedCoins map[string]time.Time // Map of detected coins and their first seen timestamp
}

// NewNewCoinStrategy creates a new instance of NewCoinStrategy
func NewNewCoinStrategy(config map[string]interface{}) (*NewCoinStrategy, error) {
	// Parse and validate configuration
	cfg := NewCoinConfig{
		MinVolume:         1000.0,  // Default $1000 minimum volume
		MinPrice:          0.00001, // Default minimum price
		MaxPrice:          10.0,    // Default maximum price
		EntryDelay:        5,       // Default 5 minutes delay
		StopLossPercent:   5.0,     // Default 5% stop loss
		TakeProfitPercent: 10.0,    // Default 10% take profit
		PositionSize:      1.0,     // Default 1% of portfolio
		MaxAge:            24,      // Default 24 hours
	}

	// Override defaults with provided config
	if v, ok := config["min_volume"].(float64); ok {
		cfg.MinVolume = v
	}
	if v, ok := config["min_price"].(float64); ok {
		cfg.MinPrice = v
	}
	if v, ok := config["max_price"].(float64); ok {
		cfg.MaxPrice = v
	}
	if v, ok := config["entry_delay"].(int); ok {
		cfg.EntryDelay = v
	}
	if v, ok := config["stop_loss_percent"].(float64); ok {
		cfg.StopLossPercent = v
	}
	if v, ok := config["take_profit_percent"].(float64); ok {
		cfg.TakeProfitPercent = v
	}
	if v, ok := config["position_size"].(float64); ok {
		cfg.PositionSize = v
	}
	if v, ok := config["max_age"].(int); ok {
		cfg.MaxAge = v
	}

	return &NewCoinStrategy{
		config:        cfg,
		detectedCoins: make(map[string]time.Time),
	}, nil
}

// GetName returns the name of the strategy
func (s *NewCoinStrategy) GetName() string {
	return "NewCoinStrategy"
}

// Initialize prepares the strategy
func (s *NewCoinStrategy) Initialize(ctx context.Context, config map[string]interface{}) error {
	// Already initialized in constructor
	return nil
}

// UpdateParameters updates the strategy parameters
func (s *NewCoinStrategy) UpdateParameters(ctx context.Context, params map[string]interface{}) error {
	newStrategy, err := NewNewCoinStrategy(params)
	if err != nil {
		return err
	}
	s.config = newStrategy.config
	return nil
}

// OnTickerUpdate processes a new ticker update
func (s *NewCoinStrategy) OnTickerUpdate(ctx context.Context, ticker *models.Ticker) (*types.Signal, error) {
	// Check if this is a new coin we haven't seen before
	firstSeen, exists := s.detectedCoins[ticker.Symbol]
	if !exists {
		s.detectedCoins[ticker.Symbol] = time.Now()
		return nil, nil // Just record it for now
	}

	// Check if the coin is still within our "new" window
	if time.Since(firstSeen).Hours() > float64(s.config.MaxAge) {
		delete(s.detectedCoins, ticker.Symbol) // Remove old coins
		return nil, nil
	}

	// Check if it meets our criteria
	if ticker.Price < s.config.MinPrice || ticker.Price > s.config.MaxPrice {
		return nil, nil
	}

	if ticker.Volume < s.config.MinVolume {
		return nil, nil
	}

	// Check if enough time has passed since detection
	if time.Since(firstSeen).Minutes() < float64(s.config.EntryDelay) {
		return nil, nil
	}

	// Generate buy signal
	signal := &types.Signal{
		Symbol:          ticker.Symbol,
		Type:            types.SignalBuy,
		Confidence:      0.8, // High confidence for new coins meeting criteria
		Price:           ticker.Price,
		StopLoss:        ticker.Price * (1 - s.config.StopLossPercent/100),
		TakeProfit:      ticker.Price * (1 + s.config.TakeProfitPercent/100),
		Timestamp:       time.Now(),
		ExpirationTime:  time.Now().Add(1 * time.Hour),
		RecommendedSize: s.config.PositionSize,
		Metadata: map[string]interface{}{
			"strategy":  "NewCoin",
			"firstSeen": firstSeen,
		},
	}

	return signal, nil
}

// OnCandleUpdate processes a new candle
func (s *NewCoinStrategy) OnCandleUpdate(ctx context.Context, candle *models.Candle) (*types.Signal, error) {
	// We primarily work with ticker updates, but could add volume analysis here
	return nil, nil
}

// OnTradeUpdate processes a new trade
func (s *NewCoinStrategy) OnTradeUpdate(ctx context.Context, trade *models.Trade) (*types.Signal, error) {
	// Not used for this strategy
	return nil, nil
}

// OnMarketDepthUpdate processes a new market depth update
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
	return []string{"ticker", "volume"}
}

// PerformBacktest runs a backtest of the strategy
func (s *NewCoinStrategy) PerformBacktest(ctx context.Context, historicalData []*models.Candle, params map[string]interface{}) ([]*types.Signal, *models.BacktestResult, error) {
	return nil, nil, fmt.Errorf("backtesting not implemented for NewCoinStrategy")
}

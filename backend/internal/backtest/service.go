package backtest

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
)

// Service provides backtest functionality
type Service struct {
	backtests map[string]*BacktestResult
	factory   *StrategyFactory
}

// NewService creates a new backtest service
func NewService() *Service {
	factory := NewStrategyFactory()

	return &Service{
		backtests: make(map[string]*BacktestResult),
		factory:   factory,
	}
}

// BacktestRequestConfig contains configuration for initiating a backtest via API/service
type BacktestRequestConfig struct {
	Strategy       string    `json:"strategy"`
	Symbol         string    `json:"symbol"` // Assuming single symbol for simplicity in this request model
	Timeframe      string    `json:"timeframe"`
	StartTime      time.Time `json:"startTime"`
	EndTime        time.Time `json:"endTime"`
	InitialCapital float64   `json:"initialCapital"`
	RiskPerTrade   float64   `json:"riskPerTrade"` // Example strategy-specific param
}

// RunBacktest runs a backtest with the given configuration
// Note: This is currently a mock implementation
func (s *Service) RunBacktest(ctx context.Context, reqConfig *BacktestRequestConfig) (*BacktestResult, error) {
	// Create strategy instance from factory
	strategy, err := s.factory.CreateStrategy(reqConfig.Strategy)
	if err != nil {
		return nil, fmt.Errorf("failed to create strategy '%s': %w", reqConfig.Strategy, err)
	}

	// TODO: Map BacktestRequestConfig to the Engine's BacktestConfig
	// This requires deciding how to handle potentially multiple symbols, data providers, etc.
	engineConfig := &BacktestConfig{
		StartTime:      reqConfig.StartTime,
		EndTime:        reqConfig.EndTime,
		InitialCapital: reqConfig.InitialCapital,
		Symbols:        []string{reqConfig.Symbol}, // Assuming single symbol for now
		Interval:       reqConfig.Timeframe,
		Strategy:       strategy,
		// Logger:         Needs a logger instance
		// DataProvider:   Needs a data provider instance
		// CommissionRate: Needs a default or configured value
		// SlippageModel:  Needs a default or configured value
	}

	// Generate a unique ID for this backtest (can be improved)
	id := fmt.Sprintf("%s-%s-%s-%d", reqConfig.Strategy, reqConfig.Symbol, reqConfig.Timeframe, time.Now().Unix())

	// -------- MOCK IMPLEMENTATION START --------
	// TODO: Replace this mock section with actual Engine execution:
	// engine := NewEngine(engineConfig)
	// result, err := engine.Run(ctx)
	// if err != nil { ... }

	// Create a mock backtest result for now
	result := &BacktestResult{
		Config:         engineConfig, // Store the engine config
		StartTime:      reqConfig.StartTime,
		EndTime:        reqConfig.EndTime,
		InitialCapital: reqConfig.InitialCapital,
		FinalCapital:   reqConfig.InitialCapital * (1 + rand.Float64()*0.5), // Random final capital
	}

	// Generate mock data
	result.EquityCurve = generateMockEquityCurve(reqConfig.StartTime, reqConfig.EndTime, reqConfig.InitialCapital, result.FinalCapital)
	result.DrawdownCurve = generateMockDrawdownCurve(result.EquityCurve)
	result.ClosedPositions = generateMockTrades(reqConfig.StartTime, reqConfig.EndTime, reqConfig.Symbol, 20)
	result.Trades = generateMockOrders(reqConfig.StartTime, reqConfig.EndTime, reqConfig.Symbol, 20)

	// Calculate mock performance metrics
	analyzer := NewPerformanceAnalyzer()
	metrics, err := analyzer.CalculateMetrics(result) // Use the mock result
	if err != nil {
		return nil, fmt.Errorf("failed to calculate mock metrics: %w", err)
	}
	result.PerformanceMetrics = metrics
	// -------- MOCK IMPLEMENTATION END --------

	// Store the result (using the generated ID)
	s.backtests[id] = result

	return result, nil
}

// GetBacktestResult gets a backtest result by ID
func (s *Service) GetBacktestResult(ctx context.Context, id string) (*BacktestResult, error) {
	result, ok := s.backtests[id]
	if !ok {
		return nil, fmt.Errorf("backtest with ID %s not found", id)
	}
	return result, nil
}

// ListBacktestResults lists all backtest results
func (s *Service) ListBacktestResults(ctx context.Context) ([]*BacktestResult, error) {
	results := make([]*BacktestResult, 0, len(s.backtests))
	for _, result := range s.backtests {
		results = append(results, result)
	}
	return results, nil
}

// generateMockEquityCurve generates a mock equity curve
func generateMockEquityCurve(startTime, endTime time.Time, initialCapital, finalCapital float64) []*EquityPoint {
	// Calculate the number of days between start and end
	days := int(endTime.Sub(startTime).Hours() / 24)
	if days < 1 {
		days = 1
	}

	// Create equity curve with daily points
	equityCurve := make([]*EquityPoint, 0, days)

	// Add initial point
	equityCurve = append(equityCurve, &EquityPoint{
		Timestamp: startTime,
		Equity:    initialCapital,
	})

	// Calculate daily change
	totalChange := finalCapital - initialCapital
	avgDailyChange := totalChange / float64(days)
	volatility := avgDailyChange * 0.5 // 50% volatility

	// Generate daily points
	currentEquity := initialCapital
	for i := 1; i < days; i++ {
		// Calculate random daily change with some volatility
		dailyChange := avgDailyChange + rand.Float64()*volatility*2 - volatility
		currentEquity += dailyChange

		// Ensure equity doesn't go below 0
		if currentEquity < 0 {
			currentEquity = initialCapital * 0.1 // 90% drawdown at worst
		}

		// Add point
		equityCurve = append(equityCurve, &EquityPoint{
			Timestamp: startTime.Add(time.Duration(i) * 24 * time.Hour),
			Equity:    currentEquity,
		})
	}

	// Add final point
	equityCurve = append(equityCurve, &EquityPoint{
		Timestamp: endTime,
		Equity:    finalCapital,
	})

	return equityCurve
}

// generateMockDrawdownCurve generates a mock drawdown curve
func generateMockDrawdownCurve(equityCurve []*EquityPoint) []*DrawdownPoint {
	drawdownCurve := make([]*DrawdownPoint, 0, len(equityCurve))
	highWaterMark := equityCurve[0].Equity

	for _, point := range equityCurve {
		// Update high water mark
		if point.Equity > highWaterMark {
			highWaterMark = point.Equity
		}

		// Calculate drawdown
		drawdown := highWaterMark - point.Equity

		// Add point
		drawdownCurve = append(drawdownCurve, &DrawdownPoint{
			Timestamp: point.Timestamp,
			Drawdown:  drawdown,
		})
	}

	return drawdownCurve
}

// generateMockTrades generates mock trades
func generateMockTrades(startTime, endTime time.Time, symbol string, count int) []*models.ClosedPosition {
	trades := make([]*models.ClosedPosition, 0, count)
	duration := endTime.Sub(startTime)
	timePerTrade := duration / time.Duration(count)

	for i := 0; i < count; i++ {
		// Calculate open and close times
		openTime := startTime.Add(time.Duration(i) * timePerTrade)
		closeTime := openTime.Add(timePerTrade / 2)

		// Determine if trade is profitable (70% win rate)
		isProfitable := rand.Float64() < 0.7

		// Calculate entry and exit prices
		entryPrice := 10000.0 + rand.Float64()*1000.0
		var exitPrice float64
		var profit float64
		if isProfitable {
			exitPrice = entryPrice * (1 + rand.Float64()*0.1) // Up to 10% profit
			profit = exitPrice - entryPrice
		} else {
			exitPrice = entryPrice * (1 - rand.Float64()*0.05) // Up to 5% loss
			profit = exitPrice - entryPrice
		}

		// Create trade
		trade := &models.ClosedPosition{
			ID:         fmt.Sprintf("mock-%d", i+1), // ID is string
			Symbol:     symbol,
			Side:       models.OrderSideBuy, // Use enum
			Quantity:   1.0,
			EntryPrice: entryPrice,
			ExitPrice:  exitPrice,
			ProfitLoss: profit, // Use ProfitLoss field
			OpenTime:   openTime,
			CloseTime:  closeTime,
		}

		trades = append(trades, trade)
	}

	return trades
}

// generateMockOrders generates mock orders
func generateMockOrders(startTime, endTime time.Time, symbol string, count int) []*models.Order {
	orders := make([]*models.Order, 0, count*2) // Buy and sell for each trade
	duration := endTime.Sub(startTime)
	timePerTrade := duration / time.Duration(count)

	for i := 0; i < count; i++ {
		// Calculate open and close times
		openTime := startTime.Add(time.Duration(i) * timePerTrade)
		closeTime := openTime.Add(timePerTrade / 2)

		// Calculate prices
		price := 10000.0 + rand.Float64()*1000.0

		// Create buy order
		buyOrder := &models.Order{
			ID:        fmt.Sprintf("mock-buy-%d", i+1), // ID is string
			Symbol:    symbol,
			Side:      models.OrderSideBuy,   // Use enum
			Type:      models.OrderTypeLimit, // Use enum
			Price:     price,
			Quantity:  1.0,
			Status:    models.OrderStatusFilled, // Use enum
			CreatedAt: openTime,
		}

		// Create sell order
		sellOrder := &models.Order{
			ID:        fmt.Sprintf("mock-sell-%d", i+1), // ID is string
			Symbol:    symbol,
			Side:      models.OrderSideSell,                     // Use enum
			Type:      models.OrderTypeLimit,                    // Use enum
			Price:     price * (1 + (rand.Float64()*0.2 - 0.1)), // +/- 10%
			Quantity:  1.0,
			Status:    models.OrderStatusFilled, // Use enum
			CreatedAt: closeTime,
		}

		orders = append(orders, buyOrder, sellOrder)
	}

	return orders
}

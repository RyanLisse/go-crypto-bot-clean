package backtest

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"go.uber.org/zap"
)

// BacktestEventType defines the type of backtest event
type BacktestEventType string

const (
	EventStrategyInitialized BacktestEventType = "strategy_initialized"
	EventDataLoaded          BacktestEventType = "data_loaded"
	EventSignalGenerated     BacktestEventType = "signal_generated"
	EventOrderCreated        BacktestEventType = "order_created"
	EventOrderFilled         BacktestEventType = "order_filled"
	EventPositionOpened      BacktestEventType = "position_opened"
	EventPositionClosed      BacktestEventType = "position_closed"
	EventError               BacktestEventType = "error"
)

// BacktestEvent represents an event that occurred during the backtest
type BacktestEvent struct {
	Type      BacktestEventType
	Timestamp time.Time
	Symbol    string
	Data      interface{}
}

// BacktestConfig contains configuration options for a backtest
type BacktestConfig struct {
	StartTime          time.Time
	EndTime            time.Time
	InitialCapital     float64
	Symbols            []string
	Interval           string
	CommissionRate     float64
	SlippageModel      SlippageModel
	EnableShortSelling bool
	DataProvider       DataProvider
	Strategy           BacktestStrategy
	Logger             *zap.Logger
}

// BacktestResult contains the results of a backtest
type BacktestResult struct {
	Config             *BacktestConfig
	StartTime          time.Time
	EndTime            time.Time
	InitialCapital     float64
	FinalCapital       float64
	Trades             []*models.Order
	Positions          []*models.Position
	ClosedPositions    []*models.ClosedPosition
	EquityCurve        []*EquityPoint
	DrawdownCurve      []*DrawdownPoint
	Events             []*BacktestEvent
	PerformanceMetrics *PerformanceMetrics
}

// EquityPoint represents a point on the equity curve
type EquityPoint struct {
	Timestamp time.Time
	Equity    float64
}

// DrawdownPoint represents a point on the drawdown curve
type DrawdownPoint struct {
	Timestamp time.Time
	Drawdown  float64
}

// Engine implements the backtesting engine
type Engine struct {
	config          *BacktestConfig
	events          []*BacktestEvent
	positionTracker PositionTracker
	cash            float64
	equity          float64
	equityCurve     []*EquityPoint
	drawdownCurve   []*DrawdownPoint
	trades          []*models.Order
	currentPrices   map[string]float64
	logger          *zap.Logger
}

// NewEngine creates a new backtesting engine
func NewEngine(config *BacktestConfig) *Engine {
	logger := config.Logger
	if logger == nil {
		logger, _ = zap.NewDevelopment()
	}

	return &Engine{
		config:          config,
		events:          make([]*BacktestEvent, 0),
		positionTracker: NewPositionTracker(),
		cash:            config.InitialCapital,
		equity:          config.InitialCapital,
		equityCurve:     make([]*EquityPoint, 0),
		drawdownCurve:   make([]*DrawdownPoint, 0),
		trades:          make([]*models.Order, 0),
		currentPrices:   make(map[string]float64),
		logger:          logger,
	}
}

// Run executes a backtest with the given configuration
func (e *Engine) Run(ctx context.Context) (*BacktestResult, error) {
	startTime := time.Now()
	e.logger.Info("Starting backtest",
		zap.Time("start_time", e.config.StartTime),
		zap.Time("end_time", e.config.EndTime),
		zap.Float64("initial_capital", e.config.InitialCapital),
		zap.Strings("symbols", e.config.Symbols),
		zap.String("interval", e.config.Interval),
	)

	// Initialize the strategy
	err := e.config.Strategy.Initialize(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize strategy: %w", err)
	}
	e.addEvent(EventStrategyInitialized, time.Now(), "", nil)

	// Load historical data for all symbols
	data, err := e.loadHistoricalData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load historical data: %w", err)
	}
	e.addEvent(EventDataLoaded, time.Now(), "", data)

	// Sort all data points by timestamp
	timePoints := e.createTimePoints(data)

	// Process each time point
	for _, tp := range timePoints {
		// Update current prices
		for symbol, price := range tp.prices {
			e.currentPrices[symbol] = price
		}

		// Process signals from the strategy
		signals, err := e.config.Strategy.OnTick(ctx, tp.symbol, tp.timestamp, tp.data)
		if err != nil {
			e.logger.Error("Error processing tick",
				zap.Error(err),
				zap.String("symbol", tp.symbol),
				zap.Time("timestamp", tp.timestamp),
			)
			e.addEvent(EventError, tp.timestamp, tp.symbol, err)
			continue
		}

		// Process signals
		for _, signal := range signals {
			e.addEvent(EventSignalGenerated, tp.timestamp, signal.Symbol, signal)

			// Create and execute order
			order := e.createOrder(signal)
			e.addEvent(EventOrderCreated, tp.timestamp, signal.Symbol, order)

			// Execute the order
			filledOrder, err := e.executeOrder(order)
			if err != nil {
				e.logger.Error("Error executing order",
					zap.Error(err),
					zap.String("symbol", order.Symbol),
					zap.String("side", string(order.Side)),
					zap.Float64("quantity", order.Quantity),
					zap.Float64("price", order.Price),
				)
				e.addEvent(EventError, tp.timestamp, order.Symbol, err)
				continue
			}
			e.addEvent(EventOrderFilled, tp.timestamp, order.Symbol, filledOrder)

			// Notify strategy of order fill
			err = e.config.Strategy.OnOrderFilled(ctx, filledOrder)
			if err != nil {
				e.logger.Error("Error notifying strategy of order fill",
					zap.Error(err),
					zap.String("symbol", filledOrder.Symbol),
					zap.String("side", string(filledOrder.Side)),
					zap.Float64("quantity", filledOrder.Quantity),
					zap.Float64("price", filledOrder.Price),
				)
				e.addEvent(EventError, tp.timestamp, filledOrder.Symbol, err)
			}
		}

		// Update equity and drawdown
		e.updateEquity(tp.timestamp)
	}

	// Calculate final equity
	finalEquity := e.cash
	for _, position := range e.positionTracker.GetOpenPositions() {
		price, ok := e.currentPrices[position.Symbol]
		if ok {
			finalEquity += position.Quantity * price
		}
	}

	// Create backtest result
	result := &BacktestResult{
		Config:          e.config,
		StartTime:       e.config.StartTime,
		EndTime:         e.config.EndTime,
		InitialCapital:  e.config.InitialCapital,
		FinalCapital:    finalEquity,
		Trades:          e.trades,
		Positions:       e.positionTracker.GetOpenPositions(),
		ClosedPositions: e.positionTracker.GetClosedPositions(),
		EquityCurve:     e.equityCurve,
		DrawdownCurve:   e.drawdownCurve,
		Events:          e.events,
	}

	// Calculate performance metrics
	analyzer := NewPerformanceAnalyzer()
	metrics, err := analyzer.CalculateMetrics(result)
	if err != nil {
		e.logger.Error("Error calculating performance metrics", zap.Error(err))
	} else {
		result.PerformanceMetrics = metrics
	}

	e.logger.Info("Backtest completed",
		zap.Duration("duration", time.Since(startTime)),
		zap.Float64("final_capital", finalEquity),
		zap.Float64("return", (finalEquity-e.config.InitialCapital)/e.config.InitialCapital*100),
		zap.Int("trades", len(e.trades)),
		zap.Int("closed_positions", len(e.positionTracker.GetClosedPositions())),
		zap.Int("open_positions", len(e.positionTracker.GetOpenPositions())),
	)

	return result, nil
}

// GetEvents returns all events generated during the backtest
func (e *Engine) GetEvents() []*BacktestEvent {
	return e.events
}

// GetPositions returns all positions created during the backtest
func (e *Engine) GetPositions() []*models.Position {
	return e.positionTracker.GetOpenPositions()
}

// GetTrades returns all trades executed during the backtest
func (e *Engine) GetTrades() []*models.Order {
	return e.trades
}

// addEvent adds an event to the event log
func (e *Engine) addEvent(eventType BacktestEventType, timestamp time.Time, symbol string, data interface{}) {
	event := &BacktestEvent{
		Type:      eventType,
		Timestamp: timestamp,
		Symbol:    symbol,
		Data:      data,
	}
	e.events = append(e.events, event)
}

// timePoint represents a point in time with data for a specific symbol
type timePoint struct {
	timestamp time.Time
	symbol    string
	data      interface{}
	prices    map[string]float64
}

// loadHistoricalData loads historical data for all symbols
func (e *Engine) loadHistoricalData(ctx context.Context) (map[string][]*models.Kline, error) {
	result := make(map[string][]*models.Kline)

	for _, symbol := range e.config.Symbols {
		klines, err := e.config.DataProvider.GetKlines(ctx, symbol, e.config.Interval, e.config.StartTime, e.config.EndTime)
		if err != nil {
			return nil, fmt.Errorf("failed to get klines for %s: %w", symbol, err)
		}

		result[symbol] = klines
		e.logger.Info("Loaded historical data",
			zap.String("symbol", symbol),
			zap.Int("klines", len(klines)),
			zap.String("interval", e.config.Interval),
		)
	}

	return result, nil
}

// createTimePoints creates a sorted list of time points from historical data
func (e *Engine) createTimePoints(data map[string][]*models.Kline) []timePoint {
	var timePoints []timePoint

	for symbol, klines := range data {
		for _, kline := range klines {
			timestamp := kline.OpenTime
			prices := make(map[string]float64)
			prices[symbol] = kline.Close

			tp := timePoint{
				timestamp: timestamp,
				symbol:    symbol,
				data:      kline,
				prices:    prices,
			}
			timePoints = append(timePoints, tp)
		}
	}

	// Sort time points by timestamp
	sort.Slice(timePoints, func(i, j int) bool {
		return timePoints[i].timestamp.Before(timePoints[j].timestamp)
	})

	return timePoints
}

// createOrder creates an order from a signal
func (e *Engine) createOrder(signal *Signal) *models.Order {
	return &models.Order{
		ID:        uuid.New().String(),
		Symbol:    signal.Symbol,
		Side:      models.OrderSide(signal.Side),
		Type:      "MARKET",
		Quantity:  signal.Quantity,
		Price:     signal.Price,
		Status:    "NEW",
		CreatedAt: signal.Timestamp,
		Time:      signal.Timestamp,
	}
}

// executeOrder executes an order and updates positions
func (e *Engine) executeOrder(order *models.Order) (*models.Order, error) {
	// Apply slippage to the price
	slippage := 0.0
	if e.config.SlippageModel != nil {
		slippage = e.config.SlippageModel.CalculateSlippage(order.Symbol, string(order.Side), order.Quantity, order.Price, order.Time)
	}

	// Calculate execution price
	executionPrice := order.Price
	if order.Side == "BUY" {
		executionPrice += slippage
	} else {
		executionPrice -= slippage
	}

	// Calculate commission
	commission := executionPrice * order.Quantity * e.config.CommissionRate

	// Update order
	filledOrder := *order
	filledOrder.Price = executionPrice
	filledOrder.Status = "FILLED"

	// Update cash
	orderValue := executionPrice * order.Quantity
	if order.Side == "BUY" {
		// Check if we have enough cash
		if e.cash < orderValue+commission {
			return nil, fmt.Errorf("insufficient funds: required %.2f, available %.2f", orderValue+commission, e.cash)
		}
		e.cash -= orderValue + commission

		// Open position
		position, err := e.positionTracker.OpenPosition(order.Symbol, string(order.Side), executionPrice, order.Quantity, order.Time)
		if err != nil {
			return nil, fmt.Errorf("failed to open position: %w", err)
		}
		e.addEvent(EventPositionOpened, order.Time, order.Symbol, position)
	} else {
		// Find position to close
		positions := e.positionTracker.GetOpenPositions()
		var position *models.Position
		for _, p := range positions {
			if p.Symbol == order.Symbol {
				position = p
				break
			}
		}

		if position == nil {
			return nil, fmt.Errorf("no open position found for %s", order.Symbol)
		}

		// Close position
		closedPosition, err := e.positionTracker.ClosePosition(position.ID, executionPrice, order.Time)
		if err != nil {
			return nil, fmt.Errorf("failed to close position: %w", err)
		}
		e.addEvent(EventPositionClosed, order.Time, order.Symbol, closedPosition)

		// Update cash
		e.cash += orderValue - commission

		// Notify strategy of position close
		err = e.config.Strategy.OnPositionClosed(context.Background(), closedPosition)
		if err != nil {
			e.logger.Error("Error notifying strategy of position close",
				zap.Error(err),
				zap.String("symbol", closedPosition.Symbol),
				zap.Float64("profit", closedPosition.ProfitLoss),
			)
		}
	}

	// Add trade to list
	e.trades = append(e.trades, &filledOrder)

	return &filledOrder, nil
}

// updateEquity updates the equity curve and drawdown curve
func (e *Engine) updateEquity(timestamp time.Time) {
	// Calculate current equity
	equity := e.cash
	for _, position := range e.positionTracker.GetOpenPositions() {
		price, ok := e.currentPrices[position.Symbol]
		if ok {
			equity += position.Quantity * price
		}
	}

	// Update equity
	e.equity = equity

	// Add point to equity curve
	e.equityCurve = append(e.equityCurve, &EquityPoint{
		Timestamp: timestamp,
		Equity:    equity,
	})

	// Calculate drawdown
	maxEquity := e.config.InitialCapital
	for _, point := range e.equityCurve {
		if point.Equity > maxEquity {
			maxEquity = point.Equity
		}
	}

	drawdown := (maxEquity - equity) / maxEquity * 100
	e.drawdownCurve = append(e.drawdownCurve, &DrawdownPoint{
		Timestamp: timestamp,
		Drawdown:  drawdown,
	})
}

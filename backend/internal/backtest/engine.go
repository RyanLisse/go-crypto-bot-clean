package backtest

import (
	"context"
	"fmt"
	"sort"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/pkg/position_tracker"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Engine implements the backtesting engine
type Engine struct {
	config          *BacktestConfig
	events          []*BacktestEvent
	positions       map[string]*models.Position
	closedPositions []*models.ClosedPosition
	cash            float64
	equity          float64
	equityCurve     []*EquityPoint
	drawdownCurve   []*DrawdownPoint
	logger          *zap.Logger
	positionTracker *position_tracker.PositionTracker
	trades          []*models.Order
	strategy        BacktestStrategy
	currentPrices   map[string]float64
}

// NewEngine creates a new backtesting engine with the given configuration
func NewEngine(config *BacktestConfig) *Engine {
	pt := position_tracker.NewPositionTracker(config.Logger)
	return &Engine{
		config:          config,
		events:          make([]*BacktestEvent, 0),
		positions:       make(map[string]*models.Position),
		closedPositions: make([]*models.ClosedPosition, 0),
		cash:            config.InitialCapital,
		equity:          config.InitialCapital,
		equityCurve:     make([]*EquityPoint, 0),
		drawdownCurve:   make([]*DrawdownPoint, 0),
		logger:          config.Logger,
		positionTracker: pt,
		trades:          make([]*models.Order, 0),
		strategy:        config.Strategy,
		currentPrices:   make(map[string]float64),
	}
}

// Run executes the backtest
func (e *Engine) Run(ctx context.Context) (*BacktestResult, error) {
	startTime := time.Now()
	e.addEvent(BacktestStarted, startTime, "", nil)
	e.logger.Info("Starting backtest", zap.Time("start_time", startTime))

	if err := e.strategy.Initialize(ctx, nil); err != nil {
		e.addEvent(BacktestError, time.Now(), "", err)
		return nil, fmt.Errorf("failed to initialize strategy: %w", err)
	}

	historicalData, err := e.loadHistoricalData(ctx)
	if err != nil {
		e.addEvent(BacktestError, time.Now(), "", err)
		return nil, fmt.Errorf("failed to load historical data: %w", err)
	}

	timePoints := e.createTimePoints(historicalData)

	e.updateEquity(e.config.StartTime)

	for _, tp := range timePoints {
		for symbol, price := range tp.prices {
			e.currentPrices[symbol] = price
		}

		signals, err := e.strategy.OnTick(ctx, tp.symbol, tp.timestamp, tp.data)
		if err != nil {
			e.logger.Error("Strategy OnTick error", zap.Error(err), zap.Time("timestamp", tp.timestamp), zap.String("symbol", tp.symbol))
			e.addEvent(BacktestError, tp.timestamp, tp.symbol, err)
			return nil, fmt.Errorf("strategy OnTick error at %s for %s: %w", tp.timestamp, tp.symbol, err)
		}

		for _, signal := range signals {
			order := e.createOrder(signal)
			filledOrder, err := e.executeOrder(order)
			if err != nil {
				e.logger.Error("Order execution error", zap.Error(err), zap.String("order_id", order.ID))
				e.addEvent(BacktestError, order.Time, order.Symbol, err)
				return nil, fmt.Errorf("order execution failed for signal at %s: %w", signal.Timestamp, err)
			}

			if err := e.strategy.OnOrderFilled(ctx, filledOrder); err != nil {
				e.logger.Error("Strategy OnOrderFilled error", zap.Error(err), zap.String("order_id", filledOrder.ID))
			}
		}

		e.updateEquity(tp.timestamp)
	}

	closingSignals, err := e.strategy.ClosePositions(ctx)
	if err != nil {
		e.addEvent(BacktestError, time.Now(), "", fmt.Errorf("failed to get closing signals: %w", err))
		return nil, fmt.Errorf("failed to get closing signals: %w", err)
	}

	closingTime := e.config.EndTime
	if len(timePoints) > 0 {
		closingTime = timePoints[len(timePoints)-1].timestamp
	}
	for _, signal := range closingSignals {
		signal.Timestamp = closingTime
		order := e.createOrder(signal)
		filledOrder, err := e.executeOrder(order)
		if err != nil {
			e.logger.Error("Closing order execution error", zap.Error(err), zap.String("order_id", order.ID))
			e.addEvent(BacktestError, order.Time, order.Symbol, err)
		} else {
			if err := e.strategy.OnOrderFilled(ctx, filledOrder); err != nil {
				e.logger.Error("Strategy OnOrderFilled error for closing order", zap.Error(err), zap.String("order_id", filledOrder.ID))
			}
		}
	}

	e.updateEquity(closingTime)

	metrics := e.calculatePerformanceMetrics()

	finalOpenPositions := make([]*models.Position, 0)
	finalOpenPositions = append(finalOpenPositions, e.positionTracker.GetOpenPositions()...)

	e.closedPositions = e.positionTracker.GetClosedPositions()

	endTime := time.Now()
	result := &BacktestResult{
		Config:             e.config,
		StartTime:          startTime,
		EndTime:            endTime,
		InitialCapital:     e.config.InitialCapital,
		FinalCapital:       e.equity,
		Trades:             e.trades,
		Positions:          finalOpenPositions,
		ClosedPositions:    e.closedPositions,
		EquityCurve:        e.equityCurve,
		DrawdownCurve:      e.drawdownCurve,
		Events:             e.events,
		PerformanceMetrics: metrics,
	}

	e.logger.Info("Backtest completed", zap.Time("end_time", endTime), zap.Duration("duration", endTime.Sub(startTime)))
	e.addEvent(BacktestCompleted, endTime, "", result)

	return result, nil
}

// calculatePerformanceMetrics calculates various performance metrics for the backtest
func (e *Engine) calculatePerformanceMetrics() *PerformanceMetrics {
	return &PerformanceMetrics{}
}

// GetEvents returns all events generated during the backtest
func (e *Engine) GetEvents() []*BacktestEvent {
	return e.events
}

// GetPositions returns all currently open positions during the backtest execution (snapshot)
func (e *Engine) GetPositions() []*models.Position {
	return e.positionTracker.GetOpenPositions()
}

// GetTrades returns all filled orders executed during the backtest
func (e *Engine) GetTrades() []*models.Order {
	return e.trades
}

// addEvent adds an event to the event log
func (e *Engine) addEvent(eventType BacktestEventType, timestamp time.Time, symbol string, data interface{}) {
	if timestamp.IsZero() {
		timestamp = time.Now()
	}
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

	if e.config.DataProvider == nil {
		return nil, fmt.Errorf("DataProvider is not configured")
	}

	// Check if DataProvider implements the expected method (optional but good practice)
	provider, ok := e.config.DataProvider.(interface {
		GetHistoricalData(ctx context.Context, symbol string, interval string, start, end time.Time) ([]*models.Kline, error)
	})
	if !ok {
		// Handle case where the DataProvider doesn't have the exact method signature if needed,
		// or trust that the type assertion succeeded if it's guaranteed.
		// For now, let's assume it should implement GetHistoricalData.
		return nil, fmt.Errorf("DataProvider does not implement GetHistoricalData method")
	}

	for _, symbol := range e.config.Symbols {
		klines, err := provider.GetHistoricalData(ctx, symbol, e.config.Interval, e.config.StartTime, e.config.EndTime)
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
	priceMap := make(map[time.Time]map[string]float64)
	dataMap := make(map[time.Time]map[string]interface{})
	timestamps := make(map[time.Time]struct{})

	for symbol, klines := range data {
		for _, kline := range klines {
			timestamp := kline.OpenTime
			timestamps[timestamp] = struct{}{}

			if _, ok := priceMap[timestamp]; !ok {
				priceMap[timestamp] = make(map[string]float64)
			}
			priceMap[timestamp][symbol] = kline.Close

			if _, ok := dataMap[timestamp]; !ok {
				dataMap[timestamp] = make(map[string]interface{})
			}
			dataMap[timestamp][symbol] = kline
		}
	}

	sortedTimestamps := make([]time.Time, 0, len(timestamps))
	for ts := range timestamps {
		sortedTimestamps = append(sortedTimestamps, ts)
	}
	sort.Slice(sortedTimestamps, func(i, j int) bool {
		return sortedTimestamps[i].Before(sortedTimestamps[j])
	})

	for _, ts := range sortedTimestamps {
		var primarySymbol string
		if len(dataMap[ts]) > 0 {
			symbols := make([]string, 0, len(dataMap[ts]))
			for sym := range dataMap[ts] {
				symbols = append(symbols, sym)
			}
			sort.Strings(symbols)
			primarySymbol = symbols[0]
		}

		tp := timePoint{
			timestamp: ts,
			symbol:    primarySymbol,
			data:      dataMap[ts][primarySymbol],
			prices:    priceMap[ts],
		}
		timePoints = append(timePoints, tp)
	}

	return timePoints
}

// createOrder creates an order from a signal
func (e *Engine) createOrder(signal *Signal) *models.Order {
	return &models.Order{
		ID:        uuid.New().String(),
		Symbol:    signal.Symbol,
		Side:      models.OrderSide(signal.Side),
		Type:      models.OrderTypeMarket,
		Quantity:  signal.Quantity,
		Price:     signal.Price,
		Status:    models.OrderStatusNew,
		CreatedAt: signal.Timestamp,
		Time:      signal.Timestamp,
	}
}

// executeOrder executes an order and updates positions
func (e *Engine) executeOrder(order *models.Order) (*models.Order, error) {
	executionPrice, ok := e.currentPrices[order.Symbol]
	if !ok {
		return nil, fmt.Errorf("cannot execute order, no current market price available for %s at %s", order.Symbol, order.Time)
	}

	slippage := 0.0
	if e.config.SlippageModel != nil {
		slippage = e.config.SlippageModel.CalculateSlippage(executionPrice, order.Quantity, order.Side)
	}

	if order.Side == models.OrderSideBuy {
		executionPrice += slippage
	} else {
		executionPrice -= slippage
	}

	commission := executionPrice * order.Quantity * e.config.CommissionRate

	filledOrder := *order
	filledOrder.Price = executionPrice
	filledOrder.Status = models.OrderStatusFilled
	filledOrder.FilledQty = order.Quantity
	filledOrder.UpdatedAt = time.Now()

	orderValue := executionPrice * filledOrder.FilledQty
	if filledOrder.Side == models.OrderSideBuy {
		requiredCash := orderValue + commission
		if e.cash < requiredCash {
			filledOrder.Status = models.OrderStatusRejected
			return &filledOrder, fmt.Errorf("insufficient funds: required %.2f, available %.2f", requiredCash, e.cash)
		}
		e.cash -= requiredCash

		position, err := e.positionTracker.OpenPosition(order.Symbol, string(order.Side), executionPrice, order.Quantity, order.Time)
		if err != nil {
			e.logger.Error("Failed to open position via tracker", zap.Error(err), zap.String("symbol", order.Symbol))
			return &filledOrder, fmt.Errorf("failed to open position: %w", err)
		}
		e.positions[order.Symbol] = position
		e.addEvent(BacktestProgress, order.Time, order.Symbol, position)
	} else {
		e.cash += orderValue - commission

		closedPosition, err := e.positionTracker.ClosePositionBySymbol(order.Symbol, executionPrice, order.Quantity, order.Time)
		if err != nil {
			e.logger.Error("Failed to close position via tracker", zap.Error(err), zap.String("symbol", order.Symbol))
			filledOrder.Status = models.OrderStatusRejected
			return &filledOrder, fmt.Errorf("failed to close position for %s: %w", order.Symbol, err)
		}
		delete(e.positions, order.Symbol)
		e.closedPositions = append(e.closedPositions, closedPosition)
		e.addEvent(BacktestProgress, order.Time, order.Symbol, closedPosition)
	}

	e.trades = append(e.trades, &filledOrder)

	return &filledOrder, nil
}

// updateEquity updates the equity curve and drawdown curve
func (e *Engine) updateEquity(timestamp time.Time) {
	currentEquity := e.cash
	for _, position := range e.positionTracker.GetOpenPositions() {
		price, ok := e.currentPrices[position.Symbol]
		if ok {
			if position.Side == models.OrderSideBuy {
				// For long positions, calculate current value
				currentEquity += position.Quantity * price
			} else {
				// For short positions, calculate PnL and add to equity
				pnl := (position.EntryPrice - price) * position.Quantity
				currentEquity += position.Quantity*position.EntryPrice + pnl
			}
		} else {
			currentEquity += position.Quantity * position.EntryPrice
		}
	}

	e.equity = currentEquity

	e.equityCurve = append(e.equityCurve, &EquityPoint{
		Timestamp: timestamp,
		Equity:    currentEquity,
	})

	maxEquity := e.config.InitialCapital
	for _, point := range e.equityCurve {
		if point.Equity > maxEquity {
			maxEquity = point.Equity
		}
	}

	drawdown := 0.0
	if maxEquity > 0 {
		drawdown = (maxEquity - currentEquity) / maxEquity * 100
	}
	e.drawdownCurve = append(e.drawdownCurve, &DrawdownPoint{
		Timestamp: timestamp,
		Drawdown:  drawdown,
	})
}

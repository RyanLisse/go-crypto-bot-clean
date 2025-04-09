package backtest

import (
	"context"
	"fmt"
	"sort"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// EventType defines the type of event in the event-driven engine
type EventType string

const (
	EventTypeMarketData    EventType = "market_data"
	EventTypeSignal        EventType = "signal"
	EventTypeOrder         EventType = "order"
	EventTypeOrderFilled   EventType = "order_filled"
	EventTypePositionOpen  EventType = "position_open"
	EventTypePositionClose EventType = "position_close"
)

// Event represents an event in the event-driven engine
type Event struct {
	Type      EventType
	Timestamp time.Time
	Symbol    string
	Data      interface{}
}

// MarketEvent represents a market data event
type MarketEvent struct {
	Symbol    string
	Timestamp time.Time
	Data      interface{} // Can be Kline, Ticker, OrderBook, etc.
}

// EventDrivenStrategy extends the BacktestStrategy interface with event handling
type EventDrivenStrategy interface {
	BacktestStrategy

	// OnMarketEvent is called for each market event
	OnMarketEvent(ctx context.Context, event *MarketEvent) ([]*Signal, error)
}

// EventDrivenEngineConfig contains configuration options for the event-driven engine
type EventDrivenEngineConfig struct {
	StartTime          time.Time
	EndTime            time.Time
	InitialCapital     float64
	Symbols            []string
	Interval           string
	FeeModel           FeeModel
	SlippageModel      SlippageModel
	EnableShortSelling bool
	DataLoader         *DatabaseDataLoader
	Strategy           EventDrivenStrategy
	DB                 *gorm.DB
	Logger             *zap.Logger
}

// EventDrivenEngine implements an event-driven backtesting engine
type EventDrivenEngine struct {
	config          *EventDrivenEngineConfig
	events          []*Event
	positionTracker PositionTracker
	cash            float64
	equity          float64
	equityCurve     []*EquityPoint
	drawdownCurve   []*DrawdownPoint
	trades          []*models.Order
	currentPrices   map[string]float64
	logger          *zap.Logger
	db              *gorm.DB
	eventQueue      []*Event
}

// NewEventDrivenEngine creates a new event-driven backtesting engine
func NewEventDrivenEngine(config *EventDrivenEngineConfig) *EventDrivenEngine {
	logger := config.Logger
	if logger == nil {
		logger, _ = zap.NewDevelopment()
	}

	return &EventDrivenEngine{
		config:          config,
		events:          make([]*Event, 0),
		positionTracker: NewPositionTracker(),
		cash:            config.InitialCapital,
		equity:          config.InitialCapital,
		equityCurve:     make([]*EquityPoint, 0),
		drawdownCurve:   make([]*DrawdownPoint, 0),
		trades:          make([]*models.Order, 0),
		currentPrices:   make(map[string]float64),
		logger:          logger,
		db:              config.DB,
		eventQueue:      make([]*Event, 0),
	}
}

// Run executes a backtest with the event-driven engine
func (e *EventDrivenEngine) Run(ctx context.Context) (*BacktestResult, error) {
	startTime := time.Now()
	e.logger.Info("Starting event-driven backtest",
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

	// Load historical data for all symbols
	datasets, err := e.loadHistoricalData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load historical data: %w", err)
	}

	// Create market events from historical data
	e.createMarketEvents(datasets)

	// Sort events by timestamp
	sort.Slice(e.eventQueue, func(i, j int) bool {
		return e.eventQueue[i].Timestamp.Before(e.eventQueue[j].Timestamp)
	})

	// Process events
	for len(e.eventQueue) > 0 {
		// Get the next event
		event := e.eventQueue[0]
		e.eventQueue = e.eventQueue[1:]

		// Process the event
		err := e.processEvent(ctx, event)
		if err != nil {
			e.logger.Error("Error processing event",
				zap.Error(err),
				zap.String("event_type", string(event.Type)),
				zap.Time("timestamp", event.Timestamp),
				zap.String("symbol", event.Symbol),
			)
		}
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
		Config: &BacktestConfig{
			StartTime:          e.config.StartTime,
			EndTime:            e.config.EndTime,
			InitialCapital:     e.config.InitialCapital,
			Symbols:            e.config.Symbols,
			Interval:           e.config.Interval,
			CommissionRate:     0.001, // Default value, actual fees are calculated by the fee model
			SlippageModel:      e.config.SlippageModel,
			EnableShortSelling: e.config.EnableShortSelling,
			DataProvider:       nil, // Not used in event-driven engine
			Strategy:           e.config.Strategy,
			Logger:             e.logger,
		},
		StartTime:       e.config.StartTime,
		EndTime:         e.config.EndTime,
		InitialCapital:  e.config.InitialCapital,
		FinalCapital:    finalEquity,
		Trades:          e.trades,
		Positions:       e.positionTracker.GetOpenPositions(),
		ClosedPositions: e.positionTracker.GetClosedPositions(),
		EquityCurve:     e.equityCurve,
		DrawdownCurve:   e.drawdownCurve,
		Events:          nil, // We don't store BacktestEvent objects in the event-driven engine
	}

	// Calculate performance metrics
	analyzer := NewPerformanceAnalyzer()
	metrics, err := analyzer.CalculateMetrics(result)
	if err != nil {
		e.logger.Error("Error calculating performance metrics", zap.Error(err))
	} else {
		result.PerformanceMetrics = metrics
	}

	// Save backtest result to database
	err = e.saveBacktestResult(result)
	if err != nil {
		e.logger.Error("Error saving backtest result to database", zap.Error(err))
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

// loadHistoricalData loads historical data for all symbols
func (e *EventDrivenEngine) loadHistoricalData(ctx context.Context) (map[string]*DataSet, error) {
	result := make(map[string]*DataSet)

	for _, symbol := range e.config.Symbols {
		dataset, err := e.config.DataLoader.LoadData(ctx, symbol, e.config.Interval, e.config.StartTime, e.config.EndTime)
		if err != nil {
			return nil, fmt.Errorf("failed to load data for %s: %w", symbol, err)
		}

		result[symbol] = dataset
		e.logger.Info("Loaded historical data",
			zap.String("symbol", symbol),
			zap.Int("klines", len(dataset.Klines)),
			zap.String("interval", e.config.Interval),
		)
	}

	return result, nil
}

// createMarketEvents creates market events from historical data
func (e *EventDrivenEngine) createMarketEvents(datasets map[string]*DataSet) {
	for symbol, dataset := range datasets {
		for _, kline := range dataset.Klines {
			// Create market event
			marketEvent := &MarketEvent{
				Symbol:    symbol,
				Timestamp: kline.OpenTime,
				Data:      kline,
			}

			// Add to event queue
			e.eventQueue = append(e.eventQueue, &Event{
				Type:      EventTypeMarketData,
				Timestamp: kline.OpenTime,
				Symbol:    symbol,
				Data:      marketEvent,
			})

			// Update current price
			e.currentPrices[symbol] = kline.Close
		}
	}
}

// processEvent processes an event
func (e *EventDrivenEngine) processEvent(ctx context.Context, event *Event) error {
	// Add event to event log
	e.events = append(e.events, event)

	// Process event based on type
	switch event.Type {
	case EventTypeMarketData:
		return e.processMarketEvent(ctx, event)
	case EventTypeSignal:
		return e.processSignalEvent(ctx, event)
	case EventTypeOrder:
		return e.processOrderEvent(ctx, event)
	case EventTypeOrderFilled:
		return e.processOrderFilledEvent(ctx, event)
	case EventTypePositionOpen:
		return e.processPositionOpenEvent(ctx, event)
	case EventTypePositionClose:
		return e.processPositionCloseEvent(ctx, event)
	default:
		return fmt.Errorf("unknown event type: %s", event.Type)
	}
}

// processMarketEvent processes a market event
func (e *EventDrivenEngine) processMarketEvent(ctx context.Context, event *Event) error {
	marketEvent, ok := event.Data.(*MarketEvent)
	if !ok {
		return fmt.Errorf("invalid market event data")
	}

	// Update current prices if the data is a kline
	if kline, ok := marketEvent.Data.(*models.Kline); ok {
		e.currentPrices[marketEvent.Symbol] = kline.Close
	}

	// Call strategy's OnMarketEvent method
	signals, err := e.config.Strategy.OnMarketEvent(ctx, marketEvent)
	if err != nil {
		return fmt.Errorf("error calling strategy OnMarketEvent: %w", err)
	}

	// Call strategy's OnTick method for backward compatibility
	tickSignals, err := e.config.Strategy.OnTick(ctx, marketEvent.Symbol, marketEvent.Timestamp, marketEvent.Data)
	if err != nil {
		return fmt.Errorf("error calling strategy OnTick: %w", err)
	}

	// Combine signals from both methods
	allSignals := append(signals, tickSignals...)

	// Process signals
	for _, signal := range allSignals {
		// Create signal event
		signalEvent := &Event{
			Type:      EventTypeSignal,
			Timestamp: signal.Timestamp,
			Symbol:    signal.Symbol,
			Data:      signal,
		}

		// Add to event queue
		e.eventQueue = append(e.eventQueue, signalEvent)
	}

	// Update equity and drawdown
	e.updateEquity(event.Timestamp)

	return nil
}

// processSignalEvent processes a signal event
func (e *EventDrivenEngine) processSignalEvent(ctx context.Context, event *Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		signal, ok := event.Data.(*Signal)
		if !ok {
			return fmt.Errorf("invalid signal event data")
		}

		// Create order from signal
		order := e.createOrder(signal)

		// Create order event
		orderEvent := &Event{
			Type:      EventTypeOrder,
			Timestamp: event.Timestamp,
			Symbol:    signal.Symbol,
			Data:      order,
		}

		// Add to event queue
		e.eventQueue = append(e.eventQueue, orderEvent)

		return nil
	}
}

// processOrderEvent processes an order event
func (e *EventDrivenEngine) processOrderEvent(ctx context.Context, event *Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		order, ok := event.Data.(*models.Order)
		if !ok {
			return fmt.Errorf("invalid order event data")
		}

		// Execute the order
		filledOrder, err := e.executeOrder(order)
		if err != nil {
			return fmt.Errorf("error executing order: %w", err)
		}

		// Create order filled event
		orderFilledEvent := &Event{
			Type:      EventTypeOrderFilled,
			Timestamp: event.Timestamp,
			Symbol:    order.Symbol,
			Data:      filledOrder,
		}

		// Add to event queue
		e.eventQueue = append(e.eventQueue, orderFilledEvent)

		return nil
	}
}

// processOrderFilledEvent handles the order filled event
func (e *EventDrivenEngine) processOrderFilledEvent(ctx context.Context, event *Event) error {
	order, ok := event.Data.(*models.Order)
	if !ok {
		return fmt.Errorf("invalid data type for order filled event: %T", event.Data)
	}

	// Notify the strategy that the order was filled
	err := e.config.Strategy.OnOrderFilled(ctx, order)
	if err != nil {
		return fmt.Errorf("strategy OnOrderFilled error: %w", err)
	}

	// NOTE: Position tracking logic (opening/closing positions based on fills)
	// should likely happen elsewhere, perhaps triggered by the strategy's response
	// or in a dedicated position management component that listens for fill events.
	// Removing the direct position tracker updates from this function.

	return nil
}

// processPositionOpenEvent handles the position open event
// (This function seems unused based on current flow, might be for future extension)
func (e *EventDrivenEngine) processPositionOpenEvent(ctx context.Context, event *Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Implementation depends on how position open events are generated and used
		e.logger.Debug("Processing position open event", zap.Any("event", event))
		return nil
	}
}

// processPositionCloseEvent handles the position close event
// (This function seems unused based on current flow, might be for future extension)
func (e *EventDrivenEngine) processPositionCloseEvent(ctx context.Context, event *Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Implementation depends on how position close events are generated and used
		closedPosition, ok := event.Data.(*models.ClosedPosition)
		if !ok {
			return fmt.Errorf("invalid data type for position close event: %T", event.Data)
		}
		e.logger.Debug("Processing position close event", zap.Any("closedPosition", closedPosition))

		// Potentially notify the strategy about the closed position here if needed,
		// but OnPositionClosed is not part of the EventDrivenStrategy interface.
		// Consider adding it to the interface if required by strategies.

		return nil
	}
}

// createOrder creates an order from a signal
func (e *EventDrivenEngine) createOrder(signal *Signal) *models.Order {
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
func (e *EventDrivenEngine) executeOrder(order *models.Order) (*models.Order, error) {
	// Apply slippage to the price
	slippage := 0.0
	if e.config.SlippageModel != nil {
		// Use the interface defined in common_types.go
		slippage = e.config.SlippageModel.CalculateSlippage(order.Price, order.Quantity, order.Side)
	}

	// Calculate execution price
	executionPrice := order.Price
	if order.Side == "BUY" {
		executionPrice += slippage
	} else {
		executionPrice -= slippage
	}

	// Calculate fee
	fee := 0.0
	if e.config.FeeModel != nil {
		fee = e.config.FeeModel.CalculateFee(order.Symbol, string(order.Side), order.Quantity, executionPrice, order.Time)
	}

	// Update order
	filledOrder := *order
	filledOrder.Price = executionPrice
	filledOrder.Status = "FILLED"

	// Update cash
	orderValue := executionPrice * order.Quantity
	if order.Side == "BUY" {
		// Check if we have enough cash
		if e.cash < orderValue+fee {
			return nil, fmt.Errorf("insufficient funds: required %.2f, available %.2f", orderValue+fee, e.cash)
		}
		e.cash -= orderValue + fee

		// Open position
		position, err := e.positionTracker.OpenPosition(order.Symbol, string(order.Side), executionPrice, order.Quantity, order.Time)
		if err != nil {
			return nil, fmt.Errorf("failed to open position: %w", err)
		}

		// Create position open event
		positionOpenEvent := &Event{
			Type:      EventTypePositionOpen,
			Timestamp: order.Time,
			Symbol:    order.Symbol,
			Data:      position,
		}

		// Add to event queue
		e.eventQueue = append(e.eventQueue, positionOpenEvent)
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

		// Create position close event
		positionCloseEvent := &Event{
			Type:      EventTypePositionClose,
			Timestamp: order.Time,
			Symbol:    order.Symbol,
			Data:      closedPosition,
		}

		// Add to event queue
		e.eventQueue = append(e.eventQueue, positionCloseEvent)

		// Update cash
		e.cash += orderValue - fee
	}

	// Add trade to list
	e.trades = append(e.trades, &filledOrder)

	// Add event for order filled
	e.addOrderFilledEvent(&filledOrder)

	return &filledOrder, nil
}

// addOrderFilledEvent adds an order filled event to the queue
func (e *EventDrivenEngine) addOrderFilledEvent(order *models.Order) {
	event := &Event{
		Type:      EventTypeOrderFilled,
		Timestamp: order.UpdatedAt, // Use UpdatedAt for fill time if available
		Symbol:    order.Symbol,
		Data:      order,
	}
	e.eventQueue = append(e.eventQueue, event)
}

// updateEquity updates the equity curve and drawdown curve
func (e *EventDrivenEngine) updateEquity(timestamp time.Time) {
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

// saveBacktestResult saves the backtest result to the database
func (e *EventDrivenEngine) saveBacktestResult(result *BacktestResult) error {
	// Create a database model from the backtest result
	backtestResult := models.BacktestResult{
		ID:                 uuid.New().String(),
		StrategyID:         "event_driven_strategy", // This should be configurable
		StartTime:          result.StartTime,
		EndTime:            result.EndTime,
		InitialBalance:     result.InitialCapital,
		FinalBalance:       result.FinalCapital,
		ProfitLoss:         result.FinalCapital - result.InitialCapital,
		ProfitLossPercent:  (result.FinalCapital - result.InitialCapital) / result.InitialCapital * 100,
		TotalTrades:        len(result.Trades),
		WinningTrades:      result.PerformanceMetrics.WinningTrades,
		LosingTrades:       result.PerformanceMetrics.LosingTrades,
		WinRate:            result.PerformanceMetrics.WinRate,
		AverageWin:         result.PerformanceMetrics.AverageProfitTrade,
		AverageLoss:        result.PerformanceMetrics.AverageLossTrade,
		MaxDrawdown:        result.PerformanceMetrics.MaxDrawdown,
		MaxDrawdownPercent: result.PerformanceMetrics.MaxDrawdownPercent,
		SharpeRatio:        result.PerformanceMetrics.SharpeRatio,
		SortinoRatio:       result.PerformanceMetrics.SortinoRatio,
	}

	// Save to database
	return e.db.Create(&backtestResult).Error
}

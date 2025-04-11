package execution

import (
	"context"
	"log"
	"sync"
	"time"
)

// SignalFilterPlugin defines a plugin interface for validating/filtering signals.
/*
SignalFilterPlugin defines a plugin interface for validating and prioritizing signals.

Validate returns true if the signal passes the plugin's criteria.
Priority returns an integer score for prioritization; higher means higher priority.
*/
type SignalFilterPlugin interface {
	Validate(signal Signal) bool
	Priority(signal Signal) int
}

// OrderRoutingPlugin defines a plugin interface for routing orders.
/*
OrderRoutingPlugin defines a plugin interface for routing orders.

Route can modify the order or decide routing logic.
Return an error to block or reroute the order.
*/
type OrderRoutingPlugin interface {
	Route(order *Order) error
}

// AdvancedSignalProcessor processes signals asynchronously with plugin-based validation, prioritization, and order generation.
/*
AdvancedSignalProcessor asynchronously processes signals with plugin-based validation,
prioritization, order generation, routing, and risk checks.

Use RegisterFilterPlugin and RegisterRoutingPlugin to extend behavior.
Submit signals via SubmitSignal(). Start() launches processing loop.
*/
type AdvancedSignalProcessor struct {
	executor       *StrategyExecutor
	filterPlugins  []SignalFilterPlugin
	routingPlugins []OrderRoutingPlugin
	signalQueue    chan Signal
	wg             sync.WaitGroup
	ctx            context.Context
	cancel         context.CancelFunc
	mu             sync.Mutex
}

// NewAdvancedSignalProcessor creates a new AdvancedSignalProcessor.
//
// executor: the StrategyExecutor to place orders through
// queueSize: size of the buffered signal queue
func NewAdvancedSignalProcessor(executor *StrategyExecutor, queueSize int) *AdvancedSignalProcessor {
	ctx, cancel := context.WithCancel(context.Background())
	return &AdvancedSignalProcessor{
		executor:    executor,
		signalQueue: make(chan Signal, queueSize),
		ctx:         ctx,
		cancel:      cancel,
	}
}

// RegisterFilterPlugin adds a signal filter plugin.
/*
RegisterFilterPlugin adds a signal filter plugin.

Plugins can validate signals and assign priority scores.
*/
func (asp *AdvancedSignalProcessor) RegisterFilterPlugin(plugin SignalFilterPlugin) {
	asp.mu.Lock()
	defer asp.mu.Unlock()
	asp.filterPlugins = append(asp.filterPlugins, plugin)
}

// RegisterRoutingPlugin adds an order routing plugin.
/*
RegisterRoutingPlugin adds an order routing plugin.

Plugins can modify or reroute orders before placement.
*/
func (asp *AdvancedSignalProcessor) RegisterRoutingPlugin(plugin OrderRoutingPlugin) {
	asp.mu.Lock()
	defer asp.mu.Unlock()
	asp.routingPlugins = append(asp.routingPlugins, plugin)
}

// SubmitSignal enqueues a signal for processing.
/*
SubmitSignal enqueues a signal for asynchronous processing.

If the queue is full, the signal is dropped and logged.
*/
func (asp *AdvancedSignalProcessor) SubmitSignal(signal Signal) {
	select {
	case asp.signalQueue <- signal:
	default:
		log.Printf("Signal queue full, dropping signal: %+v", signal)
	}
}

// Start launches the asynchronous processing loop.
/*
Start launches the asynchronous processing loop.

Call Stop() to shut down gracefully.
*/
func (asp *AdvancedSignalProcessor) Start() {
	asp.wg.Add(1)
	go asp.processLoop()
}

// Stop stops the processor gracefully.
/*
Stop signals the processor to stop and waits for shutdown.
*/
func (asp *AdvancedSignalProcessor) Stop() {
	asp.cancel()
	asp.wg.Wait()
}

func (asp *AdvancedSignalProcessor) processLoop() {
	defer asp.wg.Done()
	for {
		select {
		case <-asp.ctx.Done():
			return
		case signal := <-asp.signalQueue:
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Recovered in signal processing: %v", r)
					}
				}()
				asp.processSignal(signal)
			}()
		}
	}
}

func (asp *AdvancedSignalProcessor) processSignal(signal Signal) {
	// Run filter plugins
	priority := 0
	valid := true
	asp.mu.Lock()
	for _, plugin := range asp.filterPlugins {
		if !plugin.Validate(signal) {
			valid = false
			break
		}
		p := plugin.Priority(signal)
		if p > priority {
			priority = p
		}
	}
	asp.mu.Unlock()

	if !valid {
		log.Printf("Signal filtered out: %+v", signal)
		return
	}

	// Prioritization logic placeholder (could use a priority queue)
	// For now, just sleep inversely proportional to priority to simulate
	time.Sleep(time.Duration(10-priority) * time.Millisecond)

	// Generate order from signal
	order, err := asp.generateOrder(signal)
	if err != nil {
		log.Printf("Order generation failed: %v", err)
		return
	}

	// Run routing plugins
	asp.mu.Lock()
	for _, plugin := range asp.routingPlugins {
		if err := plugin.Route(&order); err != nil {
			log.Printf("Routing plugin error: %v", err)
			asp.mu.Unlock()
			return
		}
	}
	asp.mu.Unlock()

	// Run risk plugins registered in executor
	for _, riskPlugin := range asp.executor.riskPlugins {
		if err := riskPlugin.BeforeOrder(&order); err != nil {
			log.Printf("Risk plugin blocked order: %v", err)
			return
		}
	}

	// Place order
	if err := asp.executor.PlaceOrder(order); err != nil {
		log.Printf("Order placement failed: %v", err)
	}
}

func (asp *AdvancedSignalProcessor) generateOrder(signal Signal) (Order, error) {
	order := Order{
		StrategyID: signal.StrategyID,
		Symbol:     "",          // default empty, try to extract below
		Side:       signal.Type, // buy/sell/hold
		Quantity:   1.0,         // default quantity
		Price:      0,
		Meta:       signal.Payload,
	}
	if symbol, ok := signal.Payload["symbol"].(string); ok {
		order.Symbol = symbol
	}
	if qty, ok := signal.Payload["quantity"].(float64); ok {
		order.Quantity = qty
	}
	if price, ok := signal.Payload["price"].(float64); ok {
		order.Price = price
	}
	return order, nil
}

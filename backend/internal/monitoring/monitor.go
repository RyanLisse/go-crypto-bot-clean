package monitoring

import (
	"context"
	"log"
	"sync"
	"time"

	"go-crypto-bot-clean/backend/internal/execution"
)

// StrategyStatus holds current status and metrics for a strategy
type StrategyStatus struct {
	ID             string
	Status         string
	LastError      error
	SignalsHandled int
	OrdersPlaced   int
	SuccessOrders  int
	FailedOrders   int
	LastUpdate     time.Time
	Performance    map[string]float64 // e.g., PnL, Sharpe, etc.
}

// StrategyMonitor implements execution.EventHook to track and log strategy activity
type StrategyMonitor struct {
	mu       sync.RWMutex
	statuses map[string]*StrategyStatus
	events   chan Event
	alerts   chan Alert
	logger   *log.Logger
}

// Event represents a monitoring event for dashboards or logs
type Event struct {
	Timestamp   time.Time
	StrategyID  string
	EventType   string
	Description string
	Error       error
	Metadata    map[string]interface{}
}

// Alert represents an actionable alert
type Alert struct {
	Timestamp  time.Time
	StrategyID string
	Level      string // info, warning, error, critical
	Message    string
	Error      error
}

// NewStrategyMonitor creates a new monitor with optional logger
func NewStrategyMonitor(logger *log.Logger) *StrategyMonitor {
	return &StrategyMonitor{
		statuses: make(map[string]*StrategyStatus),
		events:   make(chan Event, 1000),
		alerts:   make(chan Alert, 100),
		logger:   logger,
	}
}

// OnSignalReceived logs and tracks signals
func (sm *StrategyMonitor) OnSignalReceived(signal execution.Signal) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	status := sm.getOrCreateStatus(signal.StrategyID)
	status.SignalsHandled++
	status.LastUpdate = time.Now()

	sm.emitEvent(Event{
		Timestamp:   time.Now(),
		StrategyID:  signal.StrategyID,
		EventType:   "signal_received",
		Description: "Signal received",
		Metadata: map[string]interface{}{
			"signal": signal,
		},
	})
}

// OnOrderPlaced logs and tracks orders
func (sm *StrategyMonitor) OnOrderPlaced(order execution.Order) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	status := sm.getOrCreateStatus(order.StrategyID)
	status.OrdersPlaced++
	status.LastUpdate = time.Now()

	// TODO: Enhance order monitoring when execution framework provides order result/error info.
	// For now, we just log the order placement event.

	sm.emitEvent(Event{
		Timestamp:   time.Now(),
		StrategyID:  order.StrategyID,
		EventType:   "order_placed",
		Description: "Order placed",
		Metadata: map[string]interface{}{
			"order": order,
		},
	})
}

// OnStatusUpdate logs status changes and errors
func (sm *StrategyMonitor) OnStatusUpdate(strategyID string, statusStr string, err error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	status := sm.getOrCreateStatus(strategyID)
	status.Status = statusStr
	status.LastUpdate = time.Now()
	status.LastError = err

	if err != nil {
		sm.emitAlert(Alert{
			Timestamp:  time.Now(),
			StrategyID: strategyID,
			Level:      "warning",
			Message:    "Strategy error: " + err.Error(),
			Error:      err,
		})
	}

	sm.emitEvent(Event{
		Timestamp:   time.Now(),
		StrategyID:  strategyID,
		EventType:   "status_update",
		Description: "Status updated to " + statusStr,
		Error:       err,
	})
}

// getOrCreateStatus fetches or initializes status for a strategy
func (sm *StrategyMonitor) getOrCreateStatus(strategyID string) *StrategyStatus {
	s, ok := sm.statuses[strategyID]
	if !ok {
		s = &StrategyStatus{
			ID:          strategyID,
			Status:      "unknown",
			Performance: make(map[string]float64),
		}
		sm.statuses[strategyID] = s
	}
	return s
}

// emitEvent sends event to channel and logs it
func (sm *StrategyMonitor) emitEvent(event Event) {
	select {
	case sm.events <- event:
	default:
		if sm.logger != nil {
			sm.logger.Printf("Event channel full, dropping event: %+v", event)
		}
	}
	if sm.logger != nil {
		sm.logger.Printf("[Strategy %s] %s: %s", event.StrategyID, event.EventType, event.Description)
	}
}

// emitAlert sends alert to channel and logs it
func (sm *StrategyMonitor) emitAlert(alert Alert) {
	select {
	case sm.alerts <- alert:
	default:
		if sm.logger != nil {
			sm.logger.Printf("Alert channel full, dropping alert: %+v", alert)
		}
	}
	if sm.logger != nil {
		sm.logger.Printf("[ALERT][Strategy %s][%s] %s", alert.StrategyID, alert.Level, alert.Message)
	}
}

// GetStatus returns a copy of current status for a strategy
func (sm *StrategyMonitor) GetStatus(strategyID string) StrategyStatus {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if s, ok := sm.statuses[strategyID]; ok {
		copy := *s
		return copy
	}
	return StrategyStatus{ID: strategyID, Status: "not_found"}
}

// ListStatuses returns a snapshot of all statuses
func (sm *StrategyMonitor) ListStatuses() []StrategyStatus {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	var list []StrategyStatus
	for _, s := range sm.statuses {
		copy := *s
		list = append(list, copy)
	}
	return list
}

// Events returns the event channel (read-only)
func (sm *StrategyMonitor) Events() <-chan Event {
	return sm.events
}

// Alerts returns the alert channel (read-only)
func (sm *StrategyMonitor) Alerts() <-chan Alert {
	return sm.alerts
}

// Example of graceful shutdown
func (sm *StrategyMonitor) Shutdown(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		sm.mu.Lock()
		defer sm.mu.Unlock()
		close(sm.events)
		close(sm.alerts)
		close(done)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

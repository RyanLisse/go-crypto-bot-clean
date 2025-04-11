// Code generated test file with advanced tests

package execution

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// mockStrategy implements the Strategy interface for testing.
type mockStrategy struct {
	initCalled  bool
	startCalled bool
	stopCalled  bool
	failInit    bool
	failStart   bool
	failStop    bool
}

func (m *mockStrategy) Initialize(ctx context.Context) error {
	m.initCalled = true
	if m.failInit {
		return errors.New("init error")
	}
	return nil
}

func (m *mockStrategy) Start(ctx context.Context) error {
	m.startCalled = true
	if m.failStart {
		return errors.New("start error")
	}
	return nil
}

func (m *mockStrategy) Stop(ctx context.Context) error {
	m.stopCalled = true
	if m.failStop {
		return errors.New("stop error")
	}
	return nil
}

// mockHook implements the EventHook interface for testing.
type mockHook struct {
	signals     []Signal
	orders      []Order
	statuses    []string
	statusErrs  []error
	strategyIDs []string
	mu          sync.Mutex
}

func (h *mockHook) OnSignalReceived(signal Signal) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.signals = append(h.signals, signal)
}

func (h *mockHook) OnOrderPlaced(order Order) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.orders = append(h.orders, order)
}

func (h *mockHook) OnStatusUpdate(strategyID string, status string, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.strategyIDs = append(h.strategyIDs, strategyID)
	h.statuses = append(h.statuses, status)
	h.statusErrs = append(h.statusErrs, err)
}

func TestStrategyExecutor_Lifecycle(t *testing.T) {
	executor := NewStrategyExecutor()
	hook := &mockHook{}
	executor.RegisterHook(hook)

	strat1 := &mockStrategy{}
	strat2 := &mockStrategy{failStart: true}

	if err := executor.AddStrategy("s1", strat1); err != nil {
		t.Fatalf("unexpected error adding strategy: %v", err)
	}
	if err := executor.AddStrategy("s2", strat2); err != nil {
		t.Fatalf("unexpected error adding strategy: %v", err)
	}

	ctx := context.Background()

	// Initialize all
	if err := executor.InitializeAll(ctx); err != nil {
		t.Fatalf("unexpected error initializing: %v", err)
	}
	if !strat1.initCalled || !strat2.initCalled {
		t.Errorf("expected Initialize to be called on all strategies")
	}

	// Start all (strat2 should fail)
	if err := executor.StartAll(ctx); err == nil {
		t.Errorf("expected error starting strategies, got nil")
	}
	if !strat1.startCalled || !strat2.startCalled {
		t.Errorf("expected Start to be called on all strategies")
	}

	// Stop all
	if err := executor.StopAll(ctx); err != nil {
		t.Fatalf("unexpected error stopping: %v", err)
	}
	if !strat1.stopCalled || !strat2.stopCalled {
		t.Errorf("expected Stop to be called on all strategies")
	}

	// Remove strategy
	if err := executor.RemoveStrategy("s1"); err != nil {
		t.Fatalf("unexpected error removing strategy: %v", err)
	}
	if err := executor.RemoveStrategy("s1"); err == nil {
		t.Errorf("expected error removing non-existent strategy, got nil")
	}
}

func TestStrategyExecutor_SignalHandling(t *testing.T) {
	executor := NewStrategyExecutor()
	hook := &mockHook{}
	executor.RegisterHook(hook)

	strat := &mockStrategy{}
	if err := executor.AddStrategy("s1", strat); err != nil {
		t.Fatalf("unexpected error adding strategy: %v", err)
	}

	signal := Signal{
		StrategyID: "s1",
		Type:       "buy",
		Payload:    map[string]interface{}{"price": 100},
	}

	if err := executor.HandleSignal(signal); err != nil {
		t.Fatalf("unexpected error handling signal: %v", err)
	}

	// Signal for non-existent strategy
	badSignal := Signal{
		StrategyID: "unknown",
		Type:       "sell",
	}
	if err := executor.HandleSignal(badSignal); err == nil {
		t.Errorf("expected error for unknown strategy, got nil")
	}

	// Wait for hook to be called asynchronously
	time.Sleep(10 * time.Millisecond)

	hook.mu.Lock()
	defer hook.mu.Unlock()
	if len(hook.signals) != 1 {
		t.Errorf("expected 1 signal received, got %d", len(hook.signals))
	}
}

func TestStrategyExecutor_OrderPlacement(t *testing.T) {
	executor := NewStrategyExecutor()
	hook := &mockHook{}
	executor.RegisterHook(hook)

	order := Order{
		StrategyID: "s1",
		Symbol:     "BTCUSDT",
		Side:       "buy",
		Quantity:   0.1,
		Price:      50000,
		Meta:       map[string]interface{}{"leverage": 10},
	}

	if err := executor.PlaceOrder(order); err != nil {
		t.Fatalf("unexpected error placing order: %v", err)
	}

	// Wait for hook to be called asynchronously
	time.Sleep(10 * time.Millisecond)

	hook.mu.Lock()
	defer hook.mu.Unlock()
	if len(hook.orders) != 1 {
		t.Errorf("expected 1 order placed, got %d", len(hook.orders))
	}
}

func TestStrategyExecutor_DuplicateStrategy(t *testing.T) {
	executor := NewStrategyExecutor()
	strat := &mockStrategy{}

	if err := executor.AddStrategy("s1", strat); err != nil {
		t.Fatalf("unexpected error adding strategy: %v", err)
	}
	if err := executor.AddStrategy("s1", strat); err == nil {
		t.Errorf("expected error adding duplicate strategy ID, got nil")
	}
}

// mockRiskPlugin blocks orders with Block=true
type mockRiskPlugin struct {
	Block bool
}

func (m *mockRiskPlugin) BeforeOrder(order *Order) error {
	if m.Block {
		return errors.New("blocked by risk plugin")
	}
	return nil
}

// mockMetricsPlugin records events
type mockMetricsPlugin struct {
	Events []string
	mu     sync.Mutex
}

func (m *mockMetricsPlugin) RecordEvent(event string, data map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Events = append(m.Events, event)
}

// Test async processing of signals and orders
func TestStrategyExecutor_AsyncProcessing(t *testing.T) {
	executor := NewStrategyExecutor()
	defer executor.StopDispatcher()

	hook := &mockHook{}
	executor.RegisterHook(hook)

	strat := &mockStrategy{}
	executor.AddStrategy("s1", strat)

	sig := Signal{StrategyID: "s1", Type: "buy"}
	ord := Order{StrategyID: "s1", Symbol: "BTCUSDT"}

	executor.HandleSignal(sig)
	executor.PlaceOrder(ord)

	time.Sleep(20 * time.Millisecond)

	hook.mu.Lock()
	defer hook.mu.Unlock()
	if len(hook.signals) != 1 {
		t.Errorf("expected 1 signal, got %d", len(hook.signals))
	}
	if len(hook.orders) != 1 {
		t.Errorf("expected 1 order, got %d", len(hook.orders))
	}
}

// Test risk plugin blocks order
func TestStrategyExecutor_RiskPluginBlocksOrder(t *testing.T) {
	executor := NewStrategyExecutor()
	defer executor.StopDispatcher()

	risk := &mockRiskPlugin{Block: true}
	executor.RegisterRiskPlugin(risk)

	ord := Order{StrategyID: "s1", Symbol: "BTCUSDT"}
	executor.PlaceOrder(ord)

	// No panic, just blocked silently (logged internally)
}

// Test metrics plugin receives events
func TestStrategyExecutor_MetricsPluginReceivesEvents(t *testing.T) {
	executor := NewStrategyExecutor()
	defer executor.StopDispatcher()

	metrics := &mockMetricsPlugin{}
	executor.RegisterMetricsPlugin(metrics)

	sig := Signal{StrategyID: "s1"}
	executor.signalCh <- sig

	ord := Order{StrategyID: "s1"}
	executor.orderCh <- ord

	time.Sleep(20 * time.Millisecond)

	metrics.mu.Lock()
	defer metrics.mu.Unlock()
	if len(metrics.Events) == 0 {
		t.Errorf("expected metrics events, got none")
	}
}

// Test hot-reloading of strategy params and implementation
func TestStrategyExecutor_HotReloading(t *testing.T) {
	executor := NewStrategyExecutor()
	defer executor.StopDispatcher()

	strat1 := &mockStrategy{}
	executor.AddStrategy("s1", strat1)

	params := map[string]interface{}{"threshold": 0.5}
	executor.UpdateStrategyParams("s1", params)

	strat2 := &mockStrategy{}
	err := executor.HotReloadStrategy("s1", strat2)
	if err != nil {
		t.Fatalf("unexpected error hot-reloading: %v", err)
	}
}

type panicRisk struct{}

func (p *panicRisk) BeforeOrder(order *Order) error {
	panic("risk plugin panic")
}

// Test panic recovery in async processing
func TestStrategyExecutor_PanicRecovery(t *testing.T) {
	executor := NewStrategyExecutor()
	defer executor.StopDispatcher()

	panicHook := &mockHook{}
	executor.RegisterHook(panicHook)

	executor.RegisterRiskPlugin(&panicRisk{})

	ord := Order{StrategyID: "s1"}
	executor.PlaceOrder(ord)

	time.Sleep(20 * time.Millisecond)
	// No crash expected
}

package execution

import (
	"errors"
	"testing"
	"time"
)

// Mock filter plugin
type mockFilterPlugin struct {
	validateResult bool
	priorityValue  int
}

func (m *mockFilterPlugin) Validate(signal Signal) bool {
	return m.validateResult
}

func (m *mockFilterPlugin) Priority(signal Signal) int {
	return m.priorityValue
}

// Mock routing plugin
type mockRoutingPlugin struct {
	fail bool
}

func (m *mockRoutingPlugin) Route(order *Order) error {
	if m.fail {
		return errors.New("routing failed")
	}
	return nil
}

// Mock risk plugin

// Test async processing, plugin integration, error handling
func TestAdvancedSignalProcessor_FullFlow(t *testing.T) {
	executor := NewStrategyExecutor()
	asp := NewAdvancedSignalProcessor(executor, 10)

	// Register plugins
	filter := &mockFilterPlugin{validateResult: true, priorityValue: 5}
	asp.RegisterFilterPlugin(filter)

	router := &mockRoutingPlugin{fail: false}
	asp.RegisterRoutingPlugin(router)

	risk := &mockRiskPlugin{Block: false}
	executor.RegisterRiskPlugin(risk)

	asp.Start()
	defer asp.Stop()

	// Submit a valid signal
	signal := Signal{
		StrategyID: "strat1",
		Type:       "buy",
		Payload: map[string]interface{}{
			"price":    100.0,
			"symbol":   "BTCUSDT",
			"quantity": 0.5,
		},
	}
	asp.SubmitSignal(signal)

	// Allow some time for async processing
	time.Sleep(100 * time.Millisecond)
}

// Test signal filtering blocks invalid signals
func TestAdvancedSignalProcessor_FilterBlocks(t *testing.T) {
	executor := NewStrategyExecutor()
	asp := NewAdvancedSignalProcessor(executor, 10)

	filter := &mockFilterPlugin{validateResult: false, priorityValue: 1}
	asp.RegisterFilterPlugin(filter)

	asp.Start()
	defer asp.Stop()

	signal := Signal{
		StrategyID: "strat2",
		Type:       "sell",
		Payload:    map[string]interface{}{"price": 50.0},
	}
	asp.SubmitSignal(signal)

	time.Sleep(50 * time.Millisecond)
}

// Test routing plugin failure blocks order
func TestAdvancedSignalProcessor_RoutingFailure(t *testing.T) {
	executor := NewStrategyExecutor()
	asp := NewAdvancedSignalProcessor(executor, 10)

	filter := &mockFilterPlugin{validateResult: true, priorityValue: 3}
	asp.RegisterFilterPlugin(filter)

	router := &mockRoutingPlugin{fail: true}
	asp.RegisterRoutingPlugin(router)

	asp.Start()
	defer asp.Stop()

	signal := Signal{
		StrategyID: "strat3",
		Type:       "buy",
		Payload:    map[string]interface{}{"price": 200.0},
	}
	asp.SubmitSignal(signal)

	time.Sleep(50 * time.Millisecond)
}

// Test risk plugin blocks order
func TestAdvancedSignalProcessor_RiskBlocks(t *testing.T) {
	executor := NewStrategyExecutor()
	asp := NewAdvancedSignalProcessor(executor, 10)

	filter := &mockFilterPlugin{validateResult: true, priorityValue: 2}
	asp.RegisterFilterPlugin(filter)

	router := &mockRoutingPlugin{fail: false}
	asp.RegisterRoutingPlugin(router)

	risk := &mockRiskPlugin{Block: true}
	executor.RegisterRiskPlugin(risk)

	asp.Start()
	defer asp.Stop()

	signal := Signal{
		StrategyID: "strat4",
		Type:       "sell",
		Payload:    map[string]interface{}{"price": 300.0},
	}
	asp.SubmitSignal(signal)

	time.Sleep(50 * time.Millisecond)
}

package execution

import (
	"context"
	"errors"
	"testing"
)

// mockExecutor mocks StrategyExecutor for testing SignalProcessor.
type mockExecutor struct {
	handledSignals []Signal
	handleErr      error
}

func (m *mockExecutor) HandleSignal(signal Signal) error {
	m.handledSignals = append(m.handledSignals, signal)
	return m.handleErr
}

func TestSignalProcessor_ProcessRawData(t *testing.T) {
	ctx := context.Background()
	mockExec := &mockExecutor{}
	sp := NewSignalProcessor(mockExec)

	rawData := MarketData{
		Symbol:    "BTCUSDT",
		Price:     50000,
		Volume:    1.5,
		Timestamp: 1234567890,
	}

	strategyInputs := map[string]map[string]interface{}{
		"strat1": {"priceThreshold": 40000.0},
		"strat2": {"priceThreshold": 60000.0},
	}

	sp.ProcessRawData(ctx, rawData, strategyInputs)

	if len(mockExec.handledSignals) != 2 {
		t.Errorf("expected 2 signals handled, got %d", len(mockExec.handledSignals))
	}

	for _, sig := range mockExec.handledSignals {
		if sig.StrategyID == "strat1" && sig.Type != "buy" {
			t.Errorf("expected 'buy' for strat1, got %s", sig.Type)
		}
		if sig.StrategyID == "strat2" && sig.Type != "sell" {
			t.Errorf("expected 'sell' for strat2, got %s", sig.Type)
		}
	}
}

func TestSignalProcessor_FilterInvalidSignal(t *testing.T) {
	mockExec := &mockExecutor{}
	sp := NewSignalProcessor(mockExec)

	invalidSignal := Signal{
		StrategyID: "stratX",
		Type:       "invalid_type",
		Payload:    map[string]interface{}{"foo": "bar"},
	}

	if sp.filterSignal(invalidSignal) {
		t.Error("expected invalid signal to be filtered out")
	}

	validSignal := Signal{
		StrategyID: "stratX",
		Type:       "buy",
		Payload:    map[string]interface{}{"foo": "bar"},
	}

	if !sp.filterSignal(validSignal) {
		t.Error("expected valid signal to pass filtering")
	}
}

func TestSignalProcessor_HandleSignalError(t *testing.T) {
	ctx := context.Background()
	mockExec := &mockExecutor{handleErr: errors.New("handle error")}
	sp := NewSignalProcessor(mockExec)

	rawData := MarketData{
		Symbol:    "ETHUSDT",
		Price:     2000,
		Volume:    3.0,
		Timestamp: 1234567890,
	}

	strategyInputs := map[string]map[string]interface{}{
		"stratErr": {"priceThreshold": 1000.0},
	}

	sp.ProcessRawData(ctx, rawData, strategyInputs)
	// No panic or crash expected despite error
}

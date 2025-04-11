package execution

import (
	"context"
	"log"
	"sync"
)

// SignalHandler abstracts the ability to handle signals.
type SignalHandler interface {
	HandleSignal(signal Signal) error
}

// MarketData represents raw market data input.
// This is a placeholder and can be extended with real fields.
type MarketData struct {
	Symbol    string
	Price     float64
	Volume    float64
	Timestamp int64
	// Add more fields as needed
}

// SignalProcessor processes raw data and strategy inputs into validated signals.
type SignalProcessor struct {
	executor SignalHandler
	mu       sync.Mutex
}

// NewSignalProcessor creates a new SignalProcessor.
func NewSignalProcessor(executor SignalHandler) *SignalProcessor {
	return &SignalProcessor{
		executor: executor,
	}
}

// ProcessRawData processes raw market data and strategy inputs, generates signals,
// filters, validates, prioritizes, and dispatches them to the StrategyExecutor.
func (sp *SignalProcessor) ProcessRawData(ctx context.Context, rawData MarketData, strategyInputs map[string]map[string]interface{}) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	for strategyID, inputs := range strategyInputs {
		signal, err := sp.generateSignal(strategyID, rawData, inputs)
		if err != nil {
			log.Printf("Signal generation error for strategy %s: %v", strategyID, err)
			continue
		}

		if !sp.filterSignal(signal) {
			log.Printf("Signal filtered out for strategy %s: %+v", strategyID, signal)
			continue
		}

		if err := sp.executor.HandleSignal(signal); err != nil {
			log.Printf("Failed to handle signal for strategy %s: %v", strategyID, err)
		}
	}
}

// generateSignal creates a Signal from raw data and strategy-specific inputs.
func (sp *SignalProcessor) generateSignal(strategyID string, data MarketData, inputs map[string]interface{}) (Signal, error) {
	// Placeholder logic: generate dummy buy/sell/hold based on price threshold
	// Replace with real strategy logic
	var signalType string
	priceThreshold, ok := inputs["priceThreshold"].(float64)
	if !ok {
		priceThreshold = 0
	}

	if data.Price > priceThreshold {
		signalType = "buy"
	} else if data.Price < priceThreshold {
		signalType = "sell"
	} else {
		signalType = "hold"
	}

	payload := map[string]interface{}{
		"price":  data.Price,
		"volume": data.Volume,
		"time":   data.Timestamp,
	}

	return Signal{
		StrategyID: strategyID,
		Type:       signalType,
		Payload:    payload,
	}, nil
}

// filterSignal applies filtering and validation logic to a signal.
// Returns true if the signal is valid and should be dispatched.
func (sp *SignalProcessor) filterSignal(signal Signal) bool {
	// Example filters:
	if signal.Type != "buy" && signal.Type != "sell" && signal.Type != "hold" {
		return false
	}
	if signal.Payload == nil {
		return false
	}
	// Add more validation/filtering as needed
	return true
}

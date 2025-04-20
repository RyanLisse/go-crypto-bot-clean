package usecase_test

import (
	"testing"
	"time"

	usecase "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"

	"github.com/stretchr/testify/assert"
)

// This function simulates the logic from SniperShotService for testing
func simulateConditionCheck(t *testing.T, condition usecase.TriggerCondition, price float64) bool {
	// Simulate the check logic from the SniperShotService.ExecuteSniper method
	triggered := false

	targetPrice := condition.TargetPrice
	operator := condition.Operator
	bufferPct := condition.PriceBufferPct

	// Apply price buffer if set
	if bufferPct > 0 {
		bufferAmount := targetPrice * bufferPct

		// Adjust target price based on operator
		switch operator {
		case ">", ">=":
			targetPrice -= bufferAmount // Lower the threshold for "above" triggers
		case "<", "<=":
			targetPrice += bufferAmount // Raise the threshold for "below" triggers
		}
	}

	// Check if condition is met based on operator
	switch operator {
	case ">":
		triggered = price > targetPrice
	case ">=":
		triggered = price >= targetPrice
	case "<":
		triggered = price < targetPrice
	case "<=":
		triggered = price <= targetPrice
	case "==":
		triggered = price == targetPrice
	}

	// If triggered, execute callbacks
	if triggered && condition.Callbacks != nil {
		for _, callback := range condition.Callbacks {
			callback(price)
		}
	}

	return triggered
}

func TestNewTriggerCondition(t *testing.T) {
	price := 50000.0
	operator := ">="

	condition := usecase.NewTriggerCondition(price, operator) // Prefix

	assert.NotNil(t, condition)
	assert.Equal(t, price, condition.TargetPrice)
	assert.Equal(t, usecase.ConditionOperator(operator), condition.Operator) // Prefix enum cast
	assert.Equal(t, 30*time.Second, condition.CheckInterval)                 // Default interval
	assert.Equal(t, 600*time.Second, condition.Timeout)                      // Default timeout
	assert.Equal(t, 0.0, condition.PriceBuffer)                              // Default buffer
	assert.Nil(t, condition.SuccessCallback)                                 // Default callback
	assert.Nil(t, condition.TimeoutCallback)                                 // Default callback
}

func TestTriggerConditionHelpers(t *testing.T) {
	t.Run("Simulates basic price check", func(t *testing.T) {
		condition := usecase.NewTriggerCondition(50000.0, ">")

		// Price below target
		assert.False(t, simulateConditionCheck(t, *condition, 49000.0))

		// Price above target
		assert.True(t, simulateConditionCheck(t, *condition, 51000.0))
	})

	t.Run("Handles callbacks during simulation", func(t *testing.T) {
		callbackTriggered := false
		callbackPrice := 0.0

		condition := usecase.NewTriggerCondition(50000.0, ">")
		condition.AddCallback(func(price float64) {
			callbackTriggered = true
			callbackPrice = price
		})

		// Trigger the condition
		simulateConditionCheck(t, *condition, 51000.0)

		assert.True(t, callbackTriggered)
		assert.Equal(t, 51000.0, callbackPrice)
	})
}

func TestTriggerCondition_Evaluate(t *testing.T) {
	tests := []struct {
		name          string
		condition     *usecase.TriggerCondition // Prefix
		currentPrice  float64
		expectedMet   bool
		expectedError bool
	}{
		{
			name:         "Greater Than - Met",
			condition:    usecase.NewTriggerCondition(50000, ">"), // Prefix
			currentPrice: 51000,
			expectedMet:  true,
		},
		// ... other test cases ...
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			met, err := tc.condition.Evaluate(tc.currentPrice)
			assert.Equal(t, tc.expectedMet, met)
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

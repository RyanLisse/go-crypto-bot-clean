package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTriggerCondition(t *testing.T) {
	t.Run("Creates with default values", func(t *testing.T) {
		condition := NewTriggerCondition(50000.0, ">")

		assert.Equal(t, 50000.0, condition.TargetPrice)
		assert.Equal(t, ">", condition.Operator)
		assert.Equal(t, 30, condition.MaxTimeoutSecs)   // default 30 seconds timeout
		assert.Equal(t, 0.0, condition.PriceBufferPct)  // default no buffer
		assert.Equal(t, 500, condition.CheckIntervalMs) // default 500ms check interval
		assert.Nil(t, condition.Callbacks)
	})

	t.Run("Sets timeout correctly", func(t *testing.T) {
		condition := NewTriggerCondition(50000.0, ">").WithTimeout(10)
		assert.Equal(t, 10, condition.MaxTimeoutSecs)
	})

	t.Run("Sets price buffer correctly", func(t *testing.T) {
		condition := NewTriggerCondition(50000.0, ">").WithPriceBuffer(0.01)
		assert.Equal(t, 0.01, condition.PriceBufferPct)
	})

	t.Run("Sets check interval correctly", func(t *testing.T) {
		condition := NewTriggerCondition(50000.0, ">").WithCheckInterval(100)
		assert.Equal(t, 100, condition.CheckIntervalMs)
	})

	t.Run("Adds callback function", func(t *testing.T) {
		condition := NewTriggerCondition(50000.0, ">")
		called := false

		condition.AddCallback(func(price float64) {
			called = true
		})

		assert.NotNil(t, condition.Callbacks)
		assert.Equal(t, 1, len(condition.Callbacks))

		// Call the callback
		condition.Callbacks[0](50001.0)
		assert.True(t, called)
	})

	t.Run("Chains callback with WithCallback", func(t *testing.T) {
		callCount := 0

		condition := NewTriggerCondition(50000.0, ">").
			WithCallback(func(price float64) {
				callCount++
			}).
			WithCallback(func(price float64) {
				callCount++
			})

		assert.NotNil(t, condition.Callbacks)
		assert.Equal(t, 2, len(condition.Callbacks))

		// Call the callbacks
		for _, cb := range condition.Callbacks {
			cb(50001.0)
		}
		assert.Equal(t, 2, callCount)
	})

	t.Run("Chain multiple modifiers", func(t *testing.T) {
		callCount := 0

		condition := NewTriggerCondition(50000.0, ">").
			WithTimeout(10).
			WithPriceBuffer(0.01).
			WithCheckInterval(100).
			WithCallback(func(price float64) {
				callCount++
			})

		assert.Equal(t, 50000.0, condition.TargetPrice)
		assert.Equal(t, ">", condition.Operator)
		assert.Equal(t, 10, condition.MaxTimeoutSecs)
		assert.Equal(t, 0.01, condition.PriceBufferPct)
		assert.Equal(t, 100, condition.CheckIntervalMs)
		assert.NotNil(t, condition.Callbacks)
		assert.Equal(t, 1, len(condition.Callbacks))

		// Call the callback
		condition.Callbacks[0](50001.0)
		assert.Equal(t, 1, callCount)
	})
}

func TestTriggerConditionCallbacks(t *testing.T) {
	// Create a basic trigger condition
	condition := NewTriggerCondition(50000.0, ">")

	// Verify initial state
	assert.Empty(t, condition.Callbacks)

	// Create a callback tracker
	var callbackCalled bool
	var receivedPrice float64

	callback := func(price float64) {
		callbackCalled = true
		receivedPrice = price
	}

	// Test AddCallback method
	condition.AddCallback(callback)
	assert.Len(t, condition.Callbacks, 1)

	// Call the callback directly to test
	testPrice := 51000.0
	condition.Callbacks[0](testPrice)

	// Verify callback was called with correct price
	assert.True(t, callbackCalled)
	assert.Equal(t, testPrice, receivedPrice)

	// Test WithCallback method (creates a new instance)
	callbackCalled = false // Reset
	receivedPrice = 0.0    // Reset

	secondCallback := func(price float64) {
		callbackCalled = true
		receivedPrice = price * 2 // Multiply by 2 to differentiate
	}

	newCondition := condition.WithCallback(secondCallback)

	// Verify original callbacks are still there
	assert.Len(t, condition.Callbacks, 1)

	// Verify new condition has both callbacks
	assert.Len(t, newCondition.Callbacks, 2)

	// Test second callback
	testPrice = 52000.0
	newCondition.Callbacks[1](testPrice)

	// Verify second callback was called with correct price
	assert.True(t, callbackCalled)
	assert.Equal(t, testPrice*2, receivedPrice) // Multiplied by 2
}

func TestTriggerConditionWithMethods(t *testing.T) {
	// Test fluent API
	condition := NewTriggerCondition(49000.0, "<").
		WithTimeout(60).
		WithPriceBuffer(0.02).
		WithCheckInterval(1000)

	// Verify values
	assert.Equal(t, 49000.0, condition.TargetPrice)
	assert.Equal(t, "<", condition.Operator)
	assert.Equal(t, 60, condition.MaxTimeoutSecs)
	assert.Equal(t, 0.02, condition.PriceBufferPct)
	assert.Equal(t, 1000, condition.CheckIntervalMs)
}

func TestTriggerConditionWithComplexPriceMovement(t *testing.T) {
	t.Run("Simulates price movement and triggers only after threshold crossed", func(t *testing.T) {
		// Create a trigger condition that activates when price goes above 50000
		targetPrice := 50000.0
		condition := NewTriggerCondition(targetPrice, ">")

		// Track when the trigger was activated and at what price
		triggered := false
		triggerPrice := 0.0

		// Add a callback that marks the trigger as activated
		condition.AddCallback(func(price float64) {
			triggered = true
			triggerPrice = price
		})

		// Simulate a sequence of price movements
		prices := []float64{
			48000.0, // Below target
			49000.0, // Below target
			49500.0, // Below target
			49900.0, // Below target
			50100.0, // Above target - should trigger
		}

		// Check each price against our condition's criteria
		for i, price := range prices {
			isTriggered := simulateConditionCheck(t, *condition, price) // dereferencing pointer

			if price <= targetPrice {
				assert.False(t, isTriggered, "Should not trigger at price %v (index: %d)", price, i)
			} else {
				assert.True(t, isTriggered, "Should trigger at price %v (index: %d)", price, i)
			}
		}

		// Verify final state
		assert.True(t, triggered, "Condition should have been triggered")
		assert.Equal(t, 50100.0, triggerPrice, "Trigger price should be the crossing point")
	})

	t.Run("Handles the 'Less Than' operator correctly", func(t *testing.T) {
		// Create a trigger condition that activates when price goes below 50000
		targetPrice := 50000.0
		condition := NewTriggerCondition(targetPrice, "<")

		// Track when the trigger was activated and at what price
		triggered := false
		triggerPrice := 0.0

		// Add a callback that marks the trigger as activated
		condition.AddCallback(func(price float64) {
			triggered = true
			triggerPrice = price
		})

		// Simulate a sequence of price movements (now decreasing)
		prices := []float64{
			52000.0, // Above target
			51000.0, // Above target
			50500.0, // Above target
			50100.0, // Above target
			49900.0, // Below target - should trigger
		}

		// Check each price against our condition's criteria
		for i, price := range prices {
			isTriggered := simulateConditionCheck(t, *condition, price) // dereferencing pointer

			if price >= targetPrice {
				assert.False(t, isTriggered, "Should not trigger at price %v (index: %d)", price, i)
			} else {
				assert.True(t, isTriggered, "Should trigger at price %v (index: %d)", price, i)
			}
		}

		// Verify final state
		assert.True(t, triggered, "Condition should have been triggered")
		assert.Equal(t, 49900.0, triggerPrice, "Trigger price should be the crossing point")
	})

	t.Run("Handles price buffer correctly", func(t *testing.T) {
		// Create a trigger condition with 1% buffer
		targetPrice := 50000.0
		bufferPct := 0.01 // 1%
		condition := NewTriggerCondition(targetPrice, ">=").WithPriceBuffer(bufferPct)

		// Calculate buffer thresholds
		bufferAmount := targetPrice * bufferPct
		effectiveThreshold := targetPrice - bufferAmount // 49500 should trigger with 1% buffer

		// Track when the trigger was activated
		triggered := false

		// Add a callback that marks the trigger as activated
		condition.AddCallback(func(price float64) {
			triggered = true
		})

		// Price just below the effective threshold - should not trigger
		belowThresholdPrice := effectiveThreshold - 10                              // 49490
		assert.False(t, simulateConditionCheck(t, *condition, belowThresholdPrice), // dereferencing pointer
			"Should not trigger at %v (below effective threshold %v)",
			belowThresholdPrice, effectiveThreshold)
		assert.False(t, triggered)

		// Price exactly at the effective threshold - should trigger
		exactThresholdPrice := effectiveThreshold
		assert.True(t, simulateConditionCheck(t, *condition, exactThresholdPrice), // dereferencing pointer
			"Should trigger at %v (exactly at effective threshold)",
			exactThresholdPrice)
		assert.True(t, triggered)
	})
}

// Helper function to simulate checking a price against a trigger condition
func simulateConditionCheck(t *testing.T, condition TriggerCondition, price float64) bool {
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

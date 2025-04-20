package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			isTriggered := simulateConditionCheck(t, *condition, price)

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
			isTriggered := simulateConditionCheck(t, *condition, price)

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
		belowThresholdPrice := effectiveThreshold - 10 // 49490
		assert.False(t, simulateConditionCheck(t, *condition, belowThresholdPrice),
			"Should not trigger at %v (below effective threshold %v)",
			belowThresholdPrice, effectiveThreshold)
		assert.False(t, triggered)

		// Price exactly at the effective threshold - should trigger
		exactThresholdPrice := effectiveThreshold
		assert.True(t, simulateConditionCheck(t, *condition, exactThresholdPrice),
			"Should trigger at %v (exactly at effective threshold)",
			exactThresholdPrice)
		assert.True(t, triggered)
	})

	t.Run("Multiple callbacks are all executed", func(t *testing.T) {
		// Create a trigger condition with multiple callbacks
		condition := NewTriggerCondition(50000.0, ">")

		// Track callback executions
		callback1Called := false
		callback2Called := false
		callback3Called := false

		// Add multiple callbacks
		condition.AddCallback(func(price float64) {
			callback1Called = true
		})

		condition.AddCallback(func(price float64) {
			callback2Called = true
		})

		condition.AddCallback(func(price float64) {
			callback3Called = true
		})

		// Trigger the condition
		simulateConditionCheck(t, *condition, 51000.0)

		// Verify all callbacks were executed
		assert.True(t, callback1Called, "First callback should have been called")
		assert.True(t, callback2Called, "Second callback should have been called")
		assert.True(t, callback3Called, "Third callback should have been called")
	})
}

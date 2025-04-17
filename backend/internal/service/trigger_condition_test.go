package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

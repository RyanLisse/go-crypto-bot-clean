package service

// TriggerCondition represents a condition that must be met before executing a trade.
type TriggerCondition struct {
	TargetPrice     float64               `json:"target_price"`      // The target price to trigger the trade
	Operator        string                `json:"operator"`          // The operator for comparison (e.g., ">", "<", ">=", "<=", "==")
	MaxTimeoutSecs  int                   `json:"max_timeout_secs"`  // Maximum time to wait for the trigger condition (0 = no timeout)
	PriceBufferPct  float64               `json:"price_buffer_pct"`  // Allow price slippage by percentage (0.01 = 1%)
	CheckIntervalMs int                   `json:"check_interval_ms"` // Time interval in ms for price checking
	Callbacks       []func(price float64) `json:"-"`                 // Callback functions to execute when condition is met
}

// AddCallback adds a callback function to be executed when the trigger condition is met
func (t *TriggerCondition) AddCallback(callback func(price float64)) {
	if t.Callbacks == nil {
		t.Callbacks = make([]func(price float64), 0)
	}
	t.Callbacks = append(t.Callbacks, callback)
}

// WithCallback creates a new TriggerCondition with the provided callback function
func (t TriggerCondition) WithCallback(callback func(price float64)) TriggerCondition {
	newCondition := t
	newCondition.AddCallback(callback)
	return newCondition
}

// NewTriggerCondition creates a new TriggerCondition with the specified parameters
func NewTriggerCondition(targetPrice float64, operator string) *TriggerCondition {
	return &TriggerCondition{
		TargetPrice:     targetPrice,
		Operator:        operator,
		MaxTimeoutSecs:  30,  // Default 30 seconds timeout
		PriceBufferPct:  0.0, // No price buffer by default
		CheckIntervalMs: 500, // Check every 500ms by default
		Callbacks:       nil,
	}
}

// WithTimeout sets the maximum timeout in seconds for the condition
func (t TriggerCondition) WithTimeout(timeoutSecs int) TriggerCondition {
	t.MaxTimeoutSecs = timeoutSecs
	return t
}

// WithPriceBuffer sets the price buffer percentage (0.01 = 1%)
func (t TriggerCondition) WithPriceBuffer(bufferPct float64) TriggerCondition {
	t.PriceBufferPct = bufferPct
	return t
}

// WithCheckInterval sets the check interval in milliseconds
func (t TriggerCondition) WithCheckInterval(checkIntervalMs int) TriggerCondition {
	t.CheckIntervalMs = checkIntervalMs
	return t
}

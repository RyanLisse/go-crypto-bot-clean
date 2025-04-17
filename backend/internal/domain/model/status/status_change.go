package status

import "time"

// Additional status constants for specific business domains
const (
	// StatusTrading indicates a component/symbol is available for trading
	StatusTrading Status = "trading"
)

// StatusChange represents a change in the status of a component
type StatusChange struct {
	// ID is a unique identifier for the status change
	ID string `json:"id"`
	// Component is the name of the component that changed status
	Component string `json:"component"`
	// OldStatus is the previous status
	OldStatus Status `json:"old_status"`
	// NewStatus is the current status
	NewStatus Status `json:"new_status"`
	// Timestamp is when the status change occurred
	Timestamp time.Time `json:"timestamp"`
	// Message provides additional details about the status change
	Message string `json:"message,omitempty"`
	// Metadata contains additional context for the status change
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// NewStatusChange creates a new status change event
func NewStatusChange(id, component string, oldStatus, newStatus Status, message string) *StatusChange {
	return &StatusChange{
		ID:        id,
		Component: component,
		OldStatus: oldStatus,
		NewStatus: newStatus,
		Timestamp: time.Now(),
		Message:   message,
		Metadata:  make(map[string]interface{}),
	}
}

// AddMetadata adds metadata to the status change
func (sc *StatusChange) AddMetadata(key string, value interface{}) {
	sc.Metadata[key] = value
}

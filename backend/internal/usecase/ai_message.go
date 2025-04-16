package usecase

import (
	"time"
)

// AIMessage represents a message in an AI conversation
type AIMessage struct {
	ID             string                 `json:"id"`
	ConversationID string                 `json:"conversation_id"`
	Role           string                 `json:"role"`
	Content        string                 `json:"content"`
	Timestamp      time.Time              `json:"timestamp"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

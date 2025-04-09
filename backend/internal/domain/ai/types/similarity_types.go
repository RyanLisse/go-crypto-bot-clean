package types

// SimilarMessage represents a similar message with its similarity score
type SimilarMessage struct {
	ConversationID string                 `json:"conversation_id"`
	MessageID      string                 `json:"message_id"`
	Content        string                 `json:"content"`
	Similarity     float64                `json:"similarity"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

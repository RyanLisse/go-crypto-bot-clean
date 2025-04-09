package handlers

import (
	"encoding/json"
	"net/http"

	"go-crypto-bot-clean/backend/internal/domain/ai/similarity"

	"github.com/go-chi/chi/v5"
)

// SimilarityHandler handles requests for similarity search
type SimilarityHandler struct {
	similaritySvc *similarity.Service
}

// NewSimilarityHandler creates a new similarity handler
func NewSimilarityHandler(similaritySvc *similarity.Service) *SimilarityHandler {
	return &SimilarityHandler{
		similaritySvc: similaritySvc,
	}
}

// IndexMessageRequest represents a request to index a message
type IndexMessageRequest struct {
	ConversationID string                 `json:"conversation_id"`
	MessageID      string                 `json:"message_id"`
	Content        string                 `json:"content"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// IndexMessage handles requests to index a message
func (h *SimilarityHandler) IndexMessage(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req IndexMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate request
	if req.ConversationID == "" {
		http.Error(w, "Missing conversation_id", http.StatusBadRequest)
		return
	}
	if req.MessageID == "" {
		http.Error(w, "Missing message_id", http.StatusBadRequest)
		return
	}
	if req.Content == "" {
		http.Error(w, "Missing content", http.StatusBadRequest)
		return
	}

	// Index message
	err := h.similaritySvc.IndexMessage(r.Context(), req.ConversationID, req.MessageID, req.Content, req.Metadata)
	if err != nil {
		http.Error(w, "Failed to index message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}

// FindSimilarMessagesRequest represents a request to find similar messages
type FindSimilarMessagesRequest struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

// FindSimilarMessages handles requests to find similar messages
func (h *SimilarityHandler) FindSimilarMessages(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req FindSimilarMessagesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Query == "" {
		http.Error(w, "Missing query", http.StatusBadRequest)
		return
	}
	if req.Limit <= 0 {
		req.Limit = 10 // Default limit
	}

	// Find similar messages
	similarMessages, err := h.similaritySvc.FindSimilarMessages(r.Context(), req.Query, req.Limit)
	if err != nil {
		http.Error(w, "Failed to find similar messages: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"messages": similarMessages,
	})
}

// DeleteConversationMessages handles requests to delete all indexed messages for a conversation
func (h *SimilarityHandler) DeleteConversationMessages(w http.ResponseWriter, r *http.Request) {
	// Get conversation ID from URL
	conversationID := chi.URLParam(r, "conversationID")
	if conversationID == "" {
		http.Error(w, "Missing conversation_id", http.StatusBadRequest)
		return
	}

	// Delete conversation messages
	err := h.similaritySvc.DeleteConversationMessages(r.Context(), conversationID)
	if err != nil {
		http.Error(w, "Failed to delete conversation messages: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}

// RegisterRoutes registers the similarity routes
func (h *SimilarityHandler) RegisterRoutes(r chi.Router) {
	r.Post("/api/similarity/index", h.IndexMessage)
	r.Post("/api/similarity/search", h.FindSimilarMessages)
	r.Delete("/api/similarity/conversations/{conversationID}", h.DeleteConversationMessages)
}

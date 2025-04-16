package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/delivery/http/response"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestAIHandler_Chat_Authenticated(t *testing.T) {
	logger := zerolog.Nop()

	mockAIService := &usecase.MockAIService{}
	mockConvRepo := &usecase.MockConversationMemoryRepository{}
	mockEmbedRepo := &usecase.MockEmbeddingRepository{}
	useCase := usecase.NewAIUsecase(mockAIService, mockConvRepo, mockEmbedRepo, logger)
	h := NewAIHandler(useCase, &logger)

	r := chi.NewRouter()
	h.RegisterRoutes(r, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			ctx = contextWithUserID(ctx, "test-user")
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	})

	body := map[string]interface{}{
		"user_id":         "test-user",
		"message":         "Hello AI!",
		"session_id":      "sess-1",
		"trading_context": map[string]interface{}{"foo": "bar"},
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/ai/chat", bytes.NewReader(b))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp.Success)
}

func contextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, "user_id", userID)
}

// Add more tests for conversation history, pagination, and error cases as needed.

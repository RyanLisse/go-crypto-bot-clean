package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/delivery/http/handler"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCredential(t *testing.T) {
	// Create a mock use case
	mockUseCase := new(mocks.MockAPICredentialUseCase)

	// Create a logger
	logger := zerolog.New(zerolog.NewTestWriter(t))

	// Create the handler
	h := handler.NewAPICredentialHandler(mockUseCase, &logger)

	// Create a request with valid MEXC API credentials
	reqBody := map[string]interface{}{
		"exchange":  "mexc",
		"apiKey":    "mx1234567890abcdef1234567890abcdef",
		"apiSecret": "mx1234567890abcdef1234567890abcdef",
		"label":     "Test Credential",
	}
	reqBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/credentials", bytes.NewReader(reqBytes))

	// Set up the context with a user ID
	ctx := context.WithValue(req.Context(), "userID", "user1")
	req = req.WithContext(ctx)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Set up the mock use case to expect a call to CreateCredential
	mockUseCase.On("CreateCredential", mock.Anything, mock.AnythingOfType("*model.APICredential")).Return(nil)

	// Call the handler
	h.CreateCredential(w, req)

	// Assert the response
	assert.Equal(t, http.StatusCreated, w.Code)

	// Parse the response
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// Assert the response data
	assert.Equal(t, true, resp["success"])
	assert.NotNil(t, resp["data"])

	// Verify that the mock was called
	mockUseCase.AssertExpectations(t)
}

func TestListCredentials(t *testing.T) {
	// Create a mock use case
	mockUseCase := new(mocks.MockAPICredentialUseCase)

	// Create a logger
	logger := zerolog.New(zerolog.NewTestWriter(t))

	// Create the handler
	h := handler.NewAPICredentialHandler(mockUseCase, &logger)

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/credentials", nil)

	// Set up the context with a user ID
	ctx := context.WithValue(req.Context(), "userID", "user1")
	req = req.WithContext(ctx)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Set up the mock use case to expect a call to ListCredentials
	mockCredentials := []*model.APICredential{
		{
			ID:        "cred1",
			UserID:    "user1",
			Exchange:  "mexc",
			APIKey:    "apiKey1",
			APISecret: "apiSecret1",
			Label:     "Credential 1",
		},
		{
			ID:        "cred2",
			UserID:    "user1",
			Exchange:  "binance",
			APIKey:    "apiKey2",
			APISecret: "apiSecret2",
			Label:     "Credential 2",
		},
	}
	mockUseCase.On("ListCredentials", mock.Anything, "user1").Return(mockCredentials, nil)

	// Call the handler
	h.ListCredentials(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// Assert the response data
	assert.Equal(t, true, resp["success"])
	assert.NotNil(t, resp["data"])
	data, ok := resp["data"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, data, 2)

	// Verify that the mock was called
	mockUseCase.AssertExpectations(t)
}

// Unused helper function - can be activated when needed for route testing
// func setupRouter(h *handler.APICredentialHandler) *chi.Mux {
// 	r := chi.NewRouter()
// 	h.RegisterRoutes(r)
// 	return r
// }

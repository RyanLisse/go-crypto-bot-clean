package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/service"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

var _ service.AuthServiceInterface = (*MockAuthService)(nil)

func (m *MockAuthService) VerifyToken(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) GetUserFromToken(ctx context.Context, token string) (*model.User, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockAuthService) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockAuthService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockAuthService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func TestAuthController_VerifyToken(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create mocks
	mockAuthService := new(MockAuthService)

	// Create controller
	controller := NewAuthController(mockAuthService, &logger)

	// Create test user
	testUser := &model.User{
		ID:        "test-user-id",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create test roles
	testRoles := []string{"user", "admin"}

	// Test token
	testToken := "valid-token"

	t.Run("Valid Token", func(t *testing.T) {
		// Mock VerifyToken
		mockAuthService.On("VerifyToken", mock.Anything, testToken).Return(testUser.ID, nil).Once()

		// Mock GetUserByID
		mockAuthService.On("GetUserByID", mock.Anything, testUser.ID).Return(testUser, nil).Once()

		// Mock GetUserRoles
		mockAuthService.On("GetUserRoles", mock.Anything, testUser.ID).Return(testRoles, nil).Once()

		// Create request
		requestBody := `{"token":"valid-token"}`
		req := httptest.NewRequest("POST", "/auth/verify", strings.NewReader(requestBody))
		res := httptest.NewRecorder()

		// Call handler
		controller.VerifyToken(res, req)

		// Check response
		assert.Equal(t, http.StatusOK, res.Code)

		// Parse response body
		var response map[string]interface{}
		err := json.Unmarshal(res.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testUser.ID, response["id"])
		assert.Equal(t, testUser.Email, response["email"])
		assert.Equal(t, testUser.Name, response["name"])
		assert.Equal(t, []interface{}{"user", "admin"}, response["roles"])
	})

	t.Run("Invalid Token", func(t *testing.T) {
		// Mock VerifyToken to return an error
		mockAuthService.On("VerifyToken", mock.Anything, "invalid-token").Return("", assert.AnError).Once()

		// Create request
		requestBody := `{"token":"invalid-token"}`
		req := httptest.NewRequest("POST", "/auth/verify", strings.NewReader(requestBody))
		res := httptest.NewRecorder()

		// Call handler
		controller.VerifyToken(res, req)

		// Check response
		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		// Create request with invalid JSON
		requestBody := `{"token":}`
		req := httptest.NewRequest("POST", "/auth/verify", strings.NewReader(requestBody))
		res := httptest.NewRecorder()

		// Call handler
		controller.VerifyToken(res, req)

		// Check response
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	// Verify all mocks
	mockAuthService.AssertExpectations(t)
}

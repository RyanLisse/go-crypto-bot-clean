package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of the AuthService
type MockAuthService struct {
	mock.Mock
}

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

func TestEnhancedClerkMiddleware_NoAuthHeader(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create mock auth service
	mockAuthService := new(MockAuthService)

	// Create middleware
	middleware := NewEnhancedClerkMiddleware(mockAuthService, &logger)

	// Create a test handler that will be wrapped
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Should be called since we're not providing an auth header
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test passed"))
	})

	// Create a request without auth header
	req := httptest.NewRequest("GET", "/test", nil)
	res := httptest.NewRecorder()

	// Apply the middleware
	handler := middleware.Middleware()(testHandler)

	// Send the request
	handler.ServeHTTP(res, req)

	// Check the response
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "test passed", res.Body.String())

	// Verify mocks
	mockAuthService.AssertExpectations(t)
}

func TestEnhancedClerkMiddleware_WithValidToken(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create mock auth service
	mockAuthService := new(MockAuthService)

	// Create middleware
	middleware := NewEnhancedClerkMiddleware(mockAuthService, &logger)

	// Create test user
	testUser := &model.User{
		ID:    "test-user-id",
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Create test roles
	testRoles := []string{"user", "admin"}

	// Mock GetUserFromToken
	mockAuthService.On("GetUserFromToken", mock.Anything, "valid-token").Return(testUser, nil)

	// Mock GetUserRoles
	mockAuthService.On("GetUserRoles", mock.Anything, testUser.ID).Return(testRoles, nil)

	// Create a test handler that will be wrapped
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that user ID is in context
		userID, ok := GetUserIDFromContext(r.Context())
		assert.True(t, ok)
		assert.Equal(t, testUser.ID, userID)

		// Check that roles are in context
		roles, ok := r.Context().Value(RoleKey).([]string)
		assert.True(t, ok)
		assert.Equal(t, testRoles, roles)

		// Check that user is in context
		user, ok := GetUserFromContext(r.Context())
		assert.True(t, ok)
		assert.Equal(t, testUser, user)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test passed"))
	})

	// Create a request with auth header
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	res := httptest.NewRecorder()

	// Apply the middleware
	handler := middleware.Middleware()(testHandler)

	// Send the request
	handler.ServeHTTP(res, req)

	// Check the response
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "test passed", res.Body.String())

	// Verify mocks
	mockAuthService.AssertExpectations(t)
}

func TestEnhancedClerkMiddleware_WithInvalidToken(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create mock auth service
	mockAuthService := new(MockAuthService)

	// Create middleware
	middleware := NewEnhancedClerkMiddleware(mockAuthService, &logger)

	// Mock GetUserFromToken to return an error
	mockAuthService.On("GetUserFromToken", mock.Anything, "invalid-token").Return(nil, assert.AnError)

	// Create a test handler that should not be called
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with invalid token")
	})

	// Create a request with invalid auth header
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	res := httptest.NewRecorder()

	// Apply the middleware
	handler := middleware.Middleware()(testHandler)

	// Send the request
	handler.ServeHTTP(res, req)

	// Check that it returned unauthorized
	assert.Equal(t, http.StatusUnauthorized, res.Code)

	// Verify mocks
	mockAuthService.AssertExpectations(t)
}

func TestEnhancedClerkMiddleware_RequireAuthentication(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create mock auth service
	mockAuthService := new(MockAuthService)

	// Create middleware
	middleware := NewEnhancedClerkMiddleware(mockAuthService, &logger)

	// Create a test handler that will be wrapped
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test passed"))
	})

	t.Run("With Authentication", func(t *testing.T) {
		// Create a request with user ID in context
		req := httptest.NewRequest("GET", "/test", nil)
		ctx := context.WithValue(req.Context(), UserIDKey, "test-user-id")
		req = req.WithContext(ctx)
		res := httptest.NewRecorder()

		// Apply the middleware
		handler := middleware.RequireAuthentication(testHandler)

		// Send the request
		handler.ServeHTTP(res, req)

		// Check the response
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, "test passed", res.Body.String())
	})

	t.Run("Without Authentication", func(t *testing.T) {
		// Create a request without user ID in context
		req := httptest.NewRequest("GET", "/test", nil)
		res := httptest.NewRecorder()

		// Apply the middleware
		handler := middleware.RequireAuthentication(testHandler)

		// Send the request
		handler.ServeHTTP(res, req)

		// Check that it returned unauthorized
		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})
}

func TestEnhancedClerkMiddleware_RequireRole(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create mock auth service
	mockAuthService := new(MockAuthService)

	// Create middleware
	middleware := NewEnhancedClerkMiddleware(mockAuthService, &logger)

	// Create a test handler that will be wrapped
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test passed"))
	})

	t.Run("With Required Role", func(t *testing.T) {
		// Create a request with roles in context
		req := httptest.NewRequest("GET", "/test", nil)
		ctx := context.WithValue(req.Context(), RoleKey, []string{"user", "admin"})
		req = req.WithContext(ctx)
		res := httptest.NewRecorder()

		// Apply the middleware
		handler := middleware.RequireRole("admin")(testHandler)

		// Send the request
		handler.ServeHTTP(res, req)

		// Check the response
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, "test passed", res.Body.String())
	})

	t.Run("Without Required Role", func(t *testing.T) {
		// Create a request with roles in context
		req := httptest.NewRequest("GET", "/test", nil)
		ctx := context.WithValue(req.Context(), RoleKey, []string{"user"})
		req = req.WithContext(ctx)
		res := httptest.NewRecorder()

		// Apply the middleware
		handler := middleware.RequireRole("admin")(testHandler)

		// Send the request
		handler.ServeHTTP(res, req)

		// Check that it returned forbidden
		assert.Equal(t, http.StatusForbidden, res.Code)
	})

	t.Run("Without Roles", func(t *testing.T) {
		// Create a request without roles in context
		req := httptest.NewRequest("GET", "/test", nil)
		res := httptest.NewRecorder()

		// Apply the middleware
		handler := middleware.RequireRole("admin")(testHandler)

		// Send the request
		handler.ServeHTTP(res, req)

		// Check that it returned unauthorized
		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})
}

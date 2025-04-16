package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Define error constants
var (
	ErrInvalidToken = errors.New("invalid token")
)

// MockAuthService is a mock implementation of the AuthServiceInterface
type MockAuthService struct {
	mock.Mock
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

func (m *MockAuthService) VerifyToken(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
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

func TestAuthMiddleware_NoAuthHeader(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create a mock auth service
	mockAuthService := new(MockAuthService)

	// Create middleware
	middleware := NewAuthMiddleware(mockAuthService, &logger)

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
}

func TestAuthMiddleware_ValidAuthHeader(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create a mock auth service
	mockAuthService := new(MockAuthService)

	// Set up expectations
	mockUser := &model.User{
		ID:    "test-user-id",
		Email: "test@example.com",
		Name:  "Test User",
	}
	mockRoles := []string{"user", "admin"}

	mockAuthService.On("GetUserFromToken", mock.Anything, "valid-token").Return(mockUser, nil)
	mockAuthService.On("GetUserRoles", mock.Anything, "test-user-id").Return(mockRoles, nil)

	// Create middleware
	middleware := NewAuthMiddleware(mockAuthService, &logger)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that user context values are set
		userID, ok := r.Context().Value(UserIDKey{}).(string)
		assert.True(t, ok, "UserID should be set in context")
		assert.Equal(t, "test-user-id", userID)

		roles, ok := r.Context().Value(RolesKey{}).([]string)
		assert.True(t, ok, "Roles should be set in context")
		assert.Equal(t, mockRoles, roles)

		user, ok := r.Context().Value(UserKey{}).(*model.User)
		assert.True(t, ok, "User should be set in context")
		assert.Equal(t, mockUser, user)

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

	// Verify that the mock was called as expected
	mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create a mock auth service
	mockAuthService := new(MockAuthService)

	// Set up expectations
	mockAuthService.On("GetUserFromToken", mock.Anything, "invalid-token").Return(nil, ErrInvalidToken)

	// Create middleware
	middleware := NewAuthMiddleware(mockAuthService, &logger)

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

	// Verify that the mock was called as expected
	mockAuthService.AssertExpectations(t)
}

func TestRequireAuthentication_Authenticated(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create a mock auth service
	mockAuthService := new(MockAuthService)

	// Create middleware
	middleware := NewAuthMiddleware(mockAuthService, &logger)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test passed"))
	})

	// Create a request with user ID in context
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), UserIDKey{}, "test-user-id")
	req = req.WithContext(ctx)
	res := httptest.NewRecorder()

	// Apply the middleware
	handler := middleware.RequireAuthentication(testHandler)

	// Send the request
	handler.ServeHTTP(res, req)

	// Check the response
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "test passed", res.Body.String())
}

func TestRequireAuthentication_Unauthenticated(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create a mock auth service
	mockAuthService := new(MockAuthService)

	// Create middleware
	middleware := NewAuthMiddleware(mockAuthService, &logger)

	// Create a test handler that should not be called
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called without authentication")
	})

	// Create a request without user ID in context
	req := httptest.NewRequest("GET", "/test", nil)
	res := httptest.NewRecorder()

	// Apply the middleware
	handler := middleware.RequireAuthentication(testHandler)

	// Send the request
	handler.ServeHTTP(res, req)

	// Check that it returned unauthorized
	assert.Equal(t, http.StatusUnauthorized, res.Code)
}

func TestRequireRole_HasRole(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create a mock auth service
	mockAuthService := new(MockAuthService)

	// Create middleware
	middleware := NewAuthMiddleware(mockAuthService, &logger)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test passed"))
	})

	// Create a request with roles in context
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), RolesKey{}, []string{"user", "admin"})
	req = req.WithContext(ctx)
	res := httptest.NewRecorder()

	// Apply the middleware
	handler := middleware.RequireRole("admin")(testHandler)

	// Send the request
	handler.ServeHTTP(res, req)

	// Check the response
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "test passed", res.Body.String())
}

func TestRequireRole_DoesNotHaveRole(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create a mock auth service
	mockAuthService := new(MockAuthService)

	// Create middleware
	middleware := NewAuthMiddleware(mockAuthService, &logger)

	// Create a test handler that should not be called
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called without required role")
	})

	// Create a request with roles in context
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), RolesKey{}, []string{"user"})
	req = req.WithContext(ctx)
	res := httptest.NewRecorder()

	// Apply the middleware
	handler := middleware.RequireRole("admin")(testHandler)

	// Send the request
	handler.ServeHTTP(res, req)

	// Check that it returned forbidden
	assert.Equal(t, http.StatusForbidden, res.Code)
}

func TestRequireRole_NoRoles(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create a mock auth service
	mockAuthService := new(MockAuthService)

	// Create middleware
	middleware := NewAuthMiddleware(mockAuthService, &logger)

	// Create a test handler that should not be called
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called without roles")
	})

	// Create a request without roles in context
	req := httptest.NewRequest("GET", "/test", nil)
	res := httptest.NewRecorder()

	// Apply the middleware
	handler := middleware.RequireRole("admin")(testHandler)

	// Send the request
	handler.ServeHTTP(res, req)

	// Check that it returned unauthorized
	assert.Equal(t, http.StatusUnauthorized, res.Code)
}

func TestTestAuthMiddleware(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create middleware
	middleware := NewTestAuthMiddleware(&logger)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that user context values are set
		userID, ok := r.Context().Value(UserIDKey{}).(string)
		assert.True(t, ok, "UserID should be set in context")
		assert.Equal(t, "test_user_id", userID)

		roles, ok := r.Context().Value(RolesKey{}).([]string)
		assert.True(t, ok, "Roles should be set in context")
		assert.Equal(t, []string{"user", "admin"}, roles)

		user, ok := r.Context().Value(UserKey{}).(*model.User)
		assert.True(t, ok, "User should be set in context")
		assert.Equal(t, "test_user_id", user.ID)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test passed"))
	})

	// Create a request
	req := httptest.NewRequest("GET", "/test", nil)
	res := httptest.NewRecorder()

	// Apply the middleware
	handler := middleware.Middleware()(testHandler)

	// Send the request
	handler.ServeHTTP(res, req)

	// Check the response
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "test passed", res.Body.String())
}

func TestDisabledAuthMiddleware(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create middleware
	middleware := NewDisabledAuthMiddleware(&logger)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test passed"))
	})

	// Create a request
	req := httptest.NewRequest("GET", "/test", nil)
	res := httptest.NewRecorder()

	// Apply the middleware
	handler := middleware.Middleware()(testHandler)

	// Send the request
	handler.ServeHTTP(res, req)

	// Check the response
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "test passed", res.Body.String())
}

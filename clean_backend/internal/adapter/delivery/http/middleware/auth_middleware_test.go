package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of the AuthServiceInterface
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockAuthService) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockAuthService) GetUserFromToken(ctx context.Context, token string) (*model.User, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockAuthService) ValidateToken(ctx context.Context, token string) (bool, error) {
	args := m.Called(ctx, token)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthService) GenerateToken(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) RevokeToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func TestAuthFactory_CreateMiddleware(t *testing.T) {
	// Create a logger
	logger := zerolog.New(zerolog.NewTestWriter(t))

	// Create a mock auth service
	mockAuthService := new(MockAuthService)

	// Create a config
	cfg := &config.Config{
		Auth: config.AuthConfig{
			Provider:       "clerk",
			Disabled:       false,
			ClerkSecretKey: "test_key",
		},
	}

	// Create an auth factory
	factory := NewAuthFactory(mockAuthService, cfg, &logger)

	// Test creating different middleware types
	tests := []struct {
		name     string
		authType AuthType
		want     string
	}{
		{
			name:     "Clerk middleware",
			authType: AuthTypeClerk,
			want:     "*middleware.ClerkMiddleware",
		},
		{
			name:     "Test middleware",
			authType: AuthTypeTest,
			want:     "*middleware.TestMiddleware",
		},
		{
			name:     "Disabled middleware",
			authType: AuthTypeDisabled,
			want:     "*middleware.DisabledMiddleware",
		},
		{
			name:     "Default middleware",
			authType: "unknown",
			want:     "*middleware.ClerkMiddleware",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := factory.CreateMiddleware(tt.authType)
			assert.NotNil(t, middleware)
			// Check the type name
			typeName := fmt.Sprintf("%T", middleware)
			assert.Equal(t, tt.want, typeName)
		})
	}
}

func TestTestMiddleware(t *testing.T) {
	// Create a logger
	logger := zerolog.New(zerolog.NewTestWriter(t))

	// Create a test middleware
	middleware := NewTestMiddleware(&logger)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user ID is in context
		userID, ok := GetUserIDFromContext(r.Context())
		assert.True(t, ok)
		assert.Equal(t, "test_user_id", userID)

		// Check if roles are in context
		roles, ok := GetRolesFromContext(r.Context())
		assert.True(t, ok)
		assert.Contains(t, roles, "user")
		assert.Contains(t, roles, "admin")

		// Check if user is in context
		user, ok := GetUserFromContext(r.Context())
		assert.True(t, ok)
		assert.Equal(t, "test_user_id", user.ID)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "Test User", user.Name)

		w.WriteHeader(http.StatusOK)
	})

	// Create a request
	req := httptest.NewRequest("GET", "/", nil)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Apply the middleware
	handler := middleware.Middleware()(testHandler)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the response
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDisabledMiddleware(t *testing.T) {
	// Create a logger
	logger := zerolog.New(zerolog.NewTestWriter(t))

	// Create a disabled middleware
	middleware := NewDisabledMiddleware(&logger)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a request
	req := httptest.NewRequest("GET", "/", nil)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Apply the middleware
	handler := middleware.Middleware()(testHandler)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Test RequireAuthentication
	rr = httptest.NewRecorder()
	handler = middleware.RequireAuthentication(testHandler)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Test RequireRole
	rr = httptest.NewRecorder()
	handler = middleware.RequireRole("admin")(testHandler)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

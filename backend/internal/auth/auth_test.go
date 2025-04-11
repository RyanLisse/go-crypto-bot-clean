package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// local context key to fix compilation error
var sessionClaimsContextKey = &struct{}{}

func TestService_Authenticate(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func() *http.Request
		expectedError  error
		expectedUserID string
	}{
		{
			name: "successful authentication",
			setupMock: func() *http.Request {
				// Create a request with valid session claims
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				ctx := context.WithValue(r.Context(), sessionClaimsContextKey, map[string]interface{}{
					"sub": "user_123",
				})
				return r.WithContext(ctx)
			},
			expectedUserID: "user_123",
			expectedError:  nil,
		},
		{
			name: "missing session claims",
			setupMock: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/", nil)
			},
			expectedError: ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new service instance
			service := NewService("test_secret_key")

			// Setup the test request
			r := tt.setupMock()

			// Call Authenticate
			userData, err := service.Authenticate(r)

			// Check error
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, userData)
				return
			}

			// Check success case
			require.NoError(t, err)
			require.NotNil(t, userData)
			assert.Equal(t, tt.expectedUserID, userData.ID)
		})
	}
}

func TestService_RequireRole(t *testing.T) {
	tests := []struct {
		name          string
		role          string
		userRoles     []string
		expectedError bool
	}{
		{
			name:          "user has required role",
			role:          "admin",
			userRoles:     []string{"admin", "user"},
			expectedError: false,
		},
		{
			name:          "user does not have required role",
			role:          "admin",
			userRoles:     []string{"user"},
			expectedError: true,
		},
		{
			name:          "empty user roles",
			role:          "admin",
			userRoles:     []string{},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new service instance
			service := NewService("test_secret_key")

			// Create a test handler that will be wrapped
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create the middleware
			middleware := service.RequireRole(tt.role)
			handler := middleware(nextHandler)

			// Create a test request
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			// Add user data to context
			ctx := context.WithValue(r.Context(), UserDataKey, UserData{
				ID:    "test_user",
				Roles: tt.userRoles,
			})
			r = r.WithContext(ctx)

			// Call the handler
			handler.ServeHTTP(w, r)

			if tt.expectedError {
				assert.Equal(t, http.StatusForbidden, w.Code)
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
			}
		})
	}
}

func TestService_RequirePermission(t *testing.T) {
	tests := []struct {
		name          string
		permission    string
		userRoles     []string
		expectedError bool
	}{
		{
			name:          "user has required permission through role",
			permission:    "read:users",
			userRoles:     []string{"admin"},
			expectedError: false,
		},
		{
			name:          "user does not have required permission",
			permission:    "delete:users",
			userRoles:     []string{"user"},
			expectedError: true,
		},
		{
			name:          "empty user roles",
			permission:    "read:users",
			userRoles:     []string{},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new service instance
			service := NewService("test_secret_key")

			// Create a test handler that will be wrapped
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create the middleware
			middleware := service.RequirePermission(tt.permission)
			handler := middleware(nextHandler)

			// Create a test request
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			// Add user data to context
			ctx := context.WithValue(r.Context(), UserDataKey, UserData{
				ID:    "test_user",
				Roles: tt.userRoles,
			})
			r = r.WithContext(ctx)

			// Call the handler
			handler.ServeHTTP(w, r)

			if tt.expectedError {
				assert.Equal(t, http.StatusForbidden, w.Code)
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
			}
		})
	}
}

func TestDisabledService(t *testing.T) {
	service := NewDisabledService()

	t.Run("Authenticate returns error", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		userData, err := service.Authenticate(r)
		assert.Error(t, err)
		assert.Nil(t, userData)
		assert.Contains(t, err.Error(), "Authentication service is disabled")
	})

	t.Run("RequireRole returns error", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := service.RequireRole("admin")
		handler := middleware(nextHandler)

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, r)
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	})

	t.Run("RequirePermission returns error", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := service.RequirePermission("read:users")
		handler := middleware(nextHandler)

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, r)
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	})
}

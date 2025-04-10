package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/api/middleware/jwt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock JWT service
type mockJWTService struct {
	mock.Mock
	jwt.Service
}

func (m *mockJWTService) GenerateAccessToken(userID, email string, roles []string) (string, time.Time, error) {
	args := m.Called(userID, email, roles)
	return args.String(0), args.Get(1).(time.Time), args.Error(2)
}

func (m *mockJWTService) GenerateRefreshToken(userID string) (string, time.Time, error) {
	args := m.Called(userID)
	return args.String(0), args.Get(1).(time.Time), args.Error(2)
}

func (m *mockJWTService) ValidateAccessToken(token string) (*jwt.CustomClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.CustomClaims), args.Error(1)
}

func (m *mockJWTService) ValidateRefreshToken(token string) (*jwt.CustomClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.CustomClaims), args.Error(1)
}

func (m *mockJWTService) IsBlacklisted(token string) bool {
	args := m.Called(token)
	return args.Bool(0)
}

func (m *mockJWTService) GetRefreshTTL() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

func TestAuthMiddleware_Authenticate(t *testing.T) {
	tests := []struct {
		name           string
		setupAuth      func(r *http.Request)
		setupMock      func(*mockJWTService)
		expectedStatus int
		expectedClaims *jwt.CustomClaims
	}{
		{
			name: "valid token",
			setupAuth: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer valid_token")
			},
			setupMock: func(m *mockJWTService) {
				claims := &jwt.CustomClaims{
					UserID: "user123",
					Email:  "test@example.com",
					Roles:  []string{"user"},
				}
				m.On("ValidateAccessToken", "valid_token").Return(claims, nil)
				m.On("IsBlacklisted", "valid_token").Return(false)
			},
			expectedStatus: http.StatusOK,
			expectedClaims: &jwt.CustomClaims{
				UserID: "user123",
				Email:  "test@example.com",
				Roles:  []string{"user"},
			},
		},
		{
			name:           "missing authorization header",
			setupAuth:      func(r *http.Request) {},
			setupMock:      func(m *mockJWTService) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "invalid authorization format",
			setupAuth: func(r *http.Request) {
				r.Header.Set("Authorization", "NotBearer token")
			},
			setupMock:      func(m *mockJWTService) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "blacklisted token",
			setupAuth: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer blacklisted_token")
			},
			setupMock: func(m *mockJWTService) {
				claims := &jwt.CustomClaims{
					UserID: "user123",
					Email:  "test@example.com",
					Roles:  []string{"user"},
				}
				m.On("ValidateAccessToken", "blacklisted_token").Return(claims, nil)
				m.On("IsBlacklisted", "blacklisted_token").Return(true)
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock JWT service
			mockJWT := new(mockJWTService)
			tt.setupMock(mockJWT)

			// Create middleware
			middleware := NewAuthMiddleware(mockJWT)

			// Create test handler
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check if claims were properly set in context
				if tt.expectedClaims != nil {
					userID := r.Context().Value(UserIDKey)
					email := r.Context().Value(EmailKey)
					roles := r.Context().Value(RolesKey)

					assert.Equal(t, tt.expectedClaims.UserID, userID)
					assert.Equal(t, tt.expectedClaims.Email, email)
					assert.Equal(t, tt.expectedClaims.Roles, roles)
				}
				w.WriteHeader(http.StatusOK)
			})

			// Create test request
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			tt.setupAuth(r)
			w := httptest.NewRecorder()

			// Call middleware
			handler := middleware.Authenticate(nextHandler)
			handler.ServeHTTP(w, r)

			// Check response
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verify mock expectations
			mockJWT.AssertExpectations(t)
		})
	}
}

func TestAuthMiddleware_RequireRole(t *testing.T) {
	tests := []struct {
		name           string
		role           string
		userRoles      []string
		expectedStatus int
	}{
		{
			name:           "user has required role",
			role:           "admin",
			userRoles:      []string{"admin", "user"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "user does not have required role",
			role:           "admin",
			userRoles:      []string{"user"},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "empty user roles",
			role:           "admin",
			userRoles:      []string{},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock JWT service
			mockJWT := new(mockJWTService)

			// Create middleware
			middleware := NewAuthMiddleware(mockJWT)

			// Create test handler
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create test request with user roles in context
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			ctx := context.WithValue(r.Context(), RolesKey, tt.userRoles)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			// Call middleware
			handler := middleware.RequireRole(tt.role)(nextHandler)
			handler.ServeHTTP(w, r)

			// Check response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestAuthMiddleware_RequirePermission(t *testing.T) {
	tests := []struct {
		name           string
		permission     string
		userRoles      []string
		expectedStatus int
	}{
		{
			name:           "user has required permission through role",
			permission:     "read:users",
			userRoles:      []string{"admin"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "user does not have required permission",
			permission:     "delete:users",
			userRoles:      []string{"user"},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "empty user roles",
			permission:     "read:users",
			userRoles:      []string{},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock JWT service
			mockJWT := new(mockJWTService)

			// Create middleware
			middleware := NewAuthMiddleware(mockJWT)

			// Create test handler
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create test request with user roles in context
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			ctx := context.WithValue(r.Context(), RolesKey, tt.userRoles)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			// Call middleware
			handler := middleware.RequirePermission(tt.permission)(nextHandler)
			handler.ServeHTTP(w, r)

			// Check response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

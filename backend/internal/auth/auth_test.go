package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "missing auth header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid auth header format",
			authHeader:     "InvalidToken",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid bearer token",
			authHeader:     "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	// Create a mock service
	svc := &Service{
		// You would typically mock the clerk.Client here
		// For now, we'll just test the error cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test handler that would be wrapped by our middleware
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create a test request
			req := httptest.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Create our middleware and wrap the next handler
			handler := svc.AuthMiddleware(nextHandler)

			// Serve the request
			handler.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

func TestGetUserFromContext(t *testing.T) {
	tests := []struct {
		name        string
		setupCtx    func() context.Context
		wantUser    *UserData
		wantPresent bool
	}{
		{
			name: "user present in context",
			setupCtx: func() context.Context {
				user := &UserData{
					ID:       "test-id",
					Email:    "test@example.com",
					Username: "testuser",
				}
				return context.WithValue(context.Background(), userContextKey, user)
			},
			wantUser: &UserData{
				ID:       "test-id",
				Email:    "test@example.com",
				Username: "testuser",
			},
			wantPresent: true,
		},
		{
			name: "user not present in context",
			setupCtx: func() context.Context {
				return context.Background()
			},
			wantUser:    nil,
			wantPresent: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupCtx()
			gotUser, gotPresent := GetUserFromContext(ctx)

			if gotPresent != tt.wantPresent {
				t.Errorf("GetUserFromContext() present = %v, want %v", gotPresent, tt.wantPresent)
			}

			if tt.wantPresent {
				if gotUser.ID != tt.wantUser.ID {
					t.Errorf("GetUserFromContext() user.ID = %v, want %v", gotUser.ID, tt.wantUser.ID)
				}
				if gotUser.Email != tt.wantUser.Email {
					t.Errorf("GetUserFromContext() user.Email = %v, want %v", gotUser.Email, tt.wantUser.Email)
				}
				if gotUser.Username != tt.wantUser.Username {
					t.Errorf("GetUserFromContext() user.Username = %v, want %v", gotUser.Username, tt.wantUser.Username)
				}
			}
		})
	}
}

func TestRequireRole(t *testing.T) {
	tests := []struct {
		name           string
		setupCtx       func() context.Context
		role           string
		expectedStatus int
	}{
		{
			name: "user has required role",
			setupCtx: func() context.Context {
				user := &UserData{
					ID:    "test-id",
					Roles: []string{"admin", "user"},
				}
				return context.WithValue(context.Background(), userContextKey, user)
			},
			role:           "admin",
			expectedStatus: http.StatusOK,
		},
		{
			name: "user does not have required role",
			setupCtx: func() context.Context {
				user := &UserData{
					ID:    "test-id",
					Roles: []string{"user"},
				}
				return context.WithValue(context.Background(), userContextKey, user)
			},
			role:           "admin",
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "no user in context",
			setupCtx: func() context.Context {
				return context.Background()
			},
			role:           "admin",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	svc := &Service{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test handler that would be wrapped by our middleware
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create a test request with the context
			req := httptest.NewRequest("GET", "/", nil)
			req = req.WithContext(tt.setupCtx())

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Create our middleware and wrap the next handler
			handler := svc.RequireRole(tt.role, nextHandler)

			// Serve the request
			handler.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

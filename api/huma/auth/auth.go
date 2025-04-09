// Package auth provides the authentication endpoints for the Huma API.
package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName,omitempty"`
	LastName  string    `json:"lastName,omitempty"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
	TokenType    string    `json:"tokenType"`
	User         User      `json:"user"`
}

// RegisterEndpoints registers the authentication endpoints.
func RegisterEndpoints(api huma.API, basePath string) {
	// POST /auth/login
	huma.Register(api, huma.Operation{
		OperationID: "login",
		Method:      http.MethodPost,
		Path:        basePath + "/auth/login",
		Summary:     "Login",
		Description: "Authenticates a user and returns access and refresh tokens",
		Tags:        []string{"Authentication"},
	}, func(ctx context.Context, input *struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}) (*struct {
		Body AuthResponse
	}, error) {
		resp := &struct {
			Body AuthResponse
		}{}

		// In a real implementation, we would validate the credentials
		// and generate a JWT token. For now, we'll just return a mock response.
		resp.Body = AuthResponse{
			AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyLTEyMzQ1NiIsImVtYWlsIjoidXNlckBleGFtcGxlLmNvbSIsImlhdCI6MTUxNjIzOTAyMn0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyLTEyMzQ1NiIsInR5cCI6InJlZnJlc2giLCJpYXQiOjE1MTYyMzkwMjJ9.oK5Jw2NP5R1Cs7QTf_5jlZcxEv_YK4ZNNsGSPcLtXbE",
			ExpiresAt:    time.Now().Add(time.Hour),
			TokenType:    "Bearer",
			User: User{
				ID:        "user-123456",
				Email:     input.Email,
				Username:  "johndoe",
				FirstName: "John",
				LastName:  "Doe",
				Roles:     []string{"user"},
				CreatedAt: time.Now().AddDate(0, -1, 0),
				UpdatedAt: time.Now(),
			},
		}

		return resp, nil
	})

	// POST /auth/register
	huma.Register(api, huma.Operation{
		OperationID: "register",
		Method:      http.MethodPost,
		Path:        basePath + "/auth/register",
		Summary:     "Register",
		Description: "Registers a new user",
		Tags:        []string{"Authentication"},
	}, func(ctx context.Context, input *struct {
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		FirstName string `json:"firstName,omitempty"`
		LastName  string `json:"lastName,omitempty"`
	}) (*struct {
		Body AuthResponse
	}, error) {
		resp := &struct {
			Body AuthResponse
		}{}

		// In a real implementation, we would create a new user in the database
		// and generate a JWT token. For now, we'll just return a mock response.
		resp.Body = AuthResponse{
			AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyLW5ldyIsImVtYWlsIjoibmV3dXNlckBleGFtcGxlLmNvbSIsImlhdCI6MTUxNjIzOTAyMn0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyLW5ldyIsInR5cCI6InJlZnJlc2giLCJpYXQiOjE1MTYyMzkwMjJ9.oK5Jw2NP5R1Cs7QTf_5jlZcxEv_YK4ZNNsGSPcLtXbE",
			ExpiresAt:    time.Now().Add(time.Hour),
			TokenType:    "Bearer",
			User: User{
				ID:        "user-" + uuid.New().String()[:6],
				Email:     input.Email,
				Username:  input.Username,
				FirstName: input.FirstName,
				LastName:  input.LastName,
				Roles:     []string{"user"},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		return resp, nil
	})

	// POST /auth/refresh
	huma.Register(api, huma.Operation{
		OperationID: "refresh-token",
		Method:      http.MethodPost,
		Path:        basePath + "/auth/refresh",
		Summary:     "Refresh token",
		Description: "Refreshes an access token using a refresh token",
		Tags:        []string{"Authentication"},
	}, func(ctx context.Context, input *struct {
		RefreshToken string `json:"refreshToken"`
	}) (*struct {
		Body struct {
			AccessToken  string    `json:"accessToken"`
			RefreshToken string    `json:"refreshToken"`
			ExpiresAt    time.Time `json:"expiresAt"`
			TokenType    string    `json:"tokenType"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				AccessToken  string    `json:"accessToken"`
				RefreshToken string    `json:"refreshToken"`
				ExpiresAt    time.Time `json:"expiresAt"`
				TokenType    string    `json:"tokenType"`
			}
		}{}

		// In a real implementation, we would validate the refresh token
		// and generate a new JWT token. For now, we'll just return a mock response.
		resp.Body.AccessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyLTEyMzQ1NiIsImVtYWlsIjoidXNlckBleGFtcGxlLmNvbSIsImlhdCI6MTUxNjIzOTAyMn0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
		resp.Body.RefreshToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyLTEyMzQ1NiIsInR5cCI6InJlZnJlc2giLCJpYXQiOjE1MTYyMzkwMjJ9.oK5Jw2NP5R1Cs7QTf_5jlZcxEv_YK4ZNNsGSPcLtXbE"
		resp.Body.ExpiresAt = time.Now().Add(time.Hour)
		resp.Body.TokenType = "Bearer"

		return resp, nil
	})

	// POST /auth/logout
	huma.Register(api, huma.Operation{
		OperationID: "logout",
		Method:      http.MethodPost,
		Path:        basePath + "/auth/logout",
		Summary:     "Logout",
		Description: "Invalidates the current session",
		Tags:        []string{"Authentication"},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body struct {
			Success   bool      `json:"success"`
			Timestamp time.Time `json:"timestamp"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				Success   bool      `json:"success"`
				Timestamp time.Time `json:"timestamp"`
			}
		}{}

		// In a real implementation, we would invalidate the token in the database
		// or add it to a blacklist. For now, we'll just return a success response.
		resp.Body.Success = true
		resp.Body.Timestamp = time.Now()

		return resp, nil
	})

	// GET /auth/verify
	huma.Register(api, huma.Operation{
		OperationID: "verify-token",
		Method:      http.MethodGet,
		Path:        basePath + "/auth/verify",
		Summary:     "Verify token",
		Description: "Verifies the current access token",
		Tags:        []string{"Authentication"},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body struct {
			Valid     bool      `json:"valid"`
			ExpiresAt time.Time `json:"expiresAt"`
			User      User      `json:"user"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				Valid     bool      `json:"valid"`
				ExpiresAt time.Time `json:"expiresAt"`
				User      User      `json:"user"`
			}
		}{}

		// In a real implementation, we would validate the token from the Authorization header
		// and return the user information. For now, we'll just return a mock response.
		resp.Body.Valid = true
		resp.Body.ExpiresAt = time.Now().Add(time.Hour)
		resp.Body.User = User{
			ID:        "user-123456",
			Email:     "user@example.com",
			Username:  "johndoe",
			Roles:     []string{"user"},
			CreatedAt: time.Now().AddDate(0, -1, 0),
			UpdatedAt: time.Now(),
		}

		return resp, nil
	})
}

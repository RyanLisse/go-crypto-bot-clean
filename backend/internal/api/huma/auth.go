package huma

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

// LoginRequest represents a request to login
type LoginRequest struct {
	Body struct {
		Email    string `json:"email" doc:"Email address" example:"user@example.com" format:"email" binding:"required"`
		Password string `json:"password" doc:"Password" example:"password123" binding:"required"`
	}
}

// RegisterRequest represents a request to register a new user
type RegisterRequest struct {
	Body struct {
		Email     string `json:"email" doc:"Email address" example:"user@example.com" format:"email" binding:"required"`
		Username  string `json:"username" doc:"Username" example:"johndoe" binding:"required"`
		Password  string `json:"password" doc:"Password" example:"password123" binding:"required"`
		FirstName string `json:"firstName,omitempty" doc:"First name" example:"John"`
		LastName  string `json:"lastName,omitempty" doc:"Last name" example:"Doe"`
	}
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	Body struct {
		AccessToken  string    `json:"accessToken" doc:"JWT access token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
		RefreshToken string    `json:"refreshToken" doc:"JWT refresh token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
		ExpiresAt    time.Time `json:"expiresAt" doc:"Expiration time of the access token" example:"2023-02-02T11:00:00Z"`
		TokenType    string    `json:"tokenType" doc:"Type of token" example:"Bearer"`
		User         struct {
			ID        string   `json:"id" doc:"User ID" example:"user-123456"`
			Email     string   `json:"email" doc:"Email address" example:"user@example.com"`
			Username  string   `json:"username" doc:"Username" example:"johndoe"`
			FirstName string   `json:"firstName,omitempty" doc:"First name" example:"John"`
			LastName  string   `json:"lastName,omitempty" doc:"Last name" example:"Doe"`
			Roles     []string `json:"roles" doc:"User roles" example:"[\"user\", \"admin\"]"`
		} `json:"user" doc:"User information"`
	}
}

// RefreshTokenRequest represents a request to refresh an access token
type RefreshTokenRequest struct {
	Body struct {
		RefreshToken string `json:"refreshToken" doc:"JWT refresh token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." binding:"required"`
	}
}

// registerAuthEndpoints registers the authentication endpoints.
func registerAuthEndpoints(api huma.API, basePath string) {
	// POST /auth/login
	huma.Register(api, huma.Operation{
		OperationID: "login",
		Method:      http.MethodPost,
		Path:        basePath + "/auth/login",
		Summary:     "Login",
		Description: "Authenticates a user and returns access and refresh tokens",
		Tags:        []string{"Authentication"},
	}, func(ctx context.Context, input *LoginRequest) (*AuthResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})

	// POST /auth/register
	huma.Register(api, huma.Operation{
		OperationID: "register",
		Method:      http.MethodPost,
		Path:        basePath + "/auth/register",
		Summary:     "Register",
		Description: "Registers a new user",
		Tags:        []string{"Authentication"},
	}, func(ctx context.Context, input *RegisterRequest) (*AuthResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})

	// POST /auth/refresh
	huma.Register(api, huma.Operation{
		OperationID: "refresh-token",
		Method:      http.MethodPost,
		Path:        basePath + "/auth/refresh",
		Summary:     "Refresh token",
		Description: "Refreshes an access token using a refresh token",
		Tags:        []string{"Authentication"},
	}, func(ctx context.Context, input *RefreshTokenRequest) (*AuthResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
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
			Success   bool      `json:"success" doc:"Whether the logout was successful" example:"true"`
			Timestamp time.Time `json:"timestamp" doc:"Timestamp of the logout" example:"2023-02-02T10:00:00Z"`
		}
	}, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
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
			Valid     bool      `json:"valid" doc:"Whether the token is valid" example:"true"`
			ExpiresAt time.Time `json:"expiresAt" doc:"Expiration time of the token" example:"2023-02-02T11:00:00Z"`
			User      struct {
				ID        string   `json:"id" doc:"User ID" example:"user-123456"`
				Email     string   `json:"email" doc:"Email address" example:"user@example.com"`
				Username  string   `json:"username" doc:"Username" example:"johndoe"`
				Roles     []string `json:"roles" doc:"User roles" example:"[\"user\", \"admin\"]"`
			} `json:"user" doc:"User information"`
		}
	}, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})
}

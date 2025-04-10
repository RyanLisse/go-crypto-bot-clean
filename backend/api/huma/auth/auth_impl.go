// Package auth provides the authentication endpoints for the Huma API.
package auth

import (
	"context"
	"net/http"
	"time"

	"go-crypto-bot-clean/backend/api/service"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterAuthEndpoints registers the authentication endpoints with service implementation.
func RegisterAuthEndpoints(api huma.API, basePath string, authService *service.AuthService) {
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
		Body service.AuthResponse
	}, error) {
		// Convert API request to service request
		req := &service.LoginRequest{
			Email:    input.Email,
			Password: input.Password,
		}

		// Login
		result, err := authService.Login(ctx, req)
		if err != nil {
			return nil, err
		}

		// Return result
		return &struct {
			Body service.AuthResponse
		}{
			Body: *result,
		}, nil
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
		Body service.AuthResponse
	}, error) {
		// Convert API request to service request
		req := &service.RegisterRequest{
			Email:     input.Email,
			Username:  input.Username,
			Password:  input.Password,
			FirstName: input.FirstName,
			LastName:  input.LastName,
		}

		// Register
		result, err := authService.Register(ctx, req)
		if err != nil {
			return nil, err
		}

		// Return result
		return &struct {
			Body service.AuthResponse
		}{
			Body: *result,
		}, nil
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
		// Convert API request to service request
		req := &service.RefreshTokenRequest{
			RefreshToken: input.RefreshToken,
		}

		// Refresh token
		result, err := authService.RefreshToken(ctx, req)
		if err != nil {
			return nil, err
		}

		// Return result
		resp := &struct {
			Body struct {
				AccessToken  string    `json:"accessToken"`
				RefreshToken string    `json:"refreshToken"`
				ExpiresAt    time.Time `json:"expiresAt"`
				TokenType    string    `json:"tokenType"`
			}
		}{}
		resp.Body.AccessToken = result.AccessToken
		resp.Body.RefreshToken = result.RefreshToken
		resp.Body.ExpiresAt = result.ExpiresAt
		resp.Body.TokenType = result.TokenType

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
	}, func(ctx context.Context, input *struct {
		RefreshToken string `json:"refreshToken,omitempty"`
	}) (*struct {
		Body struct {
			Success   bool      `json:"success"`
			Timestamp time.Time `json:"timestamp"`
		}
	}, error) {
		// For now, we'll just use a mock user ID
		// In a real implementation, we would get the user ID from the JWT token
		userID := "user-123456"

		// Get refresh token from request
		var refreshToken string
		if input != nil {
			refreshToken = input.RefreshToken
		}

		// Logout
		err := authService.Logout(ctx, userID, refreshToken)
		if err != nil {
			return nil, err
		}

		// Return result
		resp := &struct {
			Body struct {
				Success   bool      `json:"success"`
				Timestamp time.Time `json:"timestamp"`
			}
		}{}
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
			Valid     bool         `json:"valid"`
			ExpiresAt time.Time    `json:"expiresAt"`
			User      service.User `json:"user"`
		}
	}, error) {
		// Get token from context
		// In a real implementation, we would get the token from the Authorization header
		token := "mock-token"

		// Verify token
		result, err := authService.VerifyToken(ctx, token)
		if err != nil {
			return nil, err
		}

		// Return result
		resp := &struct {
			Body struct {
				Valid     bool         `json:"valid"`
				ExpiresAt time.Time    `json:"expiresAt"`
				User      service.User `json:"user"`
			}
		}{}
		resp.Body.Valid = result.Valid
		resp.Body.ExpiresAt = result.ExpiresAt
		resp.Body.User = result.User

		return resp, nil
	})
}

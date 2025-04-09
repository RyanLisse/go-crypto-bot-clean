package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go-crypto-bot-clean/backend/internal/api/dto/response"
	"go-crypto-bot-clean/backend/internal/api/middleware"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	jwtSecret   string
	jwtExpiry   int
	cookieName  string
	credentials map[string]string // username -> password (in a real app, use a proper user store)
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(jwtSecret string, jwtExpiry int, cookieName string) *AuthHandler {
	// In a real application, you would use a database or other storage
	// This is just a simple example with hardcoded credentials
	credentials := map[string]string{
		"admin": "admin123",
		"user":  "user123",
	}

	return &AuthHandler{
		jwtSecret:   jwtSecret,
		jwtExpiry:   jwtExpiry,
		cookieName:  cookieName,
		credentials: credentials,
	}
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response body
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"`
}

// Login handles user login and token generation
//
//	@summary	User login
//	@description	Authenticates a user and returns a JWT token
//	@tags		Auth
//	@accept		json
//	@produce	json
//	@param		request	body		LoginRequest	true	"Login credentials"
//	@success	200		{object}	LoginResponse
//	@failure	400		{object}	response.ErrorResponse
//	@failure	401		{object}	response.ErrorResponse
//	@router		/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    "invalid_request",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Check credentials
	storedPassword, exists := h.credentials[req.Username]
	if !exists || storedPassword != req.Password {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Code:    "unauthorized",
			Message: "Invalid username or password",
		})
		return
	}

	// Determine role (in a real app, this would come from the user store)
	role := "user"
	if req.Username == "admin" {
		role = "admin"
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(req.Username, role, h.jwtSecret, h.jwtExpiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Code:    "token_generation_failed",
			Message: "Failed to generate token",
			Details: err.Error(),
		})
		return
	}

	// Calculate expiry time
	expiresAt := time.Now().Add(time.Duration(h.jwtExpiry) * time.Hour)

	// Set cookie if cookie name is provided
	if h.cookieName != "" {
		c.SetCookie(
			h.cookieName,
			token,
			h.jwtExpiry*3600, // max age in seconds
			"/",              // path
			"",               // domain
			false,            // secure
			true,             // http only
		)
	}

	// Return token in response
	c.JSON(http.StatusOK, LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		UserID:    req.Username,
		Role:      role,
	})
}

// Logout handles user logout
//
//	@summary	User logout
//	@description	Logs out a user by clearing the auth cookie
//	@tags		Auth
//	@produce	json
//	@success	200		{object}	map[string]string
//	@router		/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Clear the cookie if cookie name is provided
	if h.cookieName != "" {
		c.SetCookie(
			h.cookieName,
			"",        // empty value
			-1,        // max age: delete immediately
			"/",       // path
			"",        // domain
			false,     // secure
			true,      // http only
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

// GetCurrentUser returns the current authenticated user
//
//	@summary	Get current user
//	@description	Returns information about the currently authenticated user
//	@tags		Auth
//	@produce	json
//	@success	200		{object}	map[string]string
//	@failure	401		{object}	response.ErrorResponse
//	@router		/auth/me [get]
//	@security	BearerAuth
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// Get user ID and role from context (set by the JWT middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Code:    "unauthorized",
			Message: "Not authenticated",
		})
		return
	}

	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"role":    role,
	})
}

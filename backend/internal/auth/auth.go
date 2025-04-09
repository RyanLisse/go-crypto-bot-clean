package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/clerk/clerk-sdk-go/v2/user"
)

// Config holds the configuration for the auth package
type Config struct {
	ClerkSecretKey string
}

// UserDataKey is the context key for user data
type userDataKey struct{}

var UserDataKey = userDataKey{}

// AuthProvider defines the interface for authentication providers
type AuthProvider interface {
	// Authenticate validates the request and returns user data
	Authenticate(r *http.Request) (*UserData, error)
	// RequireRole creates middleware that checks if the user has the required role
	RequireRole(role string) func(http.Handler) http.Handler
	// RequirePermission creates middleware that checks if the user has the required permission
	RequirePermission(permission string) func(http.Handler) http.Handler
	// RequireAnyPermission creates middleware that checks if the user has any of the required permissions
	RequireAnyPermission(permissions ...string) func(http.Handler) http.Handler
	// RequireAllPermissions creates middleware that checks if the user has all required permissions
	RequireAllPermissions(permissions ...string) func(http.Handler) http.Handler
}

// Service implements the AuthProvider interface
type Service struct {
	secretKey string
}

// Authenticate validates the request and returns user data
func (s *Service) Authenticate(r *http.Request) (*UserData, error) {
	sessionClaims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		return nil, ErrUnauthorized
	}

	// Get user data using the user service
	usr, err := user.Get(r.Context(), sessionClaims.Subject)
	if err != nil {
		return nil, NewAuthError(ErrorTypeInvalidToken, "Failed to get user data", http.StatusUnauthorized)
	}

	// Map Clerk user data to our UserData struct
	email := ""
	if len(usr.EmailAddresses) > 0 {
		email = usr.EmailAddresses[0].EmailAddress
	}

	// Handle nullable string fields
	username := ""
	if usr.Username != nil {
		username = *usr.Username
	}

	firstName := ""
	if usr.FirstName != nil {
		firstName = *usr.FirstName
	}

	lastName := ""
	if usr.LastName != nil {
		lastName = *usr.LastName
	}

	imageURL := ""
	if usr.ImageURL != nil {
		imageURL = *usr.ImageURL
	}

	userData := &UserData{
		ID:              usr.ID,
		Email:           email,
		Username:        username,
		FirstName:       firstName,
		LastName:        lastName,
		ProfileImageURL: imageURL,
		Roles:           []string{"user"}, // Default role
	}

	return userData, nil
}

// NewService creates a new auth service
func NewService(secretKey string) *Service {
	clerk.SetKey(secretKey)
	return &Service{
		secretKey: secretKey,
	}
}

// UserData represents the user information from Clerk
type UserData struct {
	ID              string   `json:"id"`
	Email           string   `json:"email"`
	Username        string   `json:"username"`
	FirstName       string   `json:"first_name"`
	LastName        string   `json:"last_name"`
	ProfileImageURL string   `json:"profile_image_url"`
	Roles           []string `json:"roles"`
}

// AuthMiddleware creates a middleware that verifies the JWT token from Clerk
func (s *Service) AuthMiddleware(next http.Handler) http.Handler {
	return clerkhttp.RequireHeaderAuthorization()(next)
}

// GetUserFromContext retrieves the user data from the context
func GetUserFromContext(ctx context.Context) (*UserData, bool) {
	claims, ok := clerk.SessionClaimsFromContext(ctx)
	if !ok {
		return nil, false
	}

	usr, err := user.Get(ctx, claims.Subject)
	if err != nil {
		return nil, false
	}

	email := ""
	if len(usr.EmailAddresses) > 0 {
		email = usr.EmailAddresses[0].EmailAddress
	}

	// Handle nullable string fields
	username := ""
	if usr.Username != nil {
		username = *usr.Username
	}

	firstName := ""
	if usr.FirstName != nil {
		firstName = *usr.FirstName
	}

	lastName := ""
	if usr.LastName != nil {
		lastName = *usr.LastName
	}

	imageURL := ""
	if usr.ImageURL != nil {
		imageURL = *usr.ImageURL
	}

	// Extract roles from private metadata
	roles := []string{RoleUser} // Default role
	if usr.PrivateMetadata != nil {
		var metadata map[string]interface{}
		if err := json.Unmarshal(usr.PrivateMetadata, &metadata); err == nil {
			if rolesRaw, ok := metadata["roles"].([]interface{}); ok {
				roles = make([]string, 0, len(rolesRaw))
				for _, r := range rolesRaw {
					if role, ok := r.(string); ok && ValidateRole(role) {
						roles = append(roles, role)
					}
				}
			}
		}
	}

	userData := &UserData{
		ID:              usr.ID,
		Email:           email,
		Username:        username,
		FirstName:       firstName,
		LastName:        lastName,
		ProfileImageURL: imageURL,
		Roles:           roles,
	}

	return userData, true
}

// RequireRole creates middleware that checks if the user has the required role
func (s *Service) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userData := r.Context().Value(UserDataKey).(UserData)

			for _, userRole := range userData.Roles {
				if userRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			err := NewAuthError(ErrorTypeForbidden, "Insufficient role", http.StatusForbidden).
				WithDetails(map[string]interface{}{
					"required_role": role,
					"user_roles":    userData.Roles,
				})
			err.WriteJSON(w)
		})
	}
}

// RequirePermission creates middleware that checks if the user has the required permission
func (s *Service) RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userData := r.Context().Value(UserDataKey).(UserData)

			for _, role := range userData.Roles {
				if HasPermission(role, permission) {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
		})
	}
}

// RequireAnyPermission creates middleware that checks if the user has any of the required permissions
func (s *Service) RequireAnyPermission(permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userData := r.Context().Value(UserDataKey).(UserData)

			for _, role := range userData.Roles {
				for _, permission := range permissions {
					if HasPermission(role, permission) {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
		})
	}
}

// RequireAllPermissions creates middleware that checks if the user has all required permissions
func (s *Service) RequireAllPermissions(permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userData := r.Context().Value(UserDataKey).(UserData)

			for _, permission := range permissions {
				hasPermission := false
				for _, role := range userData.Roles {
					if HasPermission(role, permission) {
						hasPermission = true
						break
					}
				}
				if !hasPermission {
					http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Helper function to safely get string value from pointer
func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// DisabledService is a no-op implementation of the auth service
type DisabledService struct{}

// Authenticate for DisabledService returns an error indicating the service is disabled
func (s *DisabledService) Authenticate(r *http.Request) (*UserData, error) {
	return nil, NewAuthError(
		ErrorTypeServiceUnavailable,
		"Authentication service is disabled",
		http.StatusServiceUnavailable,
	).WithHelp("The authentication service is currently disabled. Authentication cannot be performed.")
}

// AuthMiddleware for DisabledService allows all requests through
func (s *DisabledService) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// RequireRole for DisabledService returns an error indicating the service is disabled
func (s *DisabledService) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create an error response indicating the service is disabled
			err := NewAuthError(
				ErrorTypeServiceUnavailable,
				"Authentication service is disabled",
				http.StatusServiceUnavailable,
			).WithDetails(map[string]interface{}{
				"requested_role": role,
				"service_status": "disabled",
			}).WithHelp("The authentication service is currently disabled. Role checks cannot be performed.")

			// Write the error response
			err.WriteJSON(w)
		})
	}
}

// RequirePermission for DisabledService returns an error indicating the service is disabled
func (s *DisabledService) RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create an error response indicating the service is disabled
			err := NewAuthError(
				ErrorTypeServiceUnavailable,
				"Authentication service is disabled",
				http.StatusServiceUnavailable,
			).WithDetails(map[string]interface{}{
				"requested_permission": permission,
				"service_status":       "disabled",
			}).WithHelp("The authentication service is currently disabled. Permission checks cannot be performed.")

			// Write the error response
			err.WriteJSON(w)
		})
	}
}

// RequireAnyPermission for DisabledService returns an error indicating the service is disabled
func (s *DisabledService) RequireAnyPermission(permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create an error response indicating the service is disabled
			err := NewAuthError(
				ErrorTypeServiceUnavailable,
				"Authentication service is disabled",
				http.StatusServiceUnavailable,
			).WithDetails(map[string]interface{}{
				"requested_permissions": permissions,
				"service_status":        "disabled",
			}).WithHelp("The authentication service is currently disabled. Permission checks cannot be performed.")

			// Write the error response
			err.WriteJSON(w)
		})
	}
}

// RequireAllPermissions for DisabledService returns an error indicating the service is disabled
func (s *DisabledService) RequireAllPermissions(permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create an error response indicating the service is disabled
			err := NewAuthError(
				ErrorTypeServiceUnavailable,
				"Authentication service is disabled",
				http.StatusServiceUnavailable,
			).WithDetails(map[string]interface{}{
				"required_permissions": permissions,
				"service_status":       "disabled",
			}).WithHelp("The authentication service is currently disabled. Permission checks cannot be performed.")

			// Write the error response
			err.WriteJSON(w)
		})
	}
}

// NewDisabledService creates a new disabled auth service for testing or development
func NewDisabledService() *DisabledService {
	return &DisabledService{}
}

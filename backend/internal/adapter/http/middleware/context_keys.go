package middleware

import "context"

type contextKey string

var (
    UserIDKey = contextKey("user_id")
    RoleKey   = contextKey("role")
)

// GetUserIDFromContext extracts the user ID from the context, returns (userID, ok)
func GetUserIDFromContext(ctx context.Context) (string, bool) {
    val := ctx.Value(UserIDKey)
    if userID, ok := val.(string); ok && userID != "" {
        return userID, true
    }
    return "", false
}

package scripts

import (
	"context"
	"fmt"
	"os"
	"time"

	repo "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/repository/gorm"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	sqlitegorm "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestAuth tests the authentication system
func TestAuth() {
	// Configure logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	logger := log.With().Str("component", "auth-tester").Logger()

	// Create a temporary database
	tempFile, err := os.CreateTemp("", "auth-test-*.db")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create temporary file")
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Connect to the database
	db, err := gorm.Open(sqlitegorm.Open(tempFile.Name()), &gorm.Config{})
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Create tables
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create tables")
	}

	// Create repositories
	userRepo := repo.NewUserRepository(db)

	// Create services
	userService := service.NewUserService(userRepo)

	// Test user service
	logger.Info().Msg("Testing user service")
	testUserService(userService, &logger)

	// Get Clerk secret key from environment
	clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkSecretKey == "" {
		logger.Warn().Msg("CLERK_SECRET_KEY environment variable not set, skipping auth service test")
		return
	}

	// Create auth service
	authService, err := service.NewAuthService(userService, clerkSecretKey)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create auth service")
	}

	// Test auth service
	logger.Info().Msg("Testing auth service")
	testAuthService(authService, &logger)
}

// testUserService tests the user service
func testUserService(userService *service.UserService, logger *zerolog.Logger) {
	ctx := context.Background()

	// Create a user
	user, err := userService.CreateUser(ctx, "test-user-id", "test@example.com", "Test User")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create user")
	}
	logger.Info().Str("id", user.ID).Str("email", user.Email).Str("name", user.Name).Msg("Created user")

	// Get user by ID
	user, err = userService.GetUserByID(ctx, "test-user-id")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get user by ID")
	}
	logger.Info().Str("id", user.ID).Str("email", user.Email).Str("name", user.Name).Msg("Got user by ID")

	// Get user by email
	user, err = userService.GetUserByEmail(ctx, "test@example.com")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get user by email")
	}
	logger.Info().Str("id", user.ID).Str("email", user.Email).Str("name", user.Name).Msg("Got user by email")

	// Update user
	user, err = userService.UpdateUser(ctx, "test-user-id", "Updated User")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to update user")
	}
	logger.Info().Str("id", user.ID).Str("email", user.Email).Str("name", user.Name).Msg("Updated user")

	// List users
	users, err := userService.ListUsers(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to list users")
	}
	logger.Info().Int("count", len(users)).Msg("Listed users")

	// Ensure user exists
	user, err = userService.EnsureUserExists(ctx, "test-user-id", "updated@example.com", "Ensured User")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to ensure user exists")
	}
	logger.Info().Str("id", user.ID).Str("email", user.Email).Str("name", user.Name).Msg("Ensured user exists")

	// Delete user
	err = userService.DeleteUser(ctx, "test-user-id")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to delete user")
	}
	logger.Info().Msg("Deleted user")
}

// testAuthService tests the auth service
func testAuthService(authService *service.AuthService, logger *zerolog.Logger) {
	ctx := context.Background()

	// Get token from environment
	token := os.Getenv("CLERK_TEST_TOKEN")
	if token == "" {
		logger.Warn().Msg("CLERK_TEST_TOKEN environment variable not set, skipping token verification test")
		return
	}

	// Verify token
	userID, err := authService.VerifyToken(ctx, token)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to verify token")
		return
	}
	logger.Info().Str("userID", userID).Msg("Verified token")

	// Get user from token
	user, err := authService.GetUserFromToken(ctx, token)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get user from token")
		return
	}
	logger.Info().Str("id", user.ID).Str("email", user.Email).Str("name", user.Name).Msg("Got user from token")

	// Get user roles
	roles, err := authService.GetUserRoles(ctx, userID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get user roles")
		return
	}
	logger.Info().Strs("roles", roles).Msg("Got user roles")
}

// GenerateClerkToken generates a Clerk token for testing
func GenerateClerkToken() {
	// Get Clerk secret key from environment
	clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkSecretKey == "" {
		fmt.Println("CLERK_SECRET_KEY environment variable not set")
		return
	}

	// Get user ID from environment
	userID := os.Getenv("CLERK_TEST_USER_ID")
	if userID == "" {
		fmt.Println("CLERK_TEST_USER_ID environment variable not set")
		return
	}

	// Create token
	token, err := createClerkToken(clerkSecretKey, userID)
	if err != nil {
		fmt.Printf("Failed to create token: %v\n", err)
		return
	}

	fmt.Printf("Token: %s\n", token)
}

// createClerkToken creates a Clerk token for testing
func createClerkToken(secretKey, userID string) (string, error) {
	// This is a simplified version of token creation
	// In a real implementation, you would use the Clerk SDK to create a token
	// or use the Clerk API to create a session
	return fmt.Sprintf("test-token-%s-%d", userID, time.Now().Unix()), nil
}

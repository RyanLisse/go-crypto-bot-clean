# Authentication System Documentation

This document provides comprehensive documentation for the authentication system used in the wallet authentication system.

## Overview

The authentication system uses Clerk for user authentication and provides a secure way to handle user authentication and authorization. It supports both basic and enhanced authentication modes, with the enhanced mode providing additional features such as database integration and role-based access control.

## Components

### 1. Authentication Service

The authentication service provides a simple interface for authenticating users and managing user-related operations. It integrates with Clerk for user authentication and provides methods for verifying tokens, getting user information, and managing user roles.

### 2. Middleware

The authentication system provides two types of middleware:

- **Basic Clerk Middleware**: A simple middleware that validates Clerk authentication tokens and sets user information in the request context.
- **Enhanced Clerk Middleware**: An advanced middleware that integrates with the database and provides additional features such as role-based access control.

### 3. User Service

The user service provides methods for managing users in the database. It handles operations such as creating, updating, and deleting users, as well as retrieving user information.

### 4. Controllers

The authentication system provides two controllers:

- **User Controller**: Handles user-related HTTP requests such as getting user information and updating user profiles.
- **Auth Controller**: Handles authentication-related HTTP requests such as verifying tokens.

## Configuration

The authentication system can be configured using environment variables:

- `CLERK_SECRET_KEY`: The Clerk secret key used for authentication.
- `USE_ENHANCED_AUTH`: Whether to use the enhanced authentication mode (default: false).

## Usage

### 1. Authentication Service

```go
// Create authentication service
authService, err := service.NewAuthService(userService, secretKey)
if err != nil {
    // Handle error
}

// Verify token
userID, err := authService.VerifyToken(ctx, token)
if err != nil {
    // Handle error
}

// Get user from token
user, err := authService.GetUserFromToken(ctx, token)
if err != nil {
    // Handle error
}

// Get user roles
roles, err := authService.GetUserRoles(ctx, userID)
if err != nil {
    // Handle error
}
```

### 2. Middleware

```go
// Create middleware
authMiddleware, err := factory.CreateEnhancedClerkMiddleware(secretKey)
if err != nil {
    // Handle error
}

// Use middleware
router.Use(authMiddleware.Middleware())

// Require authentication
router.Group(func(r chi.Router) {
    r.Use(authMiddleware.RequireAuthentication)
    // Protected routes
})

// Require role
router.Group(func(r chi.Router) {
    r.Use(authMiddleware.RequireRole("admin"))
    // Admin-only routes
})
```

### 3. User Service

```go
// Create user service
userService := service.NewUserService(userRepo)

// Get user by ID
user, err := userService.GetUserByID(ctx, id)
if err != nil {
    // Handle error
}

// Create user
user, err := userService.CreateUser(ctx, id, email, name)
if err != nil {
    // Handle error
}

// Update user
user, err := userService.UpdateUser(ctx, id, name)
if err != nil {
    // Handle error
}

// Delete user
err := userService.DeleteUser(ctx, id)
if err != nil {
    // Handle error
}
```

## Security Considerations

### 1. Token Validation

- Tokens are validated using Clerk's JWT verification.
- The token's signature is verified to ensure it was issued by Clerk.
- The token's expiration time is checked to ensure it is still valid.
- The token's subject (user ID) is extracted and used to identify the user.

### 2. Role-Based Access Control

- Users can have one or more roles.
- Roles are stored in Clerk's public metadata.
- The middleware can require specific roles for accessing certain routes.
- If a user doesn't have the required role, they will receive a 403 Forbidden response.

### 3. Database Integration

- User information is stored in the database.
- When a user authenticates, their information is retrieved from Clerk and stored in the database if it doesn't already exist.
- This ensures that the database always has up-to-date user information.

## Best Practices

1. **Use Enhanced Authentication**: The enhanced authentication mode provides additional security features such as database integration and role-based access control.

2. **Validate Tokens**: Always validate tokens before trusting the user's identity.

3. **Use Role-Based Access Control**: Use role-based access control to restrict access to sensitive operations.

4. **Keep Clerk Secret Key Secure**: The Clerk secret key should be kept secure and never committed to version control.

5. **Use HTTPS**: Always use HTTPS in production to ensure that tokens are transmitted securely.

## Troubleshooting

### 1. Token Validation Issues

- **Invalid Token**: Ensure that the token is valid and has not expired.
- **Missing Token**: Ensure that the token is included in the Authorization header.

### 2. Role-Based Access Control Issues

- **Missing Role**: Ensure that the user has the required role.
- **Role Not Found**: Ensure that the roles are stored in Clerk's public metadata.

### 3. Database Integration Issues

- **User Not Found**: Ensure that the user exists in the database.
- **Database Connection Issues**: Ensure that the database connection is working properly.

## References

- [Clerk Documentation](https://clerk.dev/docs)
- [JWT Documentation](https://jwt.io/introduction)
- [Role-Based Access Control](https://en.wikipedia.org/wiki/Role-based_access_control)

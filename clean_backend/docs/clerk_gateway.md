# Clerk Authentication Gateway

This document describes the Clerk authentication gateway implementation in the Go Crypto Bot backend.

## Overview

The application uses [Clerk](https://clerk.com/) as its authentication provider. The Clerk gateway provides:

- Session token verification
- User information retrieval
- Caching of user data for performance
- Mock implementation for development and testing

## Configuration

Clerk authentication configuration is defined in the `config.yaml` file and can be overridden with environment variables:

```yaml
auth:
  provider: "clerk"
  disabled: false
  clerk_secret_key: "${CLERK_SECRET_KEY}"
  clerk_jwt_public_key: "${CLERK_JWT_PUBLIC_KEY}"
  clerk_jwt_template: "api_auth"
  use_enhanced: false
```

### Environment Variables

- `AUTH_PROVIDER`: Set to "clerk" to use Clerk authentication
- `AUTH_DISABLED`: Set to "true" to disable authentication (for development only)
- `CLERK_SECRET_KEY`: Your Clerk secret key
- `CLERK_JWT_PUBLIC_KEY`: Your Clerk JWT public key (optional)
- `CLERK_JWT_TEMPLATE`: The JWT template to use (default: "api_auth")
- `MOCK_AUTH_SERVICE`: Set to "true" to use mock authentication (for development and testing)

## Implementation Details

### Gateway Interface

The Clerk gateway implements the `ClerkGateway` interface defined in the domain layer:

```go
type ClerkGateway interface {
    VerifySession(ctx context.Context, sessionToken string) (*model.User, error)
    GetUser(ctx context.Context, clerkUserID string) (*model.User, error)
}
```

### Session Verification

The gateway verifies session tokens by calling the Clerk API's `/v1/sessions/verify` endpoint. This ensures that:

1. The token is valid and has not been tampered with
2. The token has not expired
3. The user associated with the token exists

### User Information Retrieval

The gateway retrieves user information by calling the Clerk API's `/v1/users/{user_id}` endpoint. This provides:

1. User ID
2. Email address
3. Name
4. Creation and update timestamps

### Caching

To improve performance and reduce API calls, the gateway caches user information in memory. This cache:

1. Is keyed by user ID
2. Is thread-safe using a read-write mutex
3. Persists for the lifetime of the application

### Mock Mode

For development and testing, the gateway supports a mock mode that:

1. Bypasses actual API calls to Clerk
2. Returns predefined user information
3. Can be enabled by setting `MOCK_AUTH_SERVICE=true`

## Usage

### Setting Up Clerk

1. Sign up for a Clerk account at [clerk.com](https://clerk.com/)
2. Create a new application
3. Configure your application settings
4. Get your API keys from the Clerk dashboard

### Configuring the Application

1. Set the `CLERK_SECRET_KEY` environment variable to your Clerk secret key
2. Set `AUTH_PROVIDER` to "clerk"

### Verifying a Session

```go
// Get the session token from the request
sessionToken := req.Header.Get("Authorization")

// Verify the session
user, err := clerkGateway.VerifySession(ctx, sessionToken)
if err != nil {
    // Handle authentication error
    return err
}

// Use the authenticated user
fmt.Printf("Authenticated user: %s (%s)\n", user.Name, user.Email)
```

### Getting User Information

```go
// Get user by ID
user, err := clerkGateway.GetUser(ctx, "user_123")
if err != nil {
    // Handle error
    return err
}

// Use the user information
fmt.Printf("User: %s (%s)\n", user.Name, user.Email)
```

## Troubleshooting

### Authentication Issues

If authentication is failing, check:

1. The `CLERK_SECRET_KEY` environment variable is set correctly
2. The session token is being passed correctly in the request
3. The user exists in Clerk
4. The session has not expired

### API Rate Limiting

Clerk has rate limits on API calls. If you're hitting these limits:

1. Implement additional caching
2. Reduce the number of unnecessary API calls
3. Contact Clerk to increase your rate limits

## Limitations

- The current implementation does not support JWT verification using the public key
- User roles and permissions are not currently implemented
- The cache does not expire entries, which could lead to stale data if user information changes

## References

- [Clerk Documentation](https://clerk.com/docs)
- [Clerk API Reference](https://clerk.com/docs/reference/backend-api)

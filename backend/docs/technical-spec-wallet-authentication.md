# Wallet Authentication Technical Specification

## Introduction

This technical specification details the implementation of the Hybrid Authentication approach for the crypto bot application. This approach combines Clerk for user authentication with secure backend storage of exchange API credentials.

## System Components

### 1. Database Components

#### 1.1 API Credentials Table

Stores encrypted exchange API credentials for each user.

```go
// APICredentialEntity represents the database model for API credentials
type APICredentialEntity struct {
    ID         string    `gorm:"primaryKey;type:varchar(50)"`
    UserID     string    `gorm:"not null;index;type:varchar(50)"`
    Exchange   string    `gorm:"not null;index;type:varchar(20)"`
    APIKey     string    `gorm:"not null;type:varchar(100)"`
    APISecret  []byte    `gorm:"not null;type:blob"`  // Encrypted
    Label      string    `gorm:"type:varchar(50)"`
    CreatedAt  time.Time `gorm:"autoCreateTime"`
    UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}
```

#### 1.2 Balance Entities Table

Stores wallet balance information.

```go
// BalanceEntity represents the database model for a balance
type BalanceEntity struct {
    ID        uint      `gorm:"primaryKey"`
    WalletID  uint      `gorm:"not null;index"`
    Asset     string    `gorm:"size:20;not null"`
    Free      float64   `gorm:"type:decimal(18,8);not null"`
    Locked    float64   `gorm:"type:decimal(18,8);not null"`
    Total     float64   `gorm:"type:decimal(18,8);not null"`
    USDValue  float64   `gorm:"type:decimal(18,8);not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
```

### 2. Backend Components

#### 2.1 Clerk Authentication Middleware

```go
// ClerkMiddleware validates Clerk JWT tokens
// NOTE: Uses context keys as defined in context_keys.go to avoid collisions and ensure type safety.
type ClerkMiddleware struct {
    logger *zerolog.Logger
    config *config.Config
}

// Middleware returns a middleware function that validates Clerk JWT tokens
func (m *ClerkMiddleware) Middleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract JWT from Authorization header
            // Validate JWT (return 401 with generic error if invalid)
            // Extract user ID from JWT
            // Add user ID to request context using context_keys.go
            // Call next handler
        })
    }
}
```

#### 2.2 API Credential Repository

```go
// APICredentialRepository handles storage and retrieval of API credentials
// NOTE: All methods use context.Context for cancellation, tracing, and audit logging.
type APICredentialRepository interface {
    Save(ctx context.Context, credential *model.APICredential) error
    GetByUserIDAndExchange(ctx context.Context, userID, exchange string) (*model.APICredential, error)
    DeleteByID(ctx context.Context, id string) error
    ListByUserID(ctx context.Context, userID string) ([]*model.APICredential, error)
}
```

#### 2.3 Encryption Service

```go
// EncryptionService handles encryption and decryption of sensitive data
// NOTE: Keys are loaded from environment variables at startup and never hot-reloaded. Key rotation is supported; existing secrets must be re-encrypted after rotation.
type EncryptionService interface {
    Encrypt(plaintext string) ([]byte, error)
    Decrypt(ciphertext []byte) (string, error)
}
```

#### 2.4 API Credential Handler

```go
// APICredentialHandler handles API credential-related endpoints
// NOTE: API secrets are never returned in plaintext in any API response, even to the owner. All credential operations are logged for audit. Add "last used" timestamp for UX/debugging.
type APICredentialHandler struct {
    useCase    APICredentialUseCase
    encryption EncryptionService
    logger     *zerolog.Logger
}
```

#### 2.5 Updated Account Handler

```go
// AccountHandler handles account-related endpoints
// NOTE: Handles Clerk session expiration, missing/invalid credentials, and generic error responses. All sensitive operations are logged for audit.
type AccountHandler struct {
    useCase    AccountUseCase
    credRepo   APICredentialRepository
    encryption EncryptionService
    logger     *zerolog.Logger
}

// GetWallet returns the user's wallet
func (h *AccountHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
    // Extract user ID from context (using context_keys.go)
    // Get API credentials for user
    // Decrypt API secret (handle decryption/key errors gracefully)
    // Get wallet using credentials
    // Return wallet data (never expose secrets)
    // Log all access for audit
}
```

### 3. Frontend Components

#### 3.1 Clerk Authentication Integration

```typescript
// auth.ts
import { ClerkProvider, SignedIn, SignedOut, UserButton } from '@clerk/nextjs';

export function AuthProvider({ children }: { children: React.ReactNode }) {
  // Handles Clerk session expiration and re-authentication prompts
  return (
    <ClerkProvider>
      {children}
    </ClerkProvider>
  );
}

export function AuthGuard({ children }: { children: React.ReactNode }) {
  // Displays prompt if session is expired or missing
  return (
    <>
      <SignedIn>{children}</SignedIn>
      <SignedOut>
        <div>Please sign in to access this page</div>
      </SignedOut>
    </>
  );
}
```

#### 3.2 API Credential Management UI

```typescript
// components/ApiCredentialForm.tsx
export function ApiCredentialForm() {
  // Form state and handlers
  // Submit function that calls API
  // Never display or return API secrets in plaintext
  // Add feedback for invalid/expired credentials
  return (
    <form>
      <input type="text" placeholder="Label" />
      <input type="text" placeholder="API Key" />
      <input type="password" placeholder="API Secret" />
      <button type="submit">Save</button>
    </form>
  );
}
```

#### 3.3 Wallet Display Component

```typescript
// components/WalletDisplay.tsx
export function WalletDisplay() {
  // Fetch wallet data from API
  // Display loading state
  // Display error state (generic, no sensitive info)
  // Display wallet data
  // Show call-to-action if credentials are missing/invalid
  return (
    <div>
      <h2>Your Wallet</h2>
      <table>
        <thead>
          <tr>
            <th>Asset</th>
            <th>Free</th>
            <th>Locked</th>
            <th>Total</th>
            <th>USD Value</th>
          </tr>
        </thead>
        <tbody>
          {/* Map wallet data to rows */}
        </tbody>
      </table>
    </div>
  );
}
```

## API Endpoints

### Authentication Endpoints

- `POST /api/v1/auth/verify` - Verify Clerk JWT token

### API Credential Endpoints

- `POST /api/v1/credentials` - Create new API credential
- `GET /api/v1/credentials` - List user's API credentials
- `GET /api/v1/credentials/{id}` - Get specific API credential
- `PUT /api/v1/credentials/{id}` - Update API credential
- `DELETE /api/v1/credentials/{id}` - Delete API credential

### Wallet Endpoints

- `GET /api/v1/account/wallet` - Get user's wallet
- `POST /api/v1/account/refresh` - Refresh wallet data
- `GET /api/v1/account/balance/{asset}` - Get balance history for specific asset

## Security Considerations

### Encryption

API secrets will be encrypted using AES-256-GCM with the following approach:

1. Generate a random nonce for each encryption operation
2. Use AES-256-GCM for authenticated encryption
3. Store the nonce with the ciphertext
4. Use a master encryption key stored in environment variables

```go
func Encrypt(plaintext string, key []byte) ([]byte, error) {
    // Create cipher
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    // Create GCM
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    // Create nonce
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }
    
    // Encrypt
    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    
    return ciphertext, nil
}
```

### JWT Validation

Clerk JWT tokens will be validated using the Clerk SDK:

1. Extract JWT from Authorization header
2. Validate JWT signature using Clerk public key
3. Validate JWT claims (expiration, issuer, etc.)
4. Extract user ID from JWT

## Error Handling

The system will handle the following error scenarios:

1. Invalid/expired JWT token
2. Missing API credentials
3. Invalid API credentials
4. Exchange API errors (rate limiting, downtime, etc.)
5. Database errors

Each error will have a specific error code and user-friendly message.

## Performance Considerations

1. Implement caching for wallet data (TTL: 5 minutes)
2. Use connection pooling for database connections
3. Implement rate limiting for API endpoints
4. Use background jobs for wallet refresh operations

## Testing Strategy

1. Unit tests for each component
2. Integration tests for API endpoints
3. End-to-end tests for critical flows
4. Security testing (penetration testing, code review)
5. Performance testing

## Deployment Strategy

1. Database migrations will be run before code deployment
2. Deploy backend changes first
3. Deploy frontend changes after backend is stable
4. Monitor for errors and performance issues
5. Rollback plan in case of critical issues

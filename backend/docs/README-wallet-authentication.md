# Wallet Authentication Implementation

## Overview

This project implements a secure wallet authentication system for our cryptocurrency trading bot. It allows users to securely store their exchange API credentials and access their real wallet balances.

## Documentation

- [Implementation Plan](implementation-plan-wallet-authentication.md) - High-level overview of the implementation plan
- [Technical Specification](technical-spec-wallet-authentication.md) - Detailed technical specification
- [Tasks](../tasks/wallet-authentication-tasks.md) - Breakdown of implementation tasks

## Getting Started

### Prerequisites

- Go 1.18+
- Node.js 16+
- PostgreSQL 13+
- Clerk account (for authentication)

### Environment Variables

```
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=crypto_bot

# Encryption
MEXC_CRED_ENCRYPTION_KEY=your-encryption-key

# Clerk
CLERK_PUBLISHABLE_KEY=your-clerk-publishable-key
CLERK_SECRET_KEY=your-clerk-secret-key

# MEXC API (for testing)
MEXC_API_KEY=your-mexc-api-key
MEXC_SECRET_KEY=your-mexc-secret-key
```

### Running Migrations

```bash
go run cmd/migrate/main.go up
```

### Running the Backend

```bash
go run cmd/server/main.go
```

### Running the Frontend

```bash
cd frontend
npm install
npm run dev
```

## Architecture

### Backend

The backend is built using Go with the following components:

- **HTTP Server**: Chi router for handling HTTP requests
- **Database**: PostgreSQL with GORM for ORM
- **Authentication**: Clerk for user authentication
- **Encryption**: AES-256-GCM for encrypting API secrets
- **API Clients**: Clients for interacting with cryptocurrency exchanges

### Frontend

The frontend is built using React with the following components:

- **Authentication**: Clerk for user authentication
- **State Management**: React Query for data fetching and caching
- **UI Components**: Custom components for wallet display and credential management
- **Routing**: React Router for navigation

## Security Considerations

- API secrets are encrypted at rest using AES-256-GCM
- Encryption keys are loaded from environment variables at startup and never hot-reloaded
- Key rotation policy is established; existing secrets are re-encrypted after rotation
- API secrets are never exposed in API responses or logs
- All credential operations (create/update/delete, failed access) are logged for audit purposes
- JWT tokens are validated for every request
- Rate limiting is implemented to prevent abuse
- Input validation is performed for all user inputs
- Error messages are generic to prevent information leakage
- Audit logs are reviewed regularly and available to users (upon request) for their own actions

## Testing

### Running Backend Tests

```bash
go test ./...
```

### Running Frontend Tests

```bash
cd frontend
npm test
```

## Troubleshooting

- **Invalid Clerk JWT:** Ensure your session is active and your token is valid. Try re-authenticating.
- **Missing Encryption Key:** Check that `MEXC_CRED_ENCRYPTION_KEY` is set in your environment.
- **API Credential Not Found:** Add or update your API credentials via the UI.
- **Credential Errors:** Invalid/expired credentials will result in generic error messages. Check and update your credentials.
- **Migration Issues:** Use transactional/dry-run migrations and prepare a rollback plan.

## Deployment

See the [deployment plan](../tasks/wallet-authentication-tasks.md#phase-4-testing-and-deployment) for detailed deployment instructions.

## Contributing

1. Create a new branch for your feature
2. Implement your changes
3. Write tests for your changes (including security, error boundaries, and audit logging)
4. Follow security best practices, especially for authentication and credential code
5. Submit a pull request

## Compliance & Monitoring

- If your product is used in regulated jurisdictions, ensure compliance with relevant standards (e.g., GDPR for user data).
- Monitor failed authentication attempts, credential access patterns, and suspicious activity.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

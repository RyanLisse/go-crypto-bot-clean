# Technology Context

## Technology Stack

### Frontend
- **Framework**: React 18 with TypeScript
- **Build Tool**: Vite
- **State Management**: Redux Toolkit + TanStack Query
- **UI Components**:
  - Brutalist Design System
  - Radix UI primitives
  - Tailwind CSS
  - shadcn/ui components
- **Typography**: JetBrains Mono (monospace)
- **Authentication**: Clerk
- **Local Database**: TursoDB with Drizzle ORM
- **Testing**: Vitest + Playwright

### Backend
- **Language**: Go 1.24+
- **Database**:
  - TursoDB (primary)
  - SQLite (fallback/development)
- **API Framework**: 
  - Chi Router (primary routing)
  - Huma (OpenAPI documentation and service integration)
- **Authentication**: Clerk SDK
- **WebSocket**: Gorilla WebSocket
- **Documentation**: OpenAPI/Swagger (via Huma)

### Database
- **Primary**: TursoDB (distributed SQLite)
- **ORM**:
  - Backend: Native SQL with prepared statements
  - Frontend: Drizzle ORM
- **Migration**: Drizzle Kit
- **Sync**: TursoDB built-in sync

### External Services
- MEXC Exchange API
- Clerk Authentication
- TursoDB Cloud

### Development Tools
- Bun for package management and running
- Docker for containerization
- GitHub Actions for CI/CD
- ESLint + Prettier for code formatting
- Husky for git hooks

## Implementation Documentation

### Frontend Architecture
- Brutalist design system implementation
- Component structure and organization
- State management patterns
- Data fetching and caching with React Query
- Authentication flow with Clerk
- Responsive layout with brutalist principles
- Monospace typography and high-contrast UI
- Offline capabilities
- Testing strategy

### Backend Architecture
- Clean architecture implementation
- Service layer organization
- Repository pattern usage
- WebSocket implementation
- Error handling
- Logging and monitoring

### Database Layer
- TursoDB setup and configuration
- Migration management
- Data synchronization
- Performance optimization
- Backup and recovery

### API Layer
- REST endpoints
- WebSocket handlers
- Authentication middleware
- Rate limiting
- Error handling
- API documentation

## Key Libraries and Dependencies

### Frontend Core
- @clerk/clerk-react - Authentication
- @reduxjs/toolkit - State management
- @tanstack/react-query - Data fetching
- @radix-ui/* - UI primitives
- class-variance-authority - Component styling
- clsx - Class name utilities
- tailwind-merge - Tailwind class merging
- lucide-react - Brutalist-friendly icons
- drizzle-orm - Local database ORM
- @libsql/client - TursoDB client

### Frontend Development
- vite - Build tool
- typescript - Type checking
- vitest - Unit testing
- playwright - E2E testing
- eslint - Linting
- prettier - Code formatting
- tailwindcss - Utility CSS

### Backend Core
- github.com/go-chi/chi/v5 - Primary HTTP router
- github.com/danielgtaylor/huma/v2 - OpenAPI documentation and service integration
- github.com/clerk/clerk-sdk-go - Authentication
- github.com/gorilla/websocket - WebSocket
- github.com/mattn/go-sqlite3 - SQLite driver
- github.com/sirupsen/logrus - Logging

### Backend Development
- github.com/stretchr/testify - Testing
- github.com/spf13/viper - Configuration
- github.com/spf13/cobra - CLI tools

## Build & Deployment

### Frontend
- Vite for development and production builds
- Docker for containerization
- Environment-based configuration
- Automated testing in CI

### Backend
- Go modules for dependency management
- Docker multi-stage builds
- Configuration via environment variables
- Automated testing and linting

### Database
- TursoDB cloud deployment
- Local development with SQLite
- Automated backups
- Monitoring and maintenance

## Core Technologies
- Go 1.22+ (Primary language)
- Cobra (CLI framework)
- Standard library packages:
  - `archive/tar` for archive creation
  - `compress/gzip` for compression
  - `path/filepath` for path manipulation
  - `os` for file system operations
  - `crypto/sha256` for checksums
  - `encoding/json` for metadata serialization
  - `time` for timestamps and retention
  - `testing` for test framework

## Development Tools
- Go modules for dependency management
- `golangci-lint` for code quality
- `go test` for testing
- Git for version control

## Project Structure
```
backend/
  ├── cmd/
  │   └── cli/
  │       └── commands/
  │           ├── backup.go       # Backup command implementation
  │           └── backup_test.go  # Backup command tests
  └── pkg/
      └── backup/
          ├── service.go          # Backup service implementation
          └── types.go           # Shared types and interfaces
```

## Dependencies
- Direct dependencies:
  - `github.com/spf13/cobra` - CLI framework
  - `github.com/spf13/viper` - Configuration management
  - Standard library packages only for core functionality

## Technical Constraints
- Cross-platform compatibility required
- Minimal external dependencies
- Standard Go idioms and patterns
- Error handling through explicit returns
- Comprehensive test coverage

## Security Requirements
- Input validation for all user inputs
- Safe file system operations
- Future: Encryption for sensitive data
- Future: Secure credential storage

## Performance Considerations
- Efficient file system traversal
- Compression optimization
- Memory usage optimization for large files
- Concurrent operations where appropriate

## Testing Strategy
- Unit tests for all packages
- Integration tests for file system operations
- Table-driven test patterns
- Test coverage requirements
- Temporary test directories
- Cleanup of test artifacts

## Monitoring and Logging
- Structured logging
- Operation progress tracking
- Error tracking and reporting
- Future: Metrics collection

## Future Technical Considerations
- Encryption implementation
- Remote storage integration
- Automated scheduling
- Performance monitoring
- Cloud service integration

### API Integration (June 2024)
- All frontend API calls are now standardized to use `VITE_API_URL` (default: `http://localhost:8080/api/v1`)
- Authentication endpoints use `/api/v1/auth/*` for both Clerk and JWT flows
- CORS and environment configuration are unified for local development and production

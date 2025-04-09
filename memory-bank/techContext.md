# Technology Context

## Technology Stack

### Frontend
- **Framework**: React 18 with TypeScript
- **Build Tool**: Vite
- **State Management**: Redux Toolkit + TanStack Query
- **UI Components**: 
  - Material-UI (MUI)
  - Radix UI primitives
  - Tailwind CSS
- **Authentication**: Clerk
- **Local Database**: TursoDB with Drizzle ORM
- **Testing**: Vitest + Playwright

### Backend
- **Language**: Go 1.24+
- **Database**: 
  - TursoDB (primary)
  - SQLite (fallback/development)
- **API Framework**: Chi Router
- **Authentication**: Clerk SDK
- **WebSocket**: Gorilla WebSocket
- **Documentation**: OpenAPI/Swagger

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
- Component structure and organization
- State management patterns
- Data fetching and caching
- Authentication flow
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
- @mui/material - UI components
- @radix-ui/* - UI primitives
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
- github.com/go-chi/chi/v5 - HTTP router
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

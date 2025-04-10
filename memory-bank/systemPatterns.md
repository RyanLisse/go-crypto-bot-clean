# System Patterns

## Core Architecture: Modern Full-Stack with Clean Architecture

The system follows a modern full-stack architecture with clear separation between frontend and backend components:

1. **Frontend (React + TypeScript)**:
   - Built with Vite and React 18
   - TypeScript for type safety
   - Redux Toolkit for state management
   - TanStack Query for data fetching
   - Clerk for authentication
   - Drizzle ORM with TursoDB for local data
   - Material-UI and Radix UI for components

2. **Backend (Go)**:
   - Hexagonal architecture principles
   - Clean separation of concerns
   - Domain-driven design
   - Interface-based design
   - Dependency injection

3. **Database Layer**:
   - TursoDB as primary database
   - Drizzle ORM for frontend local data
   - SQLite as fallback/development option
   - Shadow mode for gradual migration

4. **Authentication & Authorization**:
   - Clerk for user management and authentication
   - JWT tokens for API authentication
   - Role-based access control
   - Secure session management

## Frontend Architecture

1. **Component Structure**:
   - Feature-based organization
   - Atomic design principles
   - Shared components library
   - Custom hooks for reusable logic

2. **State Management**:
   - Redux Toolkit for global state
   - React Query for server state
   - Local state with useState/useReducer
   - Context API for theme/auth

3. **Data Layer**:
   - TanStack Query for API integration
   - Drizzle ORM for local data
   - Optimistic updates
   - Offline-first capabilities

4. **Styling**:
   - Material-UI components
   - Radix UI primitives
   - Tailwind CSS for custom styling
   - CSS-in-JS with Emotion

## Backend Architecture

1. **Core Domain**:
   - Business rules and models (`internal/domain`)
   - Service interfaces as ports (`internal/domain/service`)
   - Clean architecture principles

2. **API Layer**:
   - Primary routing with Chi Router
   - OpenAPI/Swagger documentation via Huma integration
   - Service-based API documentation
   - RESTful endpoints
   - WebSocket support
   - Rate limiting and security

3. **Router Integration**:
   - Chi Router for main HTTP routing and middleware
   - Huma for OpenAPI documentation and service integration
   - Conditional Huma setup based on service availability
   - Adapter pattern for legacy Gin handlers

4. **Database Layer**:
   - TursoDB integration
   - Repository pattern
   - Migration management
   - Data synchronization

5. **External Integrations**:
   - MEXC Exchange API
   - Clerk authentication
   - WebSocket for real-time data
   - Notification services

## Data Persistence Strategy

1. **TursoDB Integration**:
   - Distributed SQLite database
   - Edge-deployed instances
   - Automatic synchronization
   - Offline-first capabilities

2. **Migration Strategy**:
   - Shadow mode deployment
   - Gradual transition from SQLite
   - Data validation and verification
   - Rollback capabilities

3. **Synchronization**:
   - Automatic sync with cloud
   - Conflict resolution
   - Offline support
   - Real-time updates

4. **Performance Optimization**:
   - Query optimization
   - Connection pooling
   - Caching strategies
   - Batch operations

## Cross-Cutting Concerns

1. **Authentication & Security**:
   - Clerk for user management
   - JWT token validation
   - CORS configuration
   - Rate limiting

2. **Error Handling**:
   - Consistent error responses
   - Error tracking and logging
   - Graceful degradation
   - User-friendly messages

3. **Logging & Monitoring**:
   - Structured logging
   - Performance metrics
   - Error tracking
   - System health monitoring

4. **Testing Strategy**:
   - Unit tests with Vitest
   - Integration tests
   - E2E tests with Playwright
   - Performance testing

## Real-Time Features

1. **WebSocket Integration**:
   - Market data streaming
   - Portfolio updates
   - Trade notifications
   - Connection management

2. **State Synchronization**:
   - Real-time UI updates
   - Optimistic updates
   - Conflict resolution
   - Offline support

## Development Workflow

1. **Local Development**:
   - Vite dev server
   - Hot module replacement
   - Local TursoDB instance
   - Development tools

2. **Testing Environment**:
   - Automated testing
   - CI/CD pipeline
   - Staging environment
   - Performance monitoring

3. **Production Deployment**:
   - Docker containerization
   - Cloud deployment
   - Monitoring and logging
   - Backup and recovery

## System Architecture

### Command Structure
- Commands follow a consistent pattern using Cobra library
- Each command is self-contained in its own file under `backend/cmd/cli/commands/`
- Commands implement validation and proper error handling
- Common utilities and shared functionality are extracted into separate packages

### Backup System
- Located in `backend/cmd/cli/commands/backup.go`
- Follows command pattern with clear separation of concerns:
  - Command definition and flag handling
  - Option validation
  - Backup execution logic
  - Metadata management
- Supports both full and incremental backup types
- Uses standard library for file system operations
- Implements compression using tar + gzip
- Maintains backup metadata for tracking and verification

### Core Patterns

#### Error Handling
- Functions return errors explicitly
- Errors are wrapped with context using `fmt.Errorf`
- Logging is used for debugging and audit trails

#### Configuration
- Command-line flags for runtime configuration
- Environment variables for sensitive settings
- Configuration validation at startup

#### Testing
- Table-driven tests for command functionality
- Integration tests for file system operations
- Temporary directories for test isolation
- Cleanup of test artifacts

#### Security
- Input validation for all command parameters
- Safe file system operations
- Planned: Encryption for sensitive data
- Planned: Secure storage mechanisms

#### Monitoring
- Progress reporting during operations
- Operation metadata collection
- Error tracking and logging
- Performance metrics collection (planned)

### Module Boundaries
- CLI commands are isolated in `commands` package
- Core functionality is separated into domain-specific packages
- Utility functions are shared via common packages
- External integrations are abstracted behind interfaces

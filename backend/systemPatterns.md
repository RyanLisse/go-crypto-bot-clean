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

## Standardized Repository Pattern

1. **Repository Interface Design**:
   - All interfaces defined in `internal/domain/port`
   - Focus on domain operations, not persistence details
   - Domain models as parameters and return values
   - Context for cancellation and tracing

2. **Base Repository Implementation**:
   - Common base repository for each persistence mechanism
   - Consistent error handling and mapping
   - Shared database connection management
   - Transaction support

3. **Specific Repository Implementations**:
   - Embed base repository
   - Implement domain interfaces
   - Map between entity and domain models
   - Use consistent error handling

4. **Entity-Model Mapping**:
   - Entities in `internal/adapter/repository/{orm}/entity`
   - ORM-specific tags and hooks
   - Bidirectional conversion methods (ToModel/FromModel)
   - Preserve data integrity during conversions

5. **Transaction Management**:
   - Transaction manager implementation using GORM
   - Transaction context propagation
   - Consistent transaction handling across repositories
   - Error handling and rollback

6. **Repository Factory**:
   - Centralized factory for repository creation
   - Dependency injection for database connections
   - Caching of repository instances
   - Type-safe repository access

7. **Mock Repositories**:
   - In-memory implementations for testing
   - Configurable error injection
   - Thread-safe operation
   - Test behavior controllers

## Unified Factory Pattern

1. **AppFactory Structure**:
   - Single entry point for creating all application components
   - Centralized dependency management
   - Composable initialization
   - Type-safe component access

2. **Component Creation Patterns**:
   - Lazy initialization with caching
   - Safe creation with error handling
   - Consistent naming conventions
   - Interface-based returns for abstraction

3. **Mock/Real Implementation Switching**:
   - Configuration-based component selection
   - Environment-aware mock detection
   - Production safeguards
   - Centralized logging of mock usage

4. **Integration with Application Startup**:
   - Simplified bootstrapping process 
   - Controlled component initialization order
   - Graceful error handling
   - Resource cleanup management

## Cross-Cutting Concerns

1. **Authentication & Security**:
   - ConsolidatedAuthMiddleware as the standard authentication middleware
   - JWT token validation with proper error handling
   - Environment-aware test middleware with production safeguards
   - CORS configuration with secure defaults
   - Rate limiting with tiered strategies

2. **Error Handling**:
   - Standardized AppError structure with HTTP status mapping
   - Centralized error middleware for consistent response formatting
   - Type-safe error context passing through request lifecycle
   - Domain-specific error types with validation
   - Consistent error logging and monitoring
   - Clear separation between user-facing and internal error details

3. **Logging & Monitoring**:
   - Structured logging with request tracking
   - Performance metrics collection
   - Error tracking and correlation
   - System health monitoring
   - Context-aware log enrichment

4. **Testing Strategy**:
   - Unit tests with Vitest/Go testing
   - Integration tests with repository mocking
   - E2E tests with Playwright
   - Performance testing
   - Standardized mock implementations for testing

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

## Text Tag Optimization System

### Tag Categories
- `@domain` - Domain models and business logic
- `@repository` - Data access and persistence
- `@service` - Business services and operations
- `@usecase` - Application use cases
- `@handler` - API handlers and controllers
- `@middleware` - HTTP middleware components
- `@util` - Utility functions and helpers
- `@config` - Configuration and settings
- `@factory` - Factory methods and dependency injection
- `@test` - Test-related code and utilities

### Tag Usage Patterns
- File headers include primary tag category
- Interface methods include relevant tags
- Implementation methods reference interface tags
- Complex logic sections include descriptive tags
- Cross-cutting concerns marked with multiple tags

### Documentation Integration
- Tags are included in Go doc comments
- Tags help generate API documentation
- Tags assist in code navigation and search
- Tags provide context for code review

### Testing Strategy
- Test files include `@test` tag with subcategory
- Test cases reference implementation tags
- Mock implementations include both `@test` and original tags
- Integration tests include multiple component tags

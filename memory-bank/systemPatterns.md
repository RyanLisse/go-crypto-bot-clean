# System Patterns

## Core Architecture: Hexagonal (Ports and Adapters)

The system follows hexagonal architecture principles with clear separation between:

1. **Core Domain** (the Hexagon):
   - Business rules and models (`internal/domain`)
   - Service interfaces defined as ports (`internal/domain/service`)

2. **Ports**:
   - Primary (Driving) Ports: API endpoints, CLI commands (interfaces defined in domain layer)
   - Secondary (Driven) Ports: Repository and exchange API interfaces (also in domain layer)

3. **Adapters**:
   - Primary (Driving) Adapters: REST handlers, gRPC services, CLI commands (in `internal/api` and `internal/cli`)
   - Secondary (Driven) Adapters: DB repositories, MEXC API client (in `internal/platform`)

## Dependency Rule
Dependencies flow inwards: cmd -> Application Layer (api/cli) -> Domain Layer <- Platform Layer. The Domain Layer depends on nothing outside itself (except potentially shared utilities or standard library types).

## Key Design Patterns

1. **Repository Pattern**: For data access abstraction
   - Interfaces defined in domain layer
   - Implementations in platform layer
   - Allows swapping database technologies

2. **Service Pattern**: For business logic encapsulation
   - Interfaces defined in domain layer
   - Implementations typically in core layer
   - Service composition for complex operations

3. **Adapter Pattern**: For external system integration
   - MEXC API client as adapter to exchange, implementing the `service.ExchangeService` interface
   - Repository implementations as adapters to database
   - Mock implementations for testing

4. **Dependency Injection**: For wiring components
   - Constructor injection for services and repositories
   - Interface-based dependency injection for flexibility and testability
   - Application bootstrapped in cmd layer

## Module Boundaries

- **Domain Module**: Contains all business entities and interface definitions
- **Service Module**: Implements business logic and orchestration
- **Platform Module**: Provides infrastructure implementations
- **API Module**: Exposes functionality via HTTP/gRPC
- **CLI Module**: Provides command-line interface

## Cross-Cutting Concerns

1. **Logging**: Centralized via dependency injection
2. **Error Handling**: Domain-specific errors defined in domain layer
3. **Configuration**: Environment-based with sensible defaults
4. **Validation**: Input validation at API boundaries, domain validation in services
5. **Testing**: Interface-based mocking for unit tests

## Interface-Based Design

The system follows a strict interface-based design approach for better testability and flexibility:

1. **Service Interfaces**: All services implement interfaces defined in the domain layer
   - `service.ExchangeService` for exchange API operations
   - `NewCoinWatcher` for new coin detection functionality
   - `account.AccountService` for account management operations
   - `position.PositionService` for position management operations
   - Services depend on interfaces, not concrete implementations

2. **Repository Interfaces**: Data access is abstracted through repository interfaces
   - `repository.NewCoinRepository` for new coin persistence operations
   - `position.PositionRepository` for position data operations
   - `account.BoughtCoinRepository` for tracking purchased coins
   - Clean separation between domain logic and data access

3. **Mock Implementations**: Test doubles implementing service interfaces
   - `MockMEXCClient` for testing without real API calls
   - `MockPositionRepository` for testing position operations
   - `MockMarketService` for testing market data operations
   - Mock repositories for testing without database dependencies

4. **Type Assertions**: Careful type assertions when working with interface types
   - Helper functions for common interface operations
   - Error handling for type assertion failures

## Core Services

1. **Account Service**: Manages account-related operations
   - Retrieves account balances from the exchange
   - Calculates portfolio value across all assets
   - Assesses position risk levels for individual symbols
   - Validates API keys with the exchange
   - Follows interface-based design with dependency injection

2. **Position Service**: Manages trading positions
   - Opens and closes trading positions
   - Sets stop-loss and take-profit levels
   - Calculates current profit/loss for positions
   - Retrieves positions by ID and gets all positions
   - Uses repository pattern for data persistence

3. **WebSocket Client**: Manages real-time market data streaming
   - Implements robust reconnection logic with exponential backoff
   - Uses token bucket rate limiting for connections and subscriptions
   - Ensures thread safety with mutex-based synchronization
   - Provides context-aware operations for cancellation and timeouts
   - Follows the observer pattern for ticker updates via channels
   - Implements automatic resubscription after reconnection
   - Integrates with market service for current price data

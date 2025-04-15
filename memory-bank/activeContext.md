# Active Context

## Recent Major Feature Completion: New Coin Detection and AutoBuy (Event-Driven)

- The "New Coin Detection and AutoBuy" feature is now fully implemented and integrated using a robust event-driven architecture.
- All four vertical slices have been completed using TDD and vertical slicing:
  1. Detection & Event: NewCoinUsecase detects new tradable coins, updates repository, and publishes domain events.
  2. Autobuy Service Logic: AutobuyService listens for events, loads config, validates, calculates, and triggers buy logic.
  3. Trade Usecase Enhancement: TradeUsecase now supports ExecuteMarketBuyOrder with SL/TP, integrating with PositionUsecase.
  4. Integration: Event bus, NewCoinUsecase, AutobuyService, and TradeUsecase are fully wired in main.go; end-to-end integration tests pass.
- The system now supports reliable, automated detection and execution of new coin listings on MEXC, with configurable risk and execution parameters.
- Comprehensive unit and integration test coverage is in place for all components and flows.
- The event-driven flow is:
  Detection → Event → Autobuy → Trade Execution → Position Management
- This marks a major milestone in automated trading capabilities and system modularity.

# Active Context

## Active Development Focus

### Current Task: Task 7 - Risk Management System Implementation (in progress)

We are currently focused on implementing the Risk Management System, which is a critical component for ensuring trading safety and compliance with user risk profiles. We've made progress on this task:

1. ✅ Subtask 7.1: Implemented Risk Control Models and Core Domain Logic
   - Created comprehensive risk controls:
     - ConcentrationControl: Prevents over-concentration in a single asset
     - LiquidityControl: Ensures trading in markets with sufficient volume
     - ExposureControl: Limits total market exposure based on risk profile
     - DrawdownControl: Monitors and limits portfolio drawdown
     - VolatilityControl: Evaluates market volatility before trading
     - PositionSizeControl: Enforces proper position sizing based on portfolio
   - Implemented RiskEvaluator to coordinate multiple risk controls
   - Developed BaseRiskControl providing common functionality
   - Created relevant domain models and interfaces

Next steps:
1. Begin work on Subtask 7.2: Develop Risk Management Repository and Persistence Layer
2. Design database schema for storing risk assessments, profiles, and constraints
3. Implement repositories for risk-related entities
4. Integrate risk management with the trading and position management systems

### Implementation Details and Decisions

- Risk controls follow a consistent interface (RiskControl) allowing easy composition and evaluation
- The system supports different risk profiles (Conservative, Moderate, Aggressive) with appropriate thresholds
- Risk evaluations produce detailed assessments with recommendations for the user
- Domain events are used to notify other system components about risk violations
- The risk system is designed to be extensible, allowing new risk controls to be easily added

### Technical Constraints and Considerations

- Risk evaluations must be performed efficiently to not slow down trading operations
- The system needs to handle real-time market data for accurate risk assessment
- Risk profiles and constraints must be persisted and easily configurable
- Integration with position and trade systems must maintain transactional integrity

### Blockers/Dependencies

- Need to design database schema for risk-related entities
- Integration with user configuration system for risk profiles

### Notes

- The implementation follows our established architecture pattern with clear separation of concerns
- Risk control models are thoroughly tested with unit tests
- The system is designed to be proactive, preventing risky trades before execution

## General Project Context
- The backend implementation now includes:
  - Complete MEXC API integration for market data and trading
  - Position management system with database persistence
  - Trade execution system with order management and persistence
  - HTTP API endpoints for all main functionality
  - Factory pattern implementation for proper dependency injection
  - Core risk management domain logic and controls

## Implementation Notes
- The Position Management System includes:
  - A comprehensive domain model with Position, PositionSide, and PositionStatus types
  - Repository implementations for all persistence operations
  - Use cases for creating, reading, updating, and closing positions
  - Service layer for business logic including position performance tracking
  - HTTP handlers implementing RESTful API endpoints for all position operations
  
- The Trade Execution System includes:
  - Order domain model with OrderSide, OrderType, OrderStatus, and TimeInForce enums
  - GORM-based OrderRepository for persistence with full CRUD operations
  - TradeUseCase implementation for business logic
  - TradeService for integration with MEXC exchange API
  - HTTP handlers for RESTful API endpoints supporting order operations
  - TradeFactory for properly managing dependencies and component creation

- The HTTP API endpoints for positions follow RESTful principles with these key routes:
  - POST /positions: Create a new position
  - GET /positions: List positions with filtering options
  - GET /positions/open: Get open positions
  - GET /positions/:id: Get a specific position
  - PUT /positions/:id: Update a position
  - PUT /positions/:id/close: Close a position
  - DELETE /positions/:id: Delete a position

- The HTTP API endpoints for trading include:
  - POST /trades/orders: Place a new order
  - DELETE /trades/orders/:symbol/:orderId: Cancel an order
  - GET /trades/orders/:symbol/:orderId: Get order status
  - GET /trades/orders/open/:symbol: Get open orders
  - GET /trades/orders/history/:symbol: Get order history
  - GET /trades/calculate/:symbol: Calculate required quantity

- All implementations include comprehensive unit tests

## Current Focus: Risk Management System - Repository and Persistence Layer (Task 7.2)

We are currently implementing the Risk Management System, focusing on the repository and persistence layer (Task 7.2). This follows the completion of the Risk Control Models and Core Domain Logic (Task 7.1).

### Completed
- ✅ Risk Control Models (Concentration, Liquidity, Exposure)
- ✅ Base risk control structure and interfaces
- ✅ Risk evaluation logic for various controls
- ✅ Risk profile definitions and constraints

### Current Task (7.2): Risk Management Repository and Persistence
We need to implement the database layer for storing and retrieving risk-related data. This includes:

1. Design database schema for:
   - Risk assessments
   - Risk profiles
   - Risk constraints
   - Risk metrics

2. Implement the following repositories:
   - RiskAssessmentRepository
   - RiskMetricsRepository
   - RiskConstraintRepository
   - RiskProfileRepository

3. Create database migrations with:
   - Proper table relationships
   - Appropriate indices for efficient queries
   - Constraints to maintain data integrity

4. Implement GORM-based persistence layer:
   - Define GORM models and tags
   - Map domain models to database entities
   - Implement efficient query methods
   - Ensure proper transaction handling

### Implementation Approach
- Follow Clean Architecture principles
- Separate domain models from database entities
- Use repository interfaces defined in the domain layer
- Implement repositories in the infrastructure layer using GORM
- Create comprehensive unit tests for the repositories

### Integration Points
- The repositories will be used by the Risk Use Case (to be implemented in Task 7.3)
- Trade Execution System will use the Risk Management System to validate trades
- Position Management System will trigger risk evaluations on position changes

### Next Steps
1. Design the database schema
2. Implement the repository interfaces
3. Create the database migrations
4. Implement the GORM-based repositories
5. Write tests for the repositories

### Technical Considerations
- Ensure proper indexing for performance optimization
- Consider caching strategies for frequently accessed risk data
- Implement proper error handling and logging
- Consider versioning for risk profiles and assessments

## Integration Milestone: Frontend-Backend API Unification (June 2024)

- The frontend and backend are now fully integrated via a unified API URL (`VITE_API_URL`, default: `http://localhost:8080/api/v1`).
- All frontend service and repository calls use the standardized API URL and `/api/v1/auth/*` endpoints for authentication.
- CORS and environment configuration are aligned for both local development and production.
- Service-level and repository-level integration tests pass; component and E2E tests require environment fixes (e.g., missing dependencies, DOM setup).
- Next step: Ensure all tests pass, fix test environment issues, and document any remaining blockers.

# Active Context Update (2024-06-10)

- Subtask 14.2 (Migrate project structure and configuration files) is complete.
  - Project structure now matches migration guide: src/app, src/components/ui, src/lib, src/test, public, etc.
  - All config files (tsconfig.json, eslint.config.mjs, postcss.config.mjs, vitest.config.ts, tailwind.config.ts) are present and correct.
  - Linter errors resolved; only a harmless warning remains.
  - Tailwind is fully enabled and configured for Next.js + Bun + shadcn/ui.
- Subtask 14.3 (Implement file-based routing system) is complete.
  - (dashboard) route group and placeholder pages for all main sections are in place.
  - Dashboard layout file created for shared structure.
  - Structure matches migration guide and is ready for UI component migration.
- **Next focus:** Subtask 14.4 — Migrate UI components to Next.js.
  - Next.js Link-based navigation added to dashboard layout.
  - Global Toaster (Sonner) added to root layout for notifications.
  - All dashboard pages are now accessible via Next.js routing.
  - Continue reviewing and adapting remaining UI components for Next.js compatibility and best practices.

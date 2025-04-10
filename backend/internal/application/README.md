# Application Layer

This directory contains the application layer of our trading system, implementing use cases by orchestrating domain objects and defining service interfaces.

## Structure

```
application/
├── services/     # Service interfaces
│   ├── trading_service.go
│   ├── market_data_service.go
│   └── risk_service.go
└── usecases/     # Use case implementations
    ├── trading/
    ├── market_data/
    └── risk/
```

## Services

The services directory contains interface definitions for our core application services:

- **TradingService**: Handles order and position management
- **MarketDataService**: Provides market data and price information
- **RiskService**: Manages risk parameters and position risk assessment

## Use Cases

The usecases directory contains the concrete implementations of our application services, organized by domain:

- **Trading**: Order placement, position management
- **Market Data**: Price updates, historical data retrieval
- **Risk**: Risk calculation, alert management

## Design Principles

1. **Use Case Isolation**: Each use case is independent and focused on a specific business operation
2. **Domain Model Usage**: Use cases work with domain models and implement business rules
3. **Port Usage**: Use cases interact with the outside world through ports (interfaces)
4. **Dependency Inversion**: Dependencies point inward, with interfaces defined in the domain layer

## Dependencies

The application layer:
- Depends on the domain layer
- Has no knowledge of infrastructure implementations
- Defines its own interfaces for external services

## Testing

Use cases should be tested with:
- Mock implementations of ports
- Stub domain models
- Focus on business logic and orchestration 
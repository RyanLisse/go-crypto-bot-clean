# Domain Layer

This directory contains the core domain logic of our trading system, following the hexagonal architecture pattern (also known as ports and adapters).

## Structure

```
domain/
├── models/       # Core domain entities
│   ├── trade.go
│   ├── order.go
│   └── position.go
└── ports/        # Interface definitions (ports)
    └── repositories.go
```

## Domain Models

The domain models represent the core business entities in our system:

- **Trade**: Represents a completed trade with price, amount, and timestamp
- **Order**: Represents a trading order (market or limit) with its status
- **Position**: Represents an open or closed trading position with P&L tracking

## Ports

The ports directory contains interfaces that define how the domain interacts with the outside world:

- **Repositories**: Define data persistence operations for our domain models

## Design Principles

1. **Domain Independence**: The domain layer has no dependencies on external frameworks or libraries
2. **Rich Domain Models**: Models contain business logic and validation rules
3. **Immutable State**: Use value objects and immutable data where possible
4. **Clear Boundaries**: Well-defined interfaces for external interactions

## Usage

The domain layer is used by the application layer, which implements use cases by orchestrating domain objects. The infrastructure layer provides concrete implementations of the ports.

## Testing

Domain models and business logic should be thoroughly tested in isolation, without dependencies on external systems or frameworks. 
# Infrastructure Layer

This directory contains the infrastructure layer of our trading system, implementing concrete adapters for external services and dependencies.

## Structure

```
infrastructure/
├── persistence/    # Database adapters
│   ├── postgres/   # PostgreSQL implementations
│   └── redis/      # Redis implementations
├── exchange/       # Exchange API adapters
│   ├── binance/    # Binance API client
│   └── gemini/     # Gemini API client
├── marketdata/     # Market data adapters
│   ├── websocket/  # WebSocket implementations
│   └── rest/       # REST API implementations
└── risk/          # Risk management adapters
    └── alerts/     # Risk alert implementations
```

## Adapters

The infrastructure layer provides concrete implementations of the interfaces (ports) defined in the domain layer:

### Persistence Adapters
- **PostgresTradeRepository**: PostgreSQL implementation of TradeRepository
- **PostgresOrderRepository**: PostgreSQL implementation of OrderRepository
- **PostgresPositionRepository**: PostgreSQL implementation of PositionRepository
- **RedisMarketDataCache**: Redis implementation for market data caching

### Exchange Adapters
- **BinanceExchangeClient**: Binance API implementation
- **GeminiExchangeClient**: Gemini API implementation

### Market Data Adapters
- **WebSocketPriceStream**: Real-time price updates via WebSocket
- **RESTMarketDataClient**: Historical data retrieval via REST API

### Risk Management Adapters
- **AlertManager**: Risk alert implementation
- **MetricsCollector**: Risk metrics collection and monitoring

## Design Principles

1. **Interface Compliance**: All adapters must implement interfaces defined in the domain layer
2. **Dependency Injection**: Adapters are injected into application services
3. **Error Handling**: Adapters translate external errors into domain errors
4. **Configuration**: Adapters are configured via environment variables or configuration files

## Dependencies

The infrastructure layer:
- Implements interfaces from the domain layer
- May depend on external libraries and frameworks
- Should not be depended upon by domain or application layers

## Testing

Adapters should be tested with:
- Integration tests with actual external services
- Mocked external dependencies for unit tests
- Error handling and edge case scenarios
- Performance and reliability tests 
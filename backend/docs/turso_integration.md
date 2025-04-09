# TursoDB Integration

This document provides an overview of the TursoDB integration in the Go Crypto Bot project.

## Overview

TursoDB is a distributed SQLite database built on libSQL (a SQLite fork) that combines the simplicity and familiarity of SQLite with distributed database capabilities. It provides:

- SQLite compatibility with distributed capabilities
- Edge-deployed database instances for low-latency access
- Embedded replicas for offline-first applications
- Automatic synchronization between instances
- Serverless deployment model

## Architecture

The TursoDB integration follows a phased approach with an abstraction layer that allows for gradual migration from SQLite to TursoDB. The key components are:

1. **Repository Interface**: A common interface for database operations that abstracts away the underlying implementation.
2. **SQLite Repository**: The existing SQLite implementation that follows the repository interface.
3. **TursoDB Repository**: A new implementation that uses TursoDB with the same interface.
4. **Repository Factory**: A factory that creates the appropriate repository based on configuration.
5. **Synchronization Manager**: A component that handles synchronization with the cloud database.
6. **Migration Manager**: A component that handles database migrations.

## Getting Started

### Prerequisites

1. Install the TursoDB CLI:
   ```bash
   curl -sSfL https://get.tur.so/install.sh | bash
   ```

2. Create a TursoDB database:
   ```bash
   turso db create my-crypto-bot-db
   ```

3. Generate an authentication token:
   ```bash
   turso db tokens create my-crypto-bot-db
   ```

### Configuration

The TursoDB integration can be configured using environment variables or a configuration file:

```bash
# Enable TursoDB
export TURSO_ENABLED=true

# TursoDB URL (from the turso db show command)
export TURSO_URL=libsql://my-crypto-bot-db.turso.io

# TursoDB authentication token
export TURSO_AUTH_TOKEN=your-token-here

# Enable synchronization (for embedded replicas)
export TURSO_SYNC_ENABLED=true

# Synchronization interval in seconds
export TURSO_SYNC_INTERVAL_SECONDS=300
```

### Usage

#### Using the Repository Factory

```go
import (
    "github.com/ryanlisse/go-crypto-bot/internal/repository"
    "github.com/ryanlisse/go-crypto-bot/internal/repository/database"
)

// Load configuration
config := database.LoadConfig(nil)

// Create database repository
db, err := database.GetRepositoryFromEnv(config)
if err != nil {
    log.Fatalf("Failed to create database repository: %v", err)
}
defer db.Close()

// Create repository factory
factory := repository.NewFactory(db)

// Create balance history repository
balanceRepo := factory.NewBalanceHistoryRepository()

// Use the repository
history, err := balanceRepo.GetBalanceHistory(ctx, startTime, endTime)
```

#### Using the Synchronization Manager

```go
// Create sync manager
syncManager := database.NewSyncManager(db, config.SyncInterval)

// Start synchronization
syncManager.Start()
defer syncManager.Stop()

// Trigger manual synchronization
err := syncManager.SyncNow(ctx)
if err != nil {
    log.Printf("Failed to synchronize: %v", err)
}
```

## Testing

To test the TursoDB integration, you can use the `turso-test` command:

```bash
# Test with SQLite
go run cmd/turso-test/main.go

# Test with TursoDB (direct connection)
go run cmd/turso-test/main.go -turso -turso-url=libsql://my-crypto-bot-db.turso.io -turso-token=your-token-here

# Test with TursoDB (embedded replica with sync)
go run cmd/turso-test/main.go -turso -turso-url=libsql://my-crypto-bot-db.turso.io -turso-token=your-token-here -sync
```

## Migration Strategy

The migration to TursoDB follows a phased approach:

1. **Phase 1: Create Abstraction Layer**
   - Implement a database interface that works with both SQLite and TursoDB
   - Start with non-critical repositories (e.g., balance history, performance metrics)
   - Keep using SQLite as the primary database

2. **Phase 2: Shadow Mode**
   - Run both databases in parallel, writing to both but reading from SQLite
   - Implement validation to ensure data consistency between systems
   - Add monitoring for performance and reliability

3. **Phase 3: Gradual Cutover**
   - Once validation is successful, gradually switch reads to TursoDB
   - Start with less critical features and monitor closely
   - Maintain SQLite as a fallback option

4. **Phase 4: Full Migration**
   - Complete migration to TursoDB for all repositories
   - Optimize configuration for trading workloads
   - Implement comprehensive monitoring and alerting

## Troubleshooting

### Common Issues

1. **Connection Errors**
   - Verify that the TursoDB URL and authentication token are correct
   - Check if your IP is allowed to access the database
   - Ensure the database exists and is running

2. **Synchronization Issues**
   - Check if synchronization is enabled
   - Verify that the local database file is writable
   - Check for network connectivity issues

3. **Migration Errors**
   - Ensure that the schema is compatible with both SQLite and TursoDB
   - Check for any SQLite-specific features that are not supported by TursoDB

### Logging

The TursoDB integration includes comprehensive logging to help diagnose issues:

- Connection status and errors
- Synchronization events and timestamps
- Migration progress and errors

## References

- [TursoDB Documentation](https://docs.turso.tech)
- [Go SDK Reference](https://docs.turso.tech/sdk/go/quickstart)
- [SQLite Compatibility Guide](https://docs.turso.tech/reference/sqlite-compatibility)
- [Embedded Replicas Documentation](https://docs.turso.tech/features/embedded-replicas)
- [Turso CLI Reference](https://docs.turso.tech/reference/turso-cli)

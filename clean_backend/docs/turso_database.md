# Turso Database Integration

This document describes the Turso database integration in the Go Crypto Bot backend.

## Overview

The application uses [Turso](https://turso.tech/) as its primary database, which is a distributed SQLite database built on libSQL. Turso provides:

- SQLite compatibility with distributed capabilities
- Local replicas that sync with the cloud
- High performance and low latency
- Simplified operations

## Configuration

Turso database configuration is defined in the `config.yaml` file and can be overridden with environment variables:

```yaml
database:
  type: "turso"
  dsn: "libsql://local.db"
  turso_url: "${TURSO_DB_URL}"
  auth_token: "${TURSO_AUTH_TOKEN}"
  max_idle_conns: 5
  max_open_conns: 10
  conn_max_lifetime_minutes: 60
  enable_logging: false
  auto_migrate: true
```

### Environment Variables

- `DATABASE_TYPE`: Set to "turso" to use Turso database
- `DATABASE_DSN`: Local database path (used as fallback)
- `TURSO_DB_URL`: URL of your Turso database (e.g., `libsql://your-db-name.turso.io`)
- `TURSO_AUTH_TOKEN`: Authentication token for your Turso database
- `TURSO_SYNC_ENABLED`: Enable/disable periodic sync (default: true)
- `TURSO_SYNC_INTERVAL_SECONDS`: Interval for periodic sync in seconds (default: 300)

## Implementation Details

### Connection Modes

The application supports two connection modes:

1. **Remote Sync Mode**: Connects to the Turso cloud database and maintains a local replica that syncs periodically.
2. **Local-Only Mode**: Uses a local SQLite database without remote synchronization (fallback mode).

### Embedded Replicas

When using Remote Sync Mode, the application creates an embedded replica of the Turso database. This provides:

- Local database access for low latency
- Automatic synchronization with the cloud database
- Offline capabilities when the cloud database is unavailable

### Periodic Synchronization

The application sets up a background goroutine that periodically syncs the local replica with the cloud database. The sync interval is configurable via environment variables.

### Fallback Mechanism

If the connection to the Turso cloud database fails, the application automatically falls back to Local-Only Mode, ensuring the application can still function without the cloud database.

## Usage

### Creating a New Turso Database

1. Sign up for a Turso account at [turso.tech](https://turso.tech/)
2. Install the Turso CLI: `brew install tursodatabase/tap/turso`
3. Authenticate: `turso auth login`
4. Create a database: `turso db create my-crypto-bot`
5. Get the database URL: `turso db show my-crypto-bot --url`
6. Create an auth token: `turso db tokens create my-crypto-bot`

### Configuring the Application

1. Set the `TURSO_DB_URL` environment variable to the database URL
2. Set the `TURSO_AUTH_TOKEN` environment variable to the auth token
3. Set `DATABASE_TYPE` to "turso"

## Troubleshooting

### Connection Issues

If the application fails to connect to the Turso database, check:

1. The `TURSO_DB_URL` and `TURSO_AUTH_TOKEN` environment variables are set correctly
2. Your network connection allows access to the Turso database
3. The auth token has not expired

### Sync Issues

If synchronization is not working:

1. Check the application logs for sync-related errors
2. Verify that `TURSO_SYNC_ENABLED` is set to "true"
3. Try increasing the log level to "debug" to see more detailed sync information

## Limitations

- The Turso integration uses the libSQL Go client, which is still under development
- Some advanced SQLite features may not be fully supported
- The free tier of Turso has limitations on database size and operations

## References

- [Turso Documentation](https://docs.turso.tech/)
- [libSQL Go Client](https://github.com/tursodatabase/go-libsql)
- [GORM LibSQL Driver](https://github.com/ekristen/gorm-libsql)

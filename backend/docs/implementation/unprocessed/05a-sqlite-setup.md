# SQLite Setup for Go Crypto Bot

This guide covers the basic setup of SQLite for the Go crypto trading bot, focusing on simplicity and performance.

## 1. Installing the SQLite Driver

Add the SQLite driver to your project:

```bash
go get github.com/mattn/go-sqlite3
```

This is a CGO-enabled package, so make sure you have a C compiler installed on your system.

## 2. Database Connection Setup

Create a simple database package to manage the SQLite connection:

```go
// internal/platform/database/sqlite.go
package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Config holds SQLite configuration
type Config struct {
	Path      string        // Database file path
	InMemory  bool          // Whether to use in-memory database (for testing)
	BusyTimeout time.Duration // Timeout when database is locked
}

// DefaultConfig returns sensible defaults
func DefaultConfig() Config {
	return Config{
		Path:      "data/cryptobot.db",
		InMemory:  false,
		BusyTimeout: 5 * time.Second,
	}
}

// Connect creates a new SQLite database connection
func Connect(cfg Config) (*sql.DB, error) {
	if !cfg.InMemory {
		// Ensure directory exists
		dir := filepath.Dir(cfg.Path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("creating database directory: %w", err)
		}
	}

	// Build DSN (connection string)
	dsn := cfg.Path
	if cfg.InMemory {
		dsn = ":memory:"
	}

	// Add connection parameters for better performance and reliability
	dsn = fmt.Sprintf("%s?_journal_mode=WAL&_busy_timeout=%d&_foreign_keys=on&_cache_size=5000",
		dsn, cfg.BusyTimeout.Milliseconds())

	// Open the database
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(1)              // SQLite only supports one writer
	db.SetMaxIdleConns(1)              // Keep one connection in the pool
	db.SetConnMaxLifetime(0)           // Connections are reused forever

	return db, nil
}

// Close safely closes the database connection
func Close(db *sql.DB) error {
	if db == nil {
		return nil
	}
	return db.Close()
}
```

## 3. Performance Optimizations

The connection string includes several SQLite-specific optimizations:

- `_journal_mode=WAL`: Uses Write-Ahead Logging for better concurrency
- `_busy_timeout=5000`: Sets a timeout of 5 seconds when the database is locked
- `_foreign_keys=on`: Enables foreign key constraints
- `_cache_size=5000`: Increases the page cache to 5000 pages (about 20MB)

## 4. Connection Management in Main Application

Use the SQLite connection in your main application:

```go
// cmd/server/main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ryanlisse/cryptobot/internal/platform/database"
)

func main() {
	// Create a context that's canceled on Ctrl+C
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Printf("Received signal: %v", sig)
		cancel()
	}()

	// Load database config
	dbConfig := database.DefaultConfig()
	
	// Connect to database
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)
	log.Println("Connected to database successfully")

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed successfully")

	// Initialize repositories...
	// Initialize services...
	// Start the server...

	// Wait for cancellation
	<-ctx.Done()
	log.Println("Shutting down gracefully...")
}
```

## 5. Testing with In-Memory Database

For testing, use an in-memory SQLite database for speed and isolation:

```go
// internal/platform/database/database_test.go
package database_test

import (
	"testing"

	"github.com/ryanlisse/cryptobot/internal/platform/database"
)

func TestConnect(t *testing.T) {
	cfg := database.Config{
		InMemory: true,
	}

	db, err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to in-memory database: %v", err)
	}
	defer database.Close(db)

	// Test database connection
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}
```

## 6. Common SQLite Gotchas

### Concurrency Limitations

SQLite allows only one writer at a time. Keep this in mind when designing your application:

- Use the WAL journal mode (as configured above)
- Keep transactions as short as possible
- Be prepared to retry on SQLITE_BUSY errors

### Database Locking

If you see "database is locked" errors:

1. Ensure you're closing all transactions properly
2. Consider increasing the busy timeout
3. Make sure you're closing all database connections

### Integer Primary Keys

SQLite's `INTEGER PRIMARY KEY` is special - it's an alias for the ROWID and is optimized. Use it for your ID columns instead of `INT` or other types.

### Boolean Values

SQLite doesn't have a boolean type. Use integers (0/1) and convert in your Go code:

```go
// In your repository
isDeleted := 0
if coin.IsDeleted {
    isDeleted = 1
}

// When reading from the database
var isDeleted int
// ... scan into isDeleted
coin.IsDeleted = isDeleted == 1
```

## 7. Database Backup Strategy

Implement a simple backup function:

```go
// internal/platform/database/backup.go
package database

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Backup creates a copy of the database file
func Backup(db *sql.DB, backupDir string) (string, error) {
	// Get database path
	var path string
	err := db.QueryRow("PRAGMA database_list").Scan(nil, &path, nil)
	if err != nil {
		return "", fmt.Errorf("getting database path: %w", err)
	}

	// Skip backup for in-memory database
	if path == ":memory:" {
		return "", nil
	}

	// Create backup directory
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("creating backup directory: %w", err)
	}

	// Create backup filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	backupName := fmt.Sprintf("cryptobot-%s.db", timestamp)
	backupPath := filepath.Join(backupDir, backupName)

	// Use SQLite's backup API through a VACUUM operation (requires mattn/go-sqlite3)
	if _, err := db.Exec(fmt.Sprintf("VACUUM INTO '%s'", backupPath)); err != nil {
		// If VACUUM INTO is not supported, fall back to file copy
		src, err := os.Open(path)
		if err != nil {
			return "", fmt.Errorf("opening source database: %w", err)
		}
		defer src.Close()

		dst, err := os.Create(backupPath)
		if err != nil {
			return "", fmt.Errorf("creating backup file: %w", err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return "", fmt.Errorf("copying database: %w", err)
		}
	}

	return backupPath, nil
}
```

## Next Steps

With the SQLite connection setup, you can now implement:

1. Database migrations (see `05b-sqlite-migrations.md`)
2. Repository implementations (see `05c-sqlite-repositories.md`)

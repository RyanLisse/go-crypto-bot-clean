# Database Layer Implementation

This document provides an overview of the database layer for the Go crypto trading bot, focusing on the repository pattern and interface definitions. For specific implementation details, refer to:

- [04a-sqlite-implementation.md](04a-sqlite-implementation.md) - SQLite setup and connections
- [04b-sqlite-migrations.md](04b-sqlite-migrations.md) - Migration system and database schema
- [04c-repository-implementations.md](04c-repository-implementations.md) - Concrete repository implementations

## 1. Repository Pattern Overview

The database layer follows the repository pattern, providing a clean abstraction over data persistence mechanisms:

```
┌────────────────────┐      ┌─────────────────┐      ┌─────────────────────┐
│                    │      │                 │      │                     │
│   Domain Services  │─────▶│   Repositories  │─────▶│   Database (SQLite) │
│                    │      │   (Interfaces)  │      │                     │
└────────────────────┘      └─────────────────┘      └─────────────────────┘
```

This approach offers several key benefits:

1. **Separation of Concerns**: Core business logic depends only on abstracted repository interfaces, not on database implementation details
2. **Testability**: Domain services can be unit tested with mocked repositories
3. **Flexibility**: The underlying database technology can be changed without affecting business logic
4. **Clean Architecture**: Follows hexagonal architecture principles with domain at the center

## 2. Repository Factory

A repository factory provides access to all repositories and centralizes their instantiation:

```go
// internal/domain/repository/factory.go
package repository

import (
	"github.com/jmoiron/sqlx"
)

// Factory creates and manages repository instances
type Factory struct {
	db *sqlx.DB
	
	boughtCoin      BoughtCoinRepository
	newCoin         NewCoinRepository
	purchaseDecision PurchaseDecisionRepository
	logEvent        LogEventRepository
	position        PositionRepository
}

// NewFactory creates a new repository factory
func NewFactory(db *sqlx.DB) *Factory {
	return &Factory{
		db: db,
	}
}

// BoughtCoin returns the BoughtCoinRepository instance
func (f *Factory) BoughtCoin() BoughtCoinRepository {
	if f.boughtCoin == nil {
		f.boughtCoin = NewSQLiteBoughtCoinRepository(f.db)
	}
	return f.boughtCoin
}

// Additional repository getters follow the same pattern...
```
    "fmt"
    "log"
    "os"
    "path/filepath"

    "github.com/jmoiron/sqlx"
    _ "github.com/mattn/go-sqlite3"
)

// DB is the global database instance
var DB *sqlx.DB

// Initialize sets up the database connection
func Initialize(dbPath string) error {
    // Ensure the directory exists
    dbDir := filepath.Dir(dbPath)
    if err := os.MkdirAll(dbDir, 0755); err != nil {
        return fmt.Errorf("failed to create database directory: %w", err)
    }

    var err error
    DB, err = sqlx.Connect("sqlite3", dbPath)
    if err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    // Set connection parameters for SQLite
    DB.SetMaxOpenConns(1) // SQLite only supports one writer at a time
    DB.SetMaxIdleConns(1)
    
    log.Printf("Connected to database at %s", dbPath)
    
    // Run migrations
    if err := RunMigrations(DB); err != nil {
        return fmt.Errorf("failed to run migrations: %w", err)
    }
    
    return nil
}

// Close closes the database connection
func Close() error {
    if DB != nil {
        return DB.Close()
    }
    return nil
}
```

## 3. Repository Interfaces

The repository interfaces define contracts for accessing domain entities in the database:

```go
// internal/domain/repository/repository.go
package repository

import (
    "context"
    "time"
    
    "github.com/ryanlisse/cryptobot-backend/internal/domain/models"
)

// BoughtCoinRepository defines the interface for bought coin operations
type BoughtCoinRepository interface {
    Store(ctx context.Context, coin *models.BoughtCoin) (int64, error)
    FindByID(ctx context.Context, id int64) (*models.BoughtCoin, error)
    FindBySymbol(ctx context.Context, symbol string) (*models.BoughtCoin, error)
    FindAll(ctx context.Context, includeDeleted bool) ([]models.BoughtCoin, error)
    Update(ctx context.Context, coin *models.BoughtCoin) error
    SoftDelete(ctx context.Context, id int64) error
    Restore(ctx context.Context, id int64) error
    DeleteAllOlderThan(ctx context.Context, threshold time.Time) (int64, error)
}

// NewCoinRepository defines the interface for new coin operations
type NewCoinRepository interface {
    Store(ctx context.Context, coin *models.NewCoin) (int64, error)
    FindByID(ctx context.Context, id int64) (*models.NewCoin, error)
    FindBySymbol(ctx context.Context, symbol string) (*models.NewCoin, error)
    FindAll(ctx context.Context) ([]models.NewCoin, error)
    FindAllWithStatus(ctx context.Context, status string) ([]models.NewCoin, error)
    UpdateStatus(ctx context.Context, id int64, status string) error
    Delete(ctx context.Context, id int64) error
}

// PurchaseDecisionRepository defines the interface for purchase decision operations
type PurchaseDecisionRepository interface {
    Store(ctx context.Context, decision *models.PurchaseDecision) (int64, error)
    FindByID(ctx context.Context, id int64) (*models.PurchaseDecision, error)
    FindBySymbol(ctx context.Context, symbol string) ([]models.PurchaseDecision, error)
    FindAll(ctx context.Context, limit, offset int) ([]models.PurchaseDecision, error)
    GetStats(ctx context.Context) (map[string]interface{}, error)
}

// LogEventRepository defines the interface for log event operations
type LogEventRepository interface {
    Store(ctx context.Context, event *models.LogEvent) (int64, error)
    FindByID(ctx context.Context, id int64) (*models.LogEvent, error)
    FindAll(ctx context.Context, limit, offset int, level, component string) ([]models.LogEvent, error)
    DeleteOlderThan(ctx context.Context, threshold time.Time) (int64, error)
    Count(ctx context.Context, level, component string) (int64, error)
}

// PositionRepository defines the interface for trading position operations
type PositionRepository interface {
    Store(ctx context.Context, position *models.Position) (int64, error)
    FindByID(ctx context.Context, id int64) (*models.Position, error)
    FindBySymbol(ctx context.Context, symbol string) (*models.Position, error)
    FindAll(ctx context.Context, status string) ([]models.Position, error)
    Update(ctx context.Context, position *models.Position) error
    Delete(ctx context.Context, id int64) error
}

// Factory creates repository instances
type Factory struct {
    db *sqlx.DB
}

// NewFactory creates a new repository factory
func NewFactory(db *sqlx.DB) *Factory {
    return &Factory{db: db}
}

// BoughtCoin creates a new BoughtCoinRepository
func (f *Factory) BoughtCoin() BoughtCoinRepository {
    return NewSQLiteBoughtCoinRepository(f.db)
}

// NewCoin creates a new NewCoinRepository
func (f *Factory) NewCoin() NewCoinRepository {
    return NewSQLiteNewCoinRepository(f.db)
}

// PurchaseDecision creates a new PurchaseDecisionRepository
func (f *Factory) PurchaseDecision() PurchaseDecisionRepository {
    return NewSQLitePurchaseDecisionRepository(f.db)
}

// LogEvent creates a new LogEventRepository
func (f *Factory) LogEvent() LogEventRepository {
    return NewSQLiteLogEventRepository(f.db)
}

// Position creates a new PositionRepository
func (f *Factory) Position() PositionRepository {
    return NewSQLitePositionRepository(f.db)
}
```

## 4. Migrations

We'll implement database migrations to ensure a consistent schema:

```go
// internal/platform/database/migrations.go
package database

import (
    "github.com/jmoiron/sqlx"
)

// RunMigrations executes all database migrations
func RunMigrations(db *sqlx.DB) error {
    migrations := []string{
        createBoughtCoinsTable,
        createNewCoinsTable,
        createPurchaseDecisionsTable,
        createLogEventsTable,
        createPositionsTable,
    }

    for _, migration := range migrations {
        if _, err := db.Exec(migration); err != nil {
            return err
        }
    }

    return nil
}

// SQL migration statements
const (
    createBoughtCoinsTable = `
    CREATE TABLE IF NOT EXISTS bought_coins (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        symbol TEXT NOT NULL UNIQUE,
        purchase_price REAL NOT NULL,
        quantity REAL NOT NULL,
        purchase_time TIMESTAMP NOT NULL,
        is_deleted BOOLEAN NOT NULL DEFAULT 0,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_bought_coins_symbol ON bought_coins(symbol);
    CREATE INDEX IF NOT EXISTS idx_bought_coins_is_deleted ON bought_coins(is_deleted);
    `

    createNewCoinsTable = `
    CREATE TABLE IF NOT EXISTS new_coins (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        symbol TEXT NOT NULL UNIQUE,
        detected_at TIMESTAMP NOT NULL,
        status TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_new_coins_symbol ON new_coins(symbol);
    CREATE INDEX IF NOT EXISTS idx_new_coins_status ON new_coins(status);
    `

    createPurchaseDecisionsTable = `
    CREATE TABLE IF NOT EXISTS purchase_decisions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        symbol TEXT NOT NULL,
        decision TEXT NOT NULL,
        reason TEXT NOT NULL,
        decision_time TIMESTAMP NOT NULL,
        created_at TIMESTAMP NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_purchase_decisions_symbol ON purchase_decisions(symbol);
    CREATE INDEX IF NOT EXISTS idx_purchase_decisions_decision ON purchase_decisions(decision);
    `

    createLogEventsTable = `
    CREATE TABLE IF NOT EXISTS log_events (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        level TEXT NOT NULL,
        message TEXT NOT NULL,
        component TEXT NOT NULL,
        timestamp TIMESTAMP NOT NULL,
        created_at TIMESTAMP NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_log_events_level ON log_events(level);
    CREATE INDEX IF NOT EXISTS idx_log_events_component ON log_events(component);
    CREATE INDEX IF NOT EXISTS idx_log_events_timestamp ON log_events(timestamp);
    `

    createPositionsTable = `
    CREATE TABLE IF NOT EXISTS positions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        symbol TEXT NOT NULL,
        entry_price REAL NOT NULL,
        current_price REAL NOT NULL,
        quantity REAL NOT NULL,
        status TEXT NOT NULL,
        stop_loss REAL,
        take_profit REAL,
        entry_time TIMESTAMP NOT NULL,
        exit_time TIMESTAMP,
        exit_price REAL,
        profit_loss REAL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_positions_symbol ON positions(symbol);
    CREATE INDEX IF NOT EXISTS idx_positions_status ON positions(status);
    `
)
```

## 5. SQLite Implementations

Now let's implement a sample repository to demonstrate the concrete implementation:

```go
// internal/platform/database/repository/sqlite_bought_coin_repository.go
package repository

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    "github.com/jmoiron/sqlx"
    
    "github.com/ryanlisse/cryptobot-backend/internal/domain/models"
)

// SQLiteBoughtCoinRepository is a SQLite implementation of BoughtCoinRepository
type SQLiteBoughtCoinRepository struct {
    db *sqlx.DB
}

// NewSQLiteBoughtCoinRepository creates a new SQLite bought coin repository
func NewSQLiteBoughtCoinRepository(db *sqlx.DB) *SQLiteBoughtCoinRepository {
    return &SQLiteBoughtCoinRepository{db: db}
}

// Store inserts a new bought coin into the database
func (r *SQLiteBoughtCoinRepository) Store(ctx context.Context, coin *models.BoughtCoin) (int64, error) {
    now := time.Now()
    coin.CreatedAt = now
    coin.UpdatedAt = now

    query := `
        INSERT INTO bought_coins (
            symbol, purchase_price, quantity, purchase_time, is_deleted, created_at, updated_at
        ) VALUES (
            ?, ?, ?, ?, ?, ?, ?
        )
    `

    result, err := r.db.ExecContext(
        ctx,
        query,
        coin.Symbol,
        coin.PurchasePrice,
        coin.Quantity,
        coin.PurchaseTime,
        coin.IsDeleted,
        coin.CreatedAt,
        coin.UpdatedAt,
    )
    if err != nil {
        return 0, fmt.Errorf("failed to insert bought coin: %w", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("failed to get last insert ID: %w", err)
    }

    coin.ID = id
    return id, nil
}

// FindByID finds a bought coin by ID
func (r *SQLiteBoughtCoinRepository) FindByID(ctx context.Context, id int64) (*models.BoughtCoin, error) {
    var coin models.BoughtCoin
    
    query := `
        SELECT id, symbol, purchase_price, quantity, purchase_time, is_deleted, created_at, updated_at
        FROM bought_coins
        WHERE id = ?
    `
    
    err := r.db.GetContext(ctx, &coin, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get bought coin by ID: %w", err)
    }
    
    return &coin, nil
}
```

## 6. Repository Factory

The Repository Factory pattern provides a clean way to instantiate repositories:

```go
// internal/platform/database/repository/factory.go
package repository

import (
    "github.com/jmoiron/sqlx"
)

// RepositoryFactory creates and initializes repositories
type RepositoryFactory struct {
    db *sqlx.DB
    
    // Cached repositories
    boughtCoinRepo     BoughtCoinRepository
    newCoinRepo        NewCoinRepository
    purchaseDecisionRepo PurchaseDecisionRepository
    logEventRepo       LogEventRepository
    positionRepo       PositionRepository
}

// NewRepositoryFactory creates a new repository factory
func NewRepositoryFactory(db *sqlx.DB) *RepositoryFactory {
    return &RepositoryFactory{db: db}
}

// BoughtCoin returns a BoughtCoinRepository
func (f *RepositoryFactory) BoughtCoin() BoughtCoinRepository {
    if f.boughtCoinRepo == nil {
        f.boughtCoinRepo = NewSQLiteBoughtCoinRepository(f.db)
    }
    return f.boughtCoinRepo
}

// NewCoin returns a NewCoinRepository
func (f *RepositoryFactory) NewCoin() NewCoinRepository {
    if f.newCoinRepo == nil {
        f.newCoinRepo = NewSQLiteNewCoinRepository(f.db)
    }
    return f.newCoinRepo
}

// PurchaseDecision returns a PurchaseDecisionRepository
func (f *RepositoryFactory) PurchaseDecision() PurchaseDecisionRepository {
    if f.purchaseDecisionRepo == nil {
        f.purchaseDecisionRepo = NewSQLitePurchaseDecisionRepository(f.db)
    }
    return f.purchaseDecisionRepo
}

// LogEvent returns a LogEventRepository
func (f *RepositoryFactory) LogEvent() LogEventRepository {
    if f.logEventRepo == nil {
        f.logEventRepo = NewSQLiteLogEventRepository(f.db)
    }
    return f.logEventRepo
}

// Position returns a PositionRepository
func (f *RepositoryFactory) Position() PositionRepository {
    if f.positionRepo == nil {
        f.positionRepo = NewSQLitePositionRepository(f.db)
    }
    return f.positionRepo
}
```

## 7. Transaction Management

For operations requiring transactions:

```go
// internal/platform/database/transaction.go
package database

import (
    "context"
    "fmt"

    "github.com/jmoiron/sqlx"
)

// RunInTransaction executes a function within a transaction
func RunInTransaction(ctx context.Context, db *sqlx.DB, fn func(*sqlx.Tx) error) error {
    tx, err := db.BeginTxx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }

    defer func() {
        if p := recover(); p != nil {
            // Rollback on panic
            _ = tx.Rollback()
            panic(p) // Re-throw panic after rollback
        } else if err != nil {
            // Rollback on error
            _ = tx.Rollback()
        }
    }()

    if err = fn(tx); err != nil {
        return err
    }

    if err = tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}
```

## 8. Usage Example

Here's an example of how to use these repositories in a service:

```go
// internal/domain/service/trade_service.go
package service

import (
    "context"
    "fmt"
    "time"

    "github.com/ryanlisse/cryptobot-backend/internal/domain/models"
    "github.com/ryanlisse/cryptobot-backend/internal/domain/repository"
)

// TradeService handles trading operations
type TradeService struct {
    boughtCoinRepo repository.BoughtCoinRepository
    positionRepo   repository.PositionRepository
}

// NewTradeService creates a new trade service
func NewTradeService(
    boughtCoinRepo repository.BoughtCoinRepository,
    positionRepo repository.PositionRepository,
) *TradeService {
    return &TradeService{
        boughtCoinRepo: boughtCoinRepo,
        positionRepo:   positionRepo,
    }
}

// PurchaseCoin purchases a coin and records it
func (s *TradeService) PurchaseCoin(
    ctx context.Context,
    symbol string,
    price, quantity float64,
) (*models.BoughtCoin, error) {
    // Create a new bought coin record
    coin := &models.BoughtCoin{
        Symbol:        symbol,
        PurchasePrice: price,
        Quantity:      quantity,
        PurchaseTime:  time.Now(),
        IsDeleted:     false,
    }

    // Store in the database
    id, err := s.boughtCoinRepo.Store(ctx, coin)
    if err != nil {
        return nil, fmt.Errorf("failed to store bought coin: %w", err)
    }

    coin.ID = id
    return coin, nil
}
```

## 9. Best Practices for Database Operations

1. **Use Transactions** for operations that modify multiple records
2. **Handle SQL No Rows** by returning nil for "not found" cases
3. **Use Prepared Statements** to prevent SQL injection
4. **Include Context** in all DB operations for cancellation and timeouts
5. **Use Indexes** on columns frequently used in WHERE clauses
6. **Keep Repositories Clean** by avoiding business logic in repositories
7. **Test Repository Implementations** thoroughly
8. **Consider Query Builder** like Squirrel for complex dynamic queries

Following these guidelines will ensure a robust, maintainable database layer for the crypto trading bot.

# Repository Implementations

This document covers the concrete SQLite implementations of the repository interfaces defined in the [Database Layer Overview](04-database-layer.md).

## 1. Repository Interface Implementation Pattern

Each repository implementation follows a consistent pattern:

1. Define a struct that embeds the database connection
2. Implement all interface methods using SQL queries
3. Use prepared statements for queries with parameters
4. Handle transactions appropriately for write operations
5. Return domain models, not database-specific structures

## 2. BoughtCoin Repository Implementation

```go
// internal/platform/sqlite/bought_coin_repository.go
package sqlite

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "time"

    "github.com/jmoiron/sqlx"
    
    "github.com/ryanlisse/cryptobot-backend/internal/domain/models"
    "github.com/ryanlisse/cryptobot-backend/internal/domain/repository"
)

// BoughtCoinRepository implements repository.BoughtCoinRepository using SQLite
type BoughtCoinRepository struct {
    db *sqlx.DB
}

// NewBoughtCoinRepository creates a new SQLite BoughtCoinRepository
func NewBoughtCoinRepository(db *sqlx.DB) repository.BoughtCoinRepository {
    return &BoughtCoinRepository{
        db: db,
    }
}

// FindAll returns all bought coins
func (r *BoughtCoinRepository) FindAll(ctx context.Context, includeDeleted bool) ([]models.BoughtCoin, error) {
    query := `
        SELECT 
            id, symbol, purchase_price, quantity, purchased_at, 
            deleted, sold_at, stop_loss, take_profit, strategy
        FROM bought_coins
    `
    
    // Add filter for deleted if not including deleted coins
    if !includeDeleted {
        query += " WHERE deleted = 0"
    }
    
    // Add ordering
    query += " ORDER BY purchased_at DESC"
    
    var coins []models.BoughtCoin
    if err := r.db.SelectContext(ctx, &coins, query); err != nil {
        return nil, fmt.Errorf("failed to find bought coins: %w", err)
    }
    
    return coins, nil
}

// FindByID returns a bought coin by ID
func (r *BoughtCoinRepository) FindByID(ctx context.Context, id int64) (*models.BoughtCoin, error) {
    query := `
        SELECT 
            id, symbol, purchase_price, quantity, purchased_at, 
            deleted, sold_at, stop_loss, take_profit, strategy
        FROM bought_coins 
        WHERE id = ?
    `
    
    var coin models.BoughtCoin
    if err := r.db.GetContext(ctx, &coin, query, id); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil // Not found, return nil without error
        }
        return nil, fmt.Errorf("failed to find bought coin by ID: %w", err)
    }
    
    return &coin, nil
}

// FindBySymbol returns a bought coin by symbol
func (r *BoughtCoinRepository) FindBySymbol(ctx context.Context, symbol string) (*models.BoughtCoin, error) {
    query := `
        SELECT 
            id, symbol, purchase_price, quantity, purchased_at, 
            deleted, sold_at, stop_loss, take_profit, strategy
        FROM bought_coins 
        WHERE symbol = ? AND deleted = 0
    `
    
    var coin models.BoughtCoin
    if err := r.db.GetContext(ctx, &coin, query, symbol); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil // Not found, return nil without error
        }
        return nil, fmt.Errorf("failed to find bought coin by symbol: %w", err)
    }
    
    return &coin, nil
}

// Store stores a bought coin
func (r *BoughtCoinRepository) Store(ctx context.Context, coin *models.BoughtCoin) (int64, error) {
    if coin.ID > 0 {
        return r.update(ctx, coin)
    }
    return r.insert(ctx, coin)
}

// insert inserts a new bought coin
func (r *BoughtCoinRepository) insert(ctx context.Context, coin *models.BoughtCoin) (int64, error) {
    query := `
        INSERT INTO bought_coins (
            symbol, purchase_price, quantity, purchased_at, 
            deleted, sold_at, stop_loss, take_profit, strategy
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
    
    // Set purchased_at if not set
    if coin.PurchasedAt.IsZero() {
        coin.PurchasedAt = time.Now()
    }
    
    result, err := r.db.ExecContext(
        ctx, 
        query,
        coin.Symbol,
        coin.PurchasePrice,
        coin.Quantity,
        coin.PurchasedAt,
        coin.IsDeleted,
        coin.SoldAt,
        coin.StopLoss,
        coin.TakeProfit,
        coin.Strategy,
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

// update updates an existing bought coin
func (r *BoughtCoinRepository) update(ctx context.Context, coin *models.BoughtCoin) (int64, error) {
    query := `
        UPDATE bought_coins SET
            symbol = ?,
            purchase_price = ?,
            quantity = ?,
            purchased_at = ?,
            deleted = ?,
            sold_at = ?,
            stop_loss = ?,
            take_profit = ?,
            strategy = ?
        WHERE id = ?
    `
    
    _, err := r.db.ExecContext(
        ctx, 
        query,
        coin.Symbol,
        coin.PurchasePrice,
        coin.Quantity,
        coin.PurchasedAt,
        coin.IsDeleted,
        coin.SoldAt,
        coin.StopLoss,
        coin.TakeProfit,
        coin.Strategy,
        coin.ID,
    )
    if err != nil {
        return 0, fmt.Errorf("failed to update bought coin: %w", err)
    }
    
    return coin.ID, nil
}

// Delete marks a bought coin as deleted
func (r *BoughtCoinRepository) Delete(ctx context.Context, id int64) error {
    query := `
        UPDATE bought_coins SET
            deleted = 1,
            sold_at = ?
        WHERE id = ?
    `
    
    _, err := r.db.ExecContext(ctx, query, time.Now(), id)
    if err != nil {
        return fmt.Errorf("failed to delete bought coin: %w", err)
    }
    
    return nil
}
```

## 3. Position Repository Implementation

The Position repository shows how to work with more complex domain models:

```go
// internal/platform/sqlite/position_repository.go
package sqlite

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "time"

    "github.com/jmoiron/sqlx"
    
    "github.com/ryanlisse/cryptobot-backend/internal/domain/models"
    "github.com/ryanlisse/cryptobot-backend/internal/domain/repository"
)

// PositionRepository implements repository.PositionRepository using SQLite
type PositionRepository struct {
    db *sqlx.DB
}

// NewPositionRepository creates a new SQLite PositionRepository
func NewPositionRepository(db *sqlx.DB) repository.PositionRepository {
    return &PositionRepository{
        db: db,
    }
}

// Store stores a position
func (r *PositionRepository) Store(ctx context.Context, position *models.Position) (int64, error) {
    if position.ID > 0 {
        return r.update(ctx, position)
    }
    return r.insert(ctx, position)
}

// insert inserts a new position
func (r *PositionRepository) insert(ctx context.Context, position *models.Position) (int64, error) {
    query := `
        INSERT INTO positions (
            symbol, entry_price, quantity, status, stop_loss, 
            take_profit, created_at, closed_at, profit_loss, profit_loss_percent, notes
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
    
    // Set created_at if not set
    if position.CreatedAt.IsZero() {
        position.CreatedAt = time.Now()
    }
    
    result, err := r.db.ExecContext(
        ctx, 
        query,
        position.Symbol,
        position.EntryPrice,
        position.Quantity,
        position.Status,
        position.StopLoss,
        position.TakeProfit,
        position.CreatedAt,
        position.ClosedAt,
        position.ProfitLoss,
        position.ProfitLossPercent,
        position.Notes,
    )
    if err != nil {
        return 0, fmt.Errorf("failed to insert position: %w", err)
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("failed to get last insert ID: %w", err)
    }
    
    position.ID = id
    return id, nil
}

// update updates an existing position
func (r *PositionRepository) update(ctx context.Context, position *models.Position) (int64, error) {
    query := `
        UPDATE positions SET
            symbol = ?,
            entry_price = ?,
            quantity = ?,
            status = ?,
            stop_loss = ?,
            take_profit = ?,
            created_at = ?,
            closed_at = ?,
            profit_loss = ?,
            profit_loss_percent = ?,
            notes = ?
        WHERE id = ?
    `
    
    _, err := r.db.ExecContext(
        ctx, 
        query,
        position.Symbol,
        position.EntryPrice,
        position.Quantity,
        position.Status,
        position.StopLoss,
        position.TakeProfit,
        position.CreatedAt,
        position.ClosedAt,
        position.ProfitLoss,
        position.ProfitLossPercent,
        position.Notes,
        position.ID,
    )
    if err != nil {
        return 0, fmt.Errorf("failed to update position: %w", err)
    }
    
    return position.ID, nil
}

// FindByID returns a position by ID
func (r *PositionRepository) FindByID(ctx context.Context, id int64) (*models.Position, error) {
    query := `
        SELECT 
            id, symbol, entry_price, quantity, status, stop_loss, 
            take_profit, created_at, closed_at, profit_loss, profit_loss_percent, notes
        FROM positions 
        WHERE id = ?
    `
    
    var position models.Position
    if err := r.db.GetContext(ctx, &position, query, id); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil // Not found, return nil without error
        }
        return nil, fmt.Errorf("failed to find position by ID: %w", err)
    }
    
    return &position, nil
}

// FindBySymbol returns positions by symbol
func (r *PositionRepository) FindBySymbol(ctx context.Context, symbol string) ([]models.Position, error) {
    query := `
        SELECT 
            id, symbol, entry_price, quantity, status, stop_loss, 
            take_profit, created_at, closed_at, profit_loss, profit_loss_percent, notes
        FROM positions 
        WHERE symbol = ?
        ORDER BY created_at DESC
    `
    
    var positions []models.Position
    if err := r.db.SelectContext(ctx, &positions, query, symbol); err != nil {
        return nil, fmt.Errorf("failed to find positions by symbol: %w", err)
    }
    
    return positions, nil
}

// FindAll returns all positions with optional status filter
func (r *PositionRepository) FindAll(ctx context.Context, status string) ([]models.Position, error) {
    query := `
        SELECT 
            id, symbol, entry_price, quantity, status, stop_loss, 
            take_profit, created_at, closed_at, profit_loss, profit_loss_percent, notes
        FROM positions
    `
    
    // Add filter for status if provided
    var args []interface{}
    if status != "" {
        query += " WHERE status = ?"
        args = append(args, status)
    }
    
    // Add ordering
    query += " ORDER BY created_at DESC"
    
    var positions []models.Position
    if err := r.db.SelectContext(ctx, &positions, query, args...); err != nil {
        return nil, fmt.Errorf("failed to find positions: %w", err)
    }
    
    return positions, nil
}
```

## 4. New Coin Repository Implementation

```go
// internal/platform/sqlite/new_coin_repository.go
package sqlite

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "time"

    "github.com/jmoiron/sqlx"
    
    "github.com/ryanlisse/cryptobot-backend/internal/domain/models"
    "github.com/ryanlisse/cryptobot-backend/internal/domain/repository"
)

// NewCoinRepository implements repository.NewCoinRepository using SQLite
type NewCoinRepository struct {
    db *sqlx.DB
}

// NewNewCoinRepository creates a new SQLite NewCoinRepository
func NewNewCoinRepository(db *sqlx.DB) repository.NewCoinRepository {
    return &NewCoinRepository{
        db: db,
    }
}

// Store stores a new coin
func (r *NewCoinRepository) Store(ctx context.Context, coin *models.NewCoin) (int64, error) {
    query := `
        INSERT INTO new_coins (symbol, discovered_at)
        VALUES (?, ?)
        ON CONFLICT (symbol) DO UPDATE SET
        discovered_at = ?
    `
    
    // Set discovered_at if not set
    if coin.DiscoveredAt.IsZero() {
        coin.DiscoveredAt = time.Now()
    }
    
    result, err := r.db.ExecContext(
        ctx, 
        query,
        coin.Symbol,
        coin.DiscoveredAt,
        coin.DiscoveredAt,
    )
    if err != nil {
        return 0, fmt.Errorf("failed to store new coin: %w", err)
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("failed to get last insert ID: %w", err)
    }
    
    coin.ID = id
    return id, nil
}

// FindAll returns all new coins
func (r *NewCoinRepository) FindAll(ctx context.Context) ([]models.NewCoin, error) {
    query := `
        SELECT id, symbol, discovered_at
        FROM new_coins
        ORDER BY discovered_at DESC
    `
    
    var coins []models.NewCoin
    if err := r.db.SelectContext(ctx, &coins, query); err != nil {
        return nil, fmt.Errorf("failed to find new coins: %w", err)
    }
    
    return coins, nil
}

// FindBySymbol returns a new coin by symbol
func (r *NewCoinRepository) FindBySymbol(ctx context.Context, symbol string) (*models.NewCoin, error) {
    query := `
        SELECT id, symbol, discovered_at
        FROM new_coins
        WHERE symbol = ?
    `
    
    var coin models.NewCoin
    if err := r.db.GetContext(ctx, &coin, query, symbol); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil // Not found, return nil without error
        }
        return nil, fmt.Errorf("failed to find new coin by symbol: %w", err)
    }
    
    return &coin, nil
}

// FindRecent returns recently discovered new coins
func (r *NewCoinRepository) FindRecent(ctx context.Context, minutes int) ([]models.NewCoin, error) {
    query := `
        SELECT id, symbol, discovered_at
        FROM new_coins
        WHERE discovered_at > datetime('now', ?)
        ORDER BY discovered_at DESC
    `
    
    timeAgo := fmt.Sprintf("-%d minutes", minutes)
    
    var coins []models.NewCoin
    if err := r.db.SelectContext(ctx, &coins, query, timeAgo); err != nil {
        return nil, fmt.Errorf("failed to find recent new coins: %w", err)
    }
    
    return coins, nil
}
```

## 5. Log Event Repository Implementation

The log event repository demonstrates handling JSON storage in SQLite:

```go
// internal/platform/sqlite/log_event_repository.go
package sqlite

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/jmoiron/sqlx"
    
    "github.com/ryanlisse/cryptobot-backend/internal/domain/models"
    "github.com/ryanlisse/cryptobot-backend/internal/domain/repository"
)

// LogEventRepository implements repository.LogEventRepository using SQLite
type LogEventRepository struct {
    db *sqlx.DB
}

// NewLogEventRepository creates a new SQLite LogEventRepository
func NewLogEventRepository(db *sqlx.DB) repository.LogEventRepository {
    return &LogEventRepository{
        db: db,
    }
}

// Store stores a log event
func (r *LogEventRepository) Store(ctx context.Context, event *models.LogEvent) (int64, error) {
    query := `
        INSERT INTO log_events (level, message, context)
        VALUES (?, ?, ?)
    `
    
    // Marshal context to JSON if present
    var contextJSON []byte
    var err error
    if event.Context != nil {
        contextJSON, err = json.Marshal(event.Context)
        if err != nil {
            return 0, fmt.Errorf("failed to marshal context: %w", err)
        }
    }
    
    result, err := r.db.ExecContext(
        ctx, 
        query,
        event.Level,
        event.Message,
        contextJSON,
    )
    if err != nil {
        return 0, fmt.Errorf("failed to store log event: %w", err)
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("failed to get last insert ID: %w", err)
    }
    
    event.ID = id
    return id, nil
}

// FindAll returns all log events with optional filtering
func (r *LogEventRepository) FindAll(ctx context.Context, level string, since time.Time, limit int) ([]models.LogEvent, error) {
    query := `
        SELECT id, timestamp, level, message, context
        FROM log_events
        WHERE 1=1
    `
    
    var args []interface{}
    
    // Add filter for level if provided
    if level != "" {
        query += " AND level = ?"
        args = append(args, level)
    }
    
    // Add filter for timestamp if provided
    if !since.IsZero() {
        query += " AND timestamp > ?"
        args = append(args, since)
    }
    
    // Add ordering and limit
    query += " ORDER BY timestamp DESC"
    
    if limit > 0 {
        query += " LIMIT ?"
        args = append(args, limit)
    }
    
    rows, err := r.db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to find log events: %w", err)
    }
    defer rows.Close()
    
    var events []models.LogEvent
    for rows.Next() {
        var event models.LogEvent
        var contextJSON []byte
        
        if err := rows.Scan(
            &event.ID,
            &event.Timestamp,
            &event.Level,
            &event.Message,
            &contextJSON,
        ); err != nil {
            return nil, fmt.Errorf("failed to scan log event: %w", err)
        }
        
        // Unmarshal context if present
        if len(contextJSON) > 0 {
            if err := json.Unmarshal(contextJSON, &event.Context); err != nil {
                return nil, fmt.Errorf("failed to unmarshal context: %w", err)
            }
        }
        
        events = append(events, event)
    }
    
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating log event rows: %w", err)
    }
    
    return events, nil
}
```

## 6. Repository Testing

Each repository implementation should be tested thoroughly:

```go
// internal/platform/sqlite/bought_coin_repository_test.go
package sqlite_test

import (
    "context"
    "testing"
    "time"

    "github.com/jmoiron/sqlx"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    "github.com/ryanlisse/cryptobot-backend/internal/domain/models"
    "github.com/ryanlisse/cryptobot-backend/internal/platform/database"
    "github.com/ryanlisse/cryptobot-backend/internal/platform/sqlite"
)

func setupTestDB(t *testing.T) *sqlx.DB {
    db, err := sqlx.Open("sqlite3", ":memory:")
    require.NoError(t, err)
    
    // Run migrations
    err = database.RunMigrations(db)
    require.NoError(t, err)
    
    return db
}

func TestBoughtCoinRepository_Store(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    repo := sqlite.NewBoughtCoinRepository(db)
    ctx := context.Background()
    
    // Test insert
    coin := models.BoughtCoin{
        Symbol:        "BTCUSDT",
        PurchasePrice: 50000.0,
        Quantity:      0.1,
        PurchasedAt:   time.Now(),
        StopLoss:      49000.0,
        TakeProfit:    52000.0,
        Strategy:      "TEST",
    }
    
    id, err := repo.Store(ctx, &coin)
    require.NoError(t, err)
    assert.Greater(t, id, int64(0))
    assert.Equal(t, id, coin.ID)
    
    // Test update
    coin.Quantity = 0.2
    id, err = repo.Store(ctx, &coin)
    require.NoError(t, err)
    assert.Equal(t, coin.ID, id)
    
    // Verify update
    foundCoin, err := repo.FindByID(ctx, id)
    require.NoError(t, err)
    require.NotNil(t, foundCoin)
    assert.Equal(t, 0.2, foundCoin.Quantity)
}

func TestBoughtCoinRepository_FindAll(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    repo := sqlite.NewBoughtCoinRepository(db)
    ctx := context.Background()
    
    // Insert test data
    coins := []models.BoughtCoin{
        {
            Symbol:        "BTCUSDT",
            PurchasePrice: 50000.0,
            Quantity:      0.1,
        },
        {
            Symbol:        "ETHUSDT",
            PurchasePrice: 3000.0,
            Quantity:      1.0,
        },
        {
            Symbol:        "SOLUSDT",
            PurchasePrice: 100.0,
            Quantity:      5.0,
            IsDeleted:     true,
        },
    }
    
    for i := range coins {
        _, err := repo.Store(ctx, &coins[i])
        require.NoError(t, err)
    }
    
    // Test FindAll without deleted
    foundCoins, err := repo.FindAll(ctx, false)
    require.NoError(t, err)
    assert.Len(t, foundCoins, 2)
    
    // Test FindAll with deleted
    foundCoins, err = repo.FindAll(ctx, true)
    require.NoError(t, err)
    assert.Len(t, foundCoins, 3)
}

func TestBoughtCoinRepository_Delete(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    repo := sqlite.NewBoughtCoinRepository(db)
    ctx := context.Background()
    
    // Insert test data
    coin := models.BoughtCoin{
        Symbol:        "BTCUSDT",
        PurchasePrice: 50000.0,
        Quantity:      0.1,
    }
    
    id, err := repo.Store(ctx, &coin)
    require.NoError(t, err)
    
    // Test Delete
    err = repo.Delete(ctx, id)
    require.NoError(t, err)
    
    // Verify deletion
    foundCoin, err := repo.FindByID(ctx, id)
    require.NoError(t, err)
    require.NotNil(t, foundCoin)
    assert.True(t, foundCoin.IsDeleted)
    assert.NotNil(t, foundCoin.SoldAt)
}
```

## 7. Transaction Management Example

For operations that require multiple database operations to be atomic, use transactions:

```go
// Example of using transactions with repositories
package service

import (
    "context"
    "errors"
    "fmt"

    "github.com/jmoiron/sqlx"
    
    "github.com/ryanlisse/cryptobot-backend/internal/domain/models"
    "github.com/ryanlisse/cryptobot-backend/internal/domain/repository"
    "github.com/ryanlisse/cryptobot-backend/internal/platform/database"
)

// TradeService handles trading operations
type TradeService struct {
    db            *sqlx.DB
    boughtCoinRepo repository.BoughtCoinRepository
    positionRepo   repository.PositionRepository
}

// NewTradeService creates a new trade service
func NewTradeService(
    db *sqlx.DB,
    boughtCoinRepo repository.BoughtCoinRepository,
    positionRepo repository.PositionRepository,
) *TradeService {
    return &TradeService{
        db:            db,
        boughtCoinRepo: boughtCoinRepo,
        positionRepo:   positionRepo,
    }
}

// CompleteTradeWithPosition creates a bought coin and position record atomically
func (s *TradeService) CompleteTradeWithPosition(
    ctx context.Context,
    symbol string,
    price float64,
    quantity float64,
    stopLoss *float64,
    takeProfit *float64,
) (int64, int64, error) {
    var boughtCoinID, positionID int64
    
    // Use transaction to ensure both operations succeed or fail together
    err := database.RunInTransaction(ctx, s.db, func(tx *sqlx.Tx) error {
        // Create bought coin
        boughtCoin := models.BoughtCoin{
            Symbol:        symbol,
            PurchasePrice: price,
            Quantity:      quantity,
            StopLoss:      stopLoss,
            TakeProfit:    takeProfit,
        }
        
        var err error
        boughtCoinID, err = s.boughtCoinRepo.Store(ctx, &boughtCoin)
        if err != nil {
            return fmt.Errorf("failed to store bought coin: %w", err)
        }
        
        // Create position
        position := models.Position{
            Symbol:     symbol,
            EntryPrice: price,
            Quantity:   quantity,
            Status:     "open",
            StopLoss:   stopLoss,
            TakeProfit: takeProfit,
        }
        
        positionID, err = s.positionRepo.Store(ctx, &position)
        if err != nil {
            return fmt.Errorf("failed to store position: %w", err)
        }
        
        return nil
    })
    
    if err != nil {
        return 0, 0, err
    }
    
    return boughtCoinID, positionID, nil
}
```

For more information on the next implementation layers, refer to:
- [06a-newcoin-service.md](06a-newcoin-service.md) - New coin detection service
- [06b-trade-service.md](06b-trade-service.md) - Trading service implementation
- [06c-portfolio-service.md](06c-portfolio-service.md) - Portfolio management service

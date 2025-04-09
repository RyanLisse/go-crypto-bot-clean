# Position Management

## Overview

The Position Management module is responsible for the lifecycle management of trading positions, ensuring disciplined trade management with appropriate risk controls. It handles position entry, scaling, stop-loss, take-profit, and trailing stop adjustments in accordance with the trading strategy.

## Component Structure

```
internal/domain/position/
├── management/
│   ├── position.go      # Core position management logic
│   ├── scaling.go       # Position scaling implementation
│   ├── stoploss.go      # Stop-loss management
│   └── takeprofit.go    # Take-profit management
├── types.go             # Position-related type definitions
└── service.go           # Position service interface
```

## Core Interfaces and Types

```go
// Position represents a trading position
type Position struct {
    ID            string    `json:"id"`
    Symbol        string    `json:"symbol"`
    Quantity      float64   `json:"quantity"`
    EntryPrice    float64   `json:"entry_price"`
    CurrentPrice  float64   `json:"current_price"`
    StopLoss      float64   `json:"stop_loss"`
    TakeProfit    float64   `json:"take_profit"`
    TrailingStop  *float64  `json:"trailing_stop,omitempty"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
    PnL           float64   `json:"pnl"`
    PnLPercentage float64   `json:"pnl_percentage"`
    Status        string    `json:"status"` // open, closed
    Orders        []Order   `json:"orders"` // Entry and scaling orders
}

// PositionService defines the interface for position management
type PositionService interface {
    // Entry methods
    EnterPosition(ctx context.Context, order Order) (*Position, error)
    ScalePosition(ctx context.Context, positionID string, order Order) (*Position, error)
    
    // Exit methods
    ExitPosition(ctx context.Context, positionID string, price float64) error
    
    // Management methods
    GetPosition(ctx context.Context, positionID string) (*Position, error)
    GetPositions(ctx context.Context, filter PositionFilter) ([]*Position, error)
    UpdateStopLoss(ctx context.Context, positionID string, price float64) error
    UpdateTakeProfit(ctx context.Context, positionID string, price float64) error
    UpdateTrailingStop(ctx context.Context, positionID string, offset float64) error
    
    // Monitoring methods
    CheckPositions(ctx context.Context) error
}

// Order represents a trade order associated with a position
type Order struct {
    ID        string    `json:"id"`
    Symbol    string    `json:"symbol"`
    Type      string    `json:"type"` // market, limit
    Side      string    `json:"side"` // buy, sell
    Price     float64   `json:"price"`
    Quantity  float64   `json:"quantity"`
    Status    string    `json:"status"` // new, filled, canceled
    CreatedAt time.Time `json:"created_at"`
    FilledAt  time.Time `json:"filled_at,omitempty"`
}

// PositionFilter used for filtering positions in queries
type PositionFilter struct {
    Symbol  string
    Status  string
    MinPnL  *float64
    MaxPnL  *float64
    FromDate *time.Time
    ToDate  *time.Time
}
```

## Key Implementation Components

### 1. Position Entry and Scaling

```go
// PositionManager implements the PositionService interface
type PositionManager struct {
    positionRepo repositories.PositionRepository
    orderService service.OrderService
    priceService service.PriceService
    logger       log.Logger
}

// EnterPosition creates a new position from an order
func (pm *PositionManager) EnterPosition(ctx context.Context, order Order) (*Position, error) {
    // Validate order
    if err := pm.validateOrder(order); err != nil {
        return nil, fmt.Errorf("invalid order: %w", err)
    }
    
    // Create position
    position := &Position{
        ID:         uuid.New().String(),
        Symbol:     order.Symbol,
        Quantity:   order.Quantity,
        EntryPrice: order.Price,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
        Status:     "open",
        Orders:     []Order{order},
    }
    
    // Set default stop-loss and take-profit if not provided
    if position.StopLoss == 0 {
        position.StopLoss = order.Price * 0.95 // 5% default stop-loss
    }
    
    if position.TakeProfit == 0 {
        position.TakeProfit = order.Price * 1.15 // 15% default take-profit
    }
    
    // Save position
    savedPosition, err := pm.positionRepo.Create(ctx, position)
    if err != nil {
        return nil, fmt.Errorf("failed to save position: %w", err)
    }
    
    pm.logger.Info("New position opened", 
        "position_id", position.ID, 
        "symbol", position.Symbol,
        "quantity", position.Quantity,
        "entry_price", position.EntryPrice)
    
    return savedPosition, nil
}

// ScalePosition adds to an existing position
func (pm *PositionManager) ScalePosition(ctx context.Context, positionID string, order Order) (*Position, error) {
    // Get existing position
    position, err := pm.positionRepo.GetByID(ctx, positionID)
    if err != nil {
        return nil, fmt.Errorf("failed to get position: %w", err)
    }
    
    // Ensure position is open
    if position.Status != "open" {
        return nil, errors.New("cannot scale a closed position")
    }
    
    // Ensure order is for the same symbol
    if position.Symbol != order.Symbol {
        return nil, errors.New("order symbol does not match position symbol")
    }
    
    // Calculate new average entry price and total quantity
    totalCost := position.EntryPrice * position.Quantity
    additionalCost := order.Price * order.Quantity
    totalQuantity := position.Quantity + order.Quantity
    
    // Update position
    position.EntryPrice = (totalCost + additionalCost) / totalQuantity
    position.Quantity = totalQuantity
    position.UpdatedAt = time.Now()
    position.Orders = append(position.Orders, order)
    
    // Save updated position
    updatedPosition, err := pm.positionRepo.Update(ctx, position)
    if err != nil {
        return nil, fmt.Errorf("failed to update position: %w", err)
    }
    
    pm.logger.Info("Position scaled", 
        "position_id", position.ID, 
        "new_quantity", position.Quantity,
        "new_entry_price", position.EntryPrice)
    
    return updatedPosition, nil
}
```

### 2. Stop-Loss and Take-Profit Management

```go
// CheckPositions monitors all open positions for stop-loss and take-profit triggers
func (pm *PositionManager) CheckPositions(ctx context.Context) error {
    // Get all open positions
    positions, err := pm.positionRepo.GetByFilter(ctx, PositionFilter{Status: "open"})
    if err != nil {
        return fmt.Errorf("failed to get open positions: %w", err)
    }
    
    for _, position := range positions {
        // Get current price for symbol
        currentPrice, err := pm.priceService.GetPrice(ctx, position.Symbol)
        if err != nil {
            pm.logger.Error("Failed to get current price", 
                "position_id", position.ID, 
                "symbol", position.Symbol,
                "error", err)
            continue
        }
        
        // Update position with current price and P&L
        position.CurrentPrice = currentPrice
        position.PnL = (currentPrice - position.EntryPrice) * position.Quantity
        position.PnLPercentage = (currentPrice - position.EntryPrice) / position.EntryPrice * 100
        
        // Check for trailing stop adjustment
        if position.TrailingStop != nil {
            pm.adjustTrailingStop(position, currentPrice)
        }
        
        // Check stop-loss
        if currentPrice <= position.StopLoss {
            pm.logger.Info("Stop-loss triggered", 
                "position_id", position.ID, 
                "symbol", position.Symbol,
                "stop_price", position.StopLoss,
                "current_price", currentPrice)
                
            if err := pm.ExitPosition(ctx, position.ID, currentPrice); err != nil {
                pm.logger.Error("Failed to exit position at stop-loss", 
                    "position_id", position.ID, 
                    "error", err)
            }
            continue
        }
        
        // Check take-profit
        if currentPrice >= position.TakeProfit {
            pm.logger.Info("Take-profit triggered", 
                "position_id", position.ID, 
                "symbol", position.Symbol,
                "take_profit", position.TakeProfit,
                "current_price", currentPrice)
                
            if err := pm.ExitPosition(ctx, position.ID, currentPrice); err != nil {
                pm.logger.Error("Failed to exit position at take-profit", 
                    "position_id", position.ID, 
                    "error", err)
            }
            continue
        }
        
        // Update position
        if _, err := pm.positionRepo.Update(ctx, position); err != nil {
            pm.logger.Error("Failed to update position", 
                "position_id", position.ID, 
                "error", err)
        }
    }
    
    return nil
}
```

### 3. Trailing Stop Management

```go
// adjustTrailingStop updates a trailing stop based on price movement
func (pm *PositionManager) adjustTrailingStop(position *Position, currentPrice float64) {
    if position.TrailingStop == nil {
        return
    }
    
    // For long positions, trailing stop should only move up
    trailingStopPrice := position.EntryPrice * (1 - *position.TrailingStop/100)
    
    // If current price has moved up enough, adjust the stop-loss
    potentialNewStop := currentPrice * (1 - *position.TrailingStop/100)
    
    // Only update if the new stop would be higher than the current one
    if potentialNewStop > position.StopLoss {
        position.StopLoss = potentialNewStop
        
        pm.logger.Info("Trailing stop adjusted", 
            "position_id", position.ID, 
            "symbol", position.Symbol,
            "new_stop_loss", position.StopLoss,
            "current_price", currentPrice)
    }
}
```

## Integration with Trade Service

The Position Management module integrates with the Trade Service to manage positions created from trades:

```go
// In the TradeService
type TradeService struct {
    // ...other fields
    positionService position.PositionService
}

func (s *TradeService) ExecutePurchase(ctx context.Context, symbol string, amount float64) (*models.BoughtCoin, error) {
    // Get price and create order
    price, err := s.priceService.GetPrice(ctx, symbol)
    if err != nil {
        return nil, err
    }
    
    quantity := amount / price
    
    // Create order
    order := position.Order{
        ID:        uuid.New().String(),
        Symbol:    symbol,
        Type:      "market",
        Side:      "buy",
        Price:     price,
        Quantity:  quantity,
        Status:    "filled",
        CreatedAt: time.Now(),
        FilledAt:  time.Now(),
    }
    
    // Enter position
    pos, err := s.positionService.EnterPosition(ctx, order)
    if err != nil {
        return nil, fmt.Errorf("failed to create position: %w", err)
    }
    
    // Create bought coin record
    boughtCoin := &models.BoughtCoin{
        Symbol:        symbol,
        PurchasePrice: price,
        Quantity:      quantity,
        PurchaseTime:  time.Now(),
        StopLossPrice: pos.StopLoss,
    }
    
    // Store bought coin
    boughtCoin, err = s.boughtCoinRepo.Store(ctx, boughtCoin)
    if err != nil {
        return nil, fmt.Errorf("failed to store bought coin: %w", err)
    }
    
    return boughtCoin, nil
}
```

## Configuration Options

The Position Management module supports several configuration options:

- **Default stop-loss percentage**: Sets the default stop-loss level for new positions
- **Default take-profit percentage**: Sets the default take-profit level for new positions
- **Trailing stop activation threshold**: Determines when to enable trailing stops
- **Position size limits**: Configures maximum position size for risk management
- **Scaling rules**: Defines how positions can be scaled (e.g., max scaling, avg price limits)

## Testing Approach

- **Unit tests** for each position management function
- **Table-driven tests** for various scenarios (stop-loss/take-profit triggers, scaling calculations)
- **Mock repositories** for testing service logic independently
- **Integration tests** with price service and order execution

## Security Considerations

- Validation of all inputs (prices, quantities, order parameters)
- Prevention of invalid operations (e.g., scaling closed positions)
- Concurrency management for position updates
- Transaction handling to prevent partial updates

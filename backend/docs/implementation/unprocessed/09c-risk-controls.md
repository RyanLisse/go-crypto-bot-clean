# Risk Controls

## Overview

The Risk Controls module implements essential risk management features to protect trading capital and ensure the bot operates within defined risk parameters. It handles drawdown monitoring, position sizing, exposure limits, and daily loss caps to prevent excessive risk-taking.

## Component Structure

```
internal/domain/risk/
├── controls/
│   ├── drawdown.go         # Drawdown calculation and monitoring
│   ├── exposure.go         # Total exposure tracking and limits
│   ├── position_sizing.go  # Safe position size calculation
│   └── daily_limits.go     # Daily loss limits enforcement
├── types.go                # Risk-related type definitions
└── service.go              # Risk service interface
```

## Core Interfaces and Types

```go
// RiskParameters contains risk management configuration
type RiskParameters struct {
    MaxDrawdownPercent float64 `json:"max_drawdown_percent"`
    RiskPerTradePercent float64 `json:"risk_per_trade_percent"`
    MaxExposurePercent float64 `json:"max_exposure_percent"`
    DailyLossLimitPercent float64 `json:"daily_loss_limit_percent"`
    MinAccountBalance float64 `json:"min_account_balance"`
}

// RiskService defines the interface for risk management
type RiskService interface {
    // Calculation methods
    CalculatePositionSize(ctx context.Context, symbol string, accountBalance float64) (float64, error)
    CalculateDrawdown(ctx context.Context) (float64, error)
    
    // Limit checking methods
    CheckExposureLimit(ctx context.Context, newOrderValue float64) (bool, error)
    CheckDailyLossLimit(ctx context.Context) (bool, error)
    
    // Risk status methods
    GetRiskStatus(ctx context.Context) (*RiskStatus, error)
    UpdateRiskParameters(ctx context.Context, params RiskParameters) error
    
    // Monitoring methods
    IsTradeAllowed(ctx context.Context, symbol string, orderValue float64) (bool, string, error)
}

// RiskStatus represents the current risk metrics
type RiskStatus struct {
    CurrentDrawdown float64 `json:"current_drawdown"`
    TotalExposure float64 `json:"total_exposure"`
    TodayPnL float64 `json:"today_pnl"`
    AccountBalance float64 `json:"account_balance"`
    TradingEnabled bool `json:"trading_enabled"`
    DisabledReason string `json:"disabled_reason,omitempty"`
}

// BalanceHistory tracks account balance over time for drawdown calculation
type BalanceHistory struct {
    ID int64 `json:"id"`
    Balance float64 `json:"balance"`
    Timestamp time.Time `json:"timestamp"`
}
```

## Key Implementation Components

### 1. Position Sizing

Position sizing ensures that each trade risks an appropriate amount of capital:

```go
// RiskManager implements the RiskService interface
type RiskManager struct {
    riskParams        RiskParameters
    balanceRepo       repositories.BalanceHistoryRepository
    positionRepo      repositories.PositionRepository
    tradeRepo         repositories.TradeRepository
    accountService    service.AccountService
    priceService      service.PriceService
    lock              sync.RWMutex
    logger            log.Logger
}

// CalculatePositionSize determines a safe position size based on risk parameters
func (rm *RiskManager) CalculatePositionSize(ctx context.Context, symbol string, accountBalance float64) (float64, error) {
    rm.lock.RLock()
    riskPercent := rm.riskParams.RiskPerTradePercent
    rm.lock.RUnlock()
    
    // Get current price and determine stop-loss placement
    currentPrice, err := rm.priceService.GetPrice(ctx, symbol)
    if err != nil {
        return 0, fmt.Errorf("failed to get price for %s: %w", symbol, err)
    }
    
    // Default stop-loss at 5% below entry
    stopLossPrice := currentPrice * 0.95
    
    // Calculate risk amount based on account balance
    riskAmount := accountBalance * (riskPercent / 100)
    
    // Calculate position size
    priceDifference := currentPrice - stopLossPrice
    riskPerUnit := priceDifference
    
    if riskPerUnit <= 0 {
        return 0, errors.New("invalid stop-loss placement, risk per unit is zero or negative")
    }
    
    // Position size = risk amount / risk per unit
    positionSize := riskAmount / riskPerUnit
    
    // Convert to coin quantity based on price
    quantity := positionSize / currentPrice
    
    rm.logger.Info("Calculated position size",
        "symbol", symbol,
        "account_balance", accountBalance,
        "risk_percent", riskPercent,
        "risk_amount", riskAmount,
        "quantity", quantity)
    
    return quantity, nil
}
```

### 2. Drawdown Calculation and Monitoring

```go
// CalculateDrawdown computes the maximum peak-to-trough drawdown
func (rm *RiskManager) CalculateDrawdown(ctx context.Context) (float64, error) {
    // Get historical balance data
    history, err := rm.balanceRepo.GetHistory(ctx, 90) // Last 90 days
    if err != nil {
        return 0, fmt.Errorf("failed to get balance history: %w", err)
    }
    
    if len(history) < 2 {
        return 0, nil // Not enough data to calculate drawdown
    }
    
    // Find peak and calculate drawdown
    var maxDrawdown float64
    var peak float64
    
    for _, entry := range history {
        if entry.Balance > peak {
            peak = entry.Balance
        }
        
        if peak > 0 {
            drawdown := (peak - entry.Balance) / peak
            if drawdown > maxDrawdown {
                maxDrawdown = drawdown
            }
        }
    }
    
    return maxDrawdown, nil
}

// CheckDrawdownLimit verifies if trading should be allowed based on drawdown
func (rm *RiskManager) checkDrawdownLimit(ctx context.Context) (bool, error) {
    drawdown, err := rm.CalculateDrawdown(ctx)
    if err != nil {
        return false, err
    }
    
    rm.lock.RLock()
    maxAllowed := rm.riskParams.MaxDrawdownPercent / 100
    rm.lock.RUnlock()
    
    allowed := drawdown < maxAllowed
    
    if !allowed {
        rm.logger.Warn("Trading disabled due to drawdown limit",
            "current_drawdown", drawdown,
            "max_allowed", maxAllowed)
    }
    
    return allowed, nil
}
```

### 3. Exposure and Daily Loss Limits

```go
// CheckExposureLimit verifies if a new order would exceed exposure limits
func (rm *RiskManager) CheckExposureLimit(ctx context.Context, newOrderValue float64) (bool, error) {
    // Get current account balance
    accountBalance, err := rm.accountService.GetBalance(ctx)
    if err != nil {
        return false, fmt.Errorf("failed to get account balance: %w", err)
    }
    
    // Get current positions to calculate total exposure
    positions, err := rm.positionRepo.GetByFilter(ctx, position.PositionFilter{Status: "open"})
    if err != nil {
        return false, fmt.Errorf("failed to get open positions: %w", err)
    }
    
    // Calculate total exposure
    var totalExposure float64
    for _, pos := range positions {
        totalExposure += pos.Quantity * pos.EntryPrice
    }
    
    // Add new order value
    potentialExposure := totalExposure + newOrderValue
    
    rm.lock.RLock()
    maxExposurePercent := rm.riskParams.MaxExposurePercent
    rm.lock.RUnlock()
    
    // Calculate maximum allowed exposure
    maxExposure := accountBalance * (maxExposurePercent / 100)
    
    // Check if new total exposure exceeds limit
    allowed := potentialExposure <= maxExposure
    
    if !allowed {
        rm.logger.Warn("Order rejected due to exposure limit",
            "current_exposure", totalExposure,
            "new_order_value", newOrderValue,
            "potential_exposure", potentialExposure,
            "max_allowed", maxExposure)
    }
    
    return allowed, nil
}

// CheckDailyLossLimit verifies if trading should be allowed based on daily P&L
func (rm *RiskManager) CheckDailyLossLimit(ctx context.Context) (bool, error) {
    // Get today's closed trades
    today := time.Now().UTC().Truncate(24 * time.Hour)
    trades, err := rm.tradeRepo.GetByFilter(ctx, repositories.TradeFilter{
        FromDate: &today,
        Status:   "closed",
    })
    if err != nil {
        return false, fmt.Errorf("failed to get today's trades: %w", err)
    }
    
    // Calculate today's P&L
    var todayPnL float64
    for _, trade := range trades {
        todayPnL += trade.PnL
    }
    
    // Get account balance
    accountBalance, err := rm.accountService.GetBalance(ctx)
    if err != nil {
        return false, fmt.Errorf("failed to get account balance: %w", err)
    }
    
    rm.lock.RLock()
    dailyLossLimitPercent := rm.riskParams.DailyLossLimitPercent
    rm.lock.RUnlock()
    
    // Calculate maximum allowed daily loss
    maxDailyLoss := accountBalance * (dailyLossLimitPercent / 100)
    
    // Check if today's losses exceed the limit
    allowed := todayPnL >= -maxDailyLoss
    
    if !allowed {
        rm.logger.Warn("Trading disabled due to daily loss limit",
            "today_pnl", todayPnL,
            "max_daily_loss", maxDailyLoss)
    }
    
    return allowed, nil
}
```

### 4. Comprehensive Trade Permission Check

```go
// IsTradeAllowed performs a comprehensive check of all risk controls
func (rm *RiskManager) IsTradeAllowed(ctx context.Context, symbol string, orderValue float64) (bool, string, error) {
    // Check account minimum balance
    accountBalance, err := rm.accountService.GetBalance(ctx)
    if err != nil {
        return false, "", fmt.Errorf("failed to get account balance: %w", err)
    }
    
    rm.lock.RLock()
    minAccountBalance := rm.riskParams.MinAccountBalance
    rm.lock.RUnlock()
    
    if accountBalance < minAccountBalance {
        reason := fmt.Sprintf("Account balance below minimum: %.2f < %.2f", accountBalance, minAccountBalance)
        return false, reason, nil
    }
    
    // Check drawdown limit
    drawdownAllowed, err := rm.checkDrawdownLimit(ctx)
    if err != nil {
        return false, "", err
    }
    
    if !drawdownAllowed {
        reason := "Maximum drawdown limit reached"
        return false, reason, nil
    }
    
    // Check exposure limit
    exposureAllowed, err := rm.CheckExposureLimit(ctx, orderValue)
    if err != nil {
        return false, "", err
    }
    
    if !exposureAllowed {
        reason := "Maximum exposure limit would be exceeded"
        return false, reason, nil
    }
    
    // Check daily loss limit
    dailyLossAllowed, err := rm.CheckDailyLossLimit(ctx)
    if err != nil {
        return false, "", err
    }
    
    if !dailyLossAllowed {
        reason := "Daily loss limit reached"
        return false, reason, nil
    }
    
    return true, "", nil
}
```

## Integration with Trade Service

The Risk Controls module integrates with the Trade Service to ensure all trades adhere to risk limits:

```go
// In the TradeService
type TradeService struct {
    // ...other fields
    riskService risk.RiskService
}

func (s *TradeService) ExecutePurchase(ctx context.Context, symbol string, amount float64) (*models.BoughtCoin, error) {
    // Get current price
    price, err := s.priceService.GetPrice(ctx, symbol)
    if err != nil {
        return nil, err
    }
    
    orderValue := amount * price
    
    // Check if this trade is allowed by risk controls
    allowed, reason, err := s.riskService.IsTradeAllowed(ctx, symbol, orderValue)
    if err != nil {
        return nil, fmt.Errorf("failed to check risk controls: %w", err)
    }
    
    if !allowed {
        s.logger.Warn("Trade rejected by risk controls",
            "symbol", symbol,
            "amount", amount,
            "reason", reason)
        return nil, fmt.Errorf("trade rejected: %s", reason)
    }
    
    // Get position size from risk service
    accountBalance, err := s.accountService.GetBalance(ctx)
    if err != nil {
        return nil, err
    }
    
    // If a specific amount wasn't requested, calculate safe position size
    if amount <= 0 {
        quantity, err := s.riskService.CalculatePositionSize(ctx, symbol, accountBalance)
        if err != nil {
            return nil, fmt.Errorf("failed to calculate position size: %w", err)
        }
        amount = quantity * price
    }
    
    // Proceed with purchase...
    // [rest of purchase implementation]
}
```

## Configuration Options

The Risk Controls module is highly configurable via the `RiskParameters` structure:

```go
// Default risk parameters
defaultRiskParams := RiskParameters{
    MaxDrawdownPercent:    20.0,  // Maximum 20% drawdown before halting trading
    RiskPerTradePercent:   1.0,   // Risk 1% of account per trade
    MaxExposurePercent:    50.0,  // Maximum 50% of account in open positions
    DailyLossLimitPercent: 5.0,   // Stop trading if daily losses exceed 5%
    MinAccountBalance:     100.0, // Minimum account balance required
}
```

These parameters can be loaded from configuration files, environment variables, or adjusted dynamically based on market conditions.

## Testing Approach

- **Unit tests** for each risk calculation function
- **Table-driven tests** for limit checks with various scenarios
- **Mock repositories** for testing service logic independently
- **Integration tests** with account and position services

## Security Considerations

- Validation of all inputs (balances, parameters, order values)
- Concurrency management using mutex locks
- Defensive error handling
- Fallback to conservative defaults if configuration is missing

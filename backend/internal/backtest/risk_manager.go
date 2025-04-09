package backtest

import (
	"context"
	"fmt"
	"math"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RiskManagerConfig contains configuration options for the risk manager
type RiskManagerConfig struct {
	// Maximum percentage of capital to risk per trade
	MaxRiskPerTrade float64

	// Maximum percentage of capital to allocate to a single position
	MaxPositionSize float64

	// Maximum percentage of capital to allocate across all positions
	MaxTotalExposure float64

	// Maximum number of concurrent positions
	MaxPositions int

	// Maximum percentage drawdown before stopping trading
	MaxDrawdown float64

	// Maximum daily loss percentage before stopping trading
	MaxDailyLoss float64

	// Whether to use trailing stops
	UseTrailingStops bool

	// Trailing stop percentage
	TrailingStopPercent float64

	// Whether to use take profit levels
	UseTakeProfits bool

	// Take profit levels (percentage of entry price)
	TakeProfitLevels []float64

	// Percentage of position to close at each take profit level
	TakeProfitSizes []float64

	// Logger
	Logger *zap.Logger

	// Database connection
	DB *gorm.DB
}

// EventDrivenRiskManager handles risk management for the event-driven backtesting engine
type EventDrivenRiskManager struct {
	config            *RiskManagerConfig
	initialCapital    float64
	currentCapital    float64
	openPositions     map[string]*models.Position
	dailyPnL          map[string]float64 // Key is date in YYYY-MM-DD format
	maxDrawdown       float64
	currentDrawdown   float64
	highWaterMark     float64
	logger            *zap.Logger
	db                *gorm.DB
	trailingStops     map[string]float64                  // Key is position ID
	takeProfitLevels  map[string][]models.TakeProfitLevel // Key is position ID
	positionSizeCache map[string]float64                  // Key is symbol
}

// NewEventDrivenRiskManager creates a new risk manager for event-driven backtesting
func NewEventDrivenRiskManager(config *RiskManagerConfig, initialCapital float64, db *gorm.DB) *EventDrivenRiskManager {
	logger := config.Logger
	if logger == nil {
		logger, _ = zap.NewDevelopment()
	}

	return &EventDrivenRiskManager{
		config:            config,
		initialCapital:    initialCapital,
		currentCapital:    initialCapital,
		openPositions:     make(map[string]*models.Position),
		dailyPnL:          make(map[string]float64),
		maxDrawdown:       0,
		currentDrawdown:   0,
		highWaterMark:     initialCapital,
		logger:            logger,
		db:                db,
		trailingStops:     make(map[string]float64),
		takeProfitLevels:  make(map[string][]models.TakeProfitLevel),
		positionSizeCache: make(map[string]float64),
	}
}

// CalculatePositionSize calculates the position size based on risk parameters
func (r *EventDrivenRiskManager) CalculatePositionSize(ctx context.Context, symbol string, entryPrice float64, stopLossPrice float64) (float64, error) {
	// Check if we have a cached position size for this symbol
	if size, ok := r.positionSizeCache[symbol]; ok {
		return size, nil
	}

	// Calculate risk per trade in currency units
	riskAmount := r.currentCapital * r.config.MaxRiskPerTrade / 100.0

	// Calculate position size based on risk and stop loss
	var positionSize float64
	if stopLossPrice > 0 && math.Abs(entryPrice-stopLossPrice) > 0 {
		// Risk per pip
		riskPerPip := riskAmount / math.Abs(entryPrice-stopLossPrice)
		positionSize = riskPerPip
	} else {
		// If no stop loss is provided, use the max position size
		positionSize = r.currentCapital * r.config.MaxPositionSize / 100.0 / entryPrice
	}

	// Check if position size exceeds max position size
	maxPositionSize := r.currentCapital * r.config.MaxPositionSize / 100.0 / entryPrice
	if positionSize > maxPositionSize {
		positionSize = maxPositionSize
	}

	// Check if adding this position would exceed max total exposure
	currentExposure := 0.0
	for _, position := range r.openPositions {
		currentExposure += position.Quantity * position.EntryPrice
	}

	maxExposure := r.currentCapital * r.config.MaxTotalExposure / 100.0
	if currentExposure+positionSize*entryPrice > maxExposure {
		// Adjust position size to fit within max exposure
		positionSize = (maxExposure - currentExposure) / entryPrice
		if positionSize <= 0 {
			return 0, fmt.Errorf("cannot open position: would exceed maximum exposure")
		}
	}

	// Check if we would exceed max positions
	if len(r.openPositions) >= r.config.MaxPositions {
		return 0, fmt.Errorf("cannot open position: would exceed maximum number of positions")
	}

	// Check if we're in a drawdown exceeding max drawdown
	if r.currentDrawdown > r.config.MaxDrawdown {
		return 0, fmt.Errorf("cannot open position: current drawdown (%.2f%%) exceeds maximum (%.2f%%)", r.currentDrawdown, r.config.MaxDrawdown)
	}

	// Check if we've exceeded max daily loss
	dateStr := time.Now().Format("2006-01-02")
	if dailyLoss, ok := r.dailyPnL[dateStr]; ok && dailyLoss < 0 {
		dailyLossPercent := math.Abs(dailyLoss) / r.initialCapital * 100.0
		if dailyLossPercent > r.config.MaxDailyLoss {
			return 0, fmt.Errorf("cannot open position: daily loss (%.2f%%) exceeds maximum (%.2f%%)", dailyLossPercent, r.config.MaxDailyLoss)
		}
	}

	// Cache the position size
	r.positionSizeCache[symbol] = positionSize

	return positionSize, nil
}

// OnPositionOpened is called when a position is opened
func (r *EventDrivenRiskManager) OnPositionOpened(ctx context.Context, position *models.Position) error {
	// Add position to open positions
	r.openPositions[position.ID] = position

	// Set up trailing stop if enabled
	if r.config.UseTrailingStops {
		r.trailingStops[position.ID] = position.EntryPrice * (1.0 - r.config.TrailingStopPercent/100.0)
	}

	// Set up take profit levels if enabled
	if r.config.UseTakeProfits {
		takeProfitLevels := make([]models.TakeProfitLevel, len(r.config.TakeProfitLevels))
		for i, level := range r.config.TakeProfitLevels {
			size := 0.0
			if i < len(r.config.TakeProfitSizes) {
				size = r.config.TakeProfitSizes[i]
			} else {
				// Default to equal distribution of remaining size
				remainingLevels := len(r.config.TakeProfitLevels) - i
				remainingSize := 100.0
				for j := 0; j < i; j++ {
					if j < len(r.config.TakeProfitSizes) {
						remainingSize -= r.config.TakeProfitSizes[j]
					}
				}
				size = remainingSize / float64(remainingLevels)
			}

			takeProfitLevels[i] = models.TakeProfitLevel{
				ID:          fmt.Sprintf("%s-%d", position.ID, i),
				PositionID:  position.ID,
				Level:       i + 1,
				Price:       position.EntryPrice * (1.0 + level/100.0),
				Percentage:  level,
				Quantity:    position.Quantity * size / 100.0,
				QuantityPct: size,
				Triggered:   false,
			}
		}
		r.takeProfitLevels[position.ID] = takeProfitLevels

		// Save take profit levels to database
		for _, level := range takeProfitLevels {
			err := r.db.Create(&level).Error
			if err != nil {
				r.logger.Error("Error saving take profit level to database",
					zap.Error(err),
					zap.String("position_id", position.ID),
					zap.Float64("price", level.Price),
				)
			}
		}
	}

	// Update capital
	r.currentCapital -= position.EntryPrice * position.Quantity

	// Save position to database
	err := r.db.Create(position).Error
	if err != nil {
		r.logger.Error("Error saving position to database",
			zap.Error(err),
			zap.String("position_id", position.ID),
			zap.String("symbol", position.Symbol),
		)
	}

	return nil
}

// OnPositionClosed is called when a position is closed
func (r *EventDrivenRiskManager) OnPositionClosed(ctx context.Context, position *models.ClosedPosition) error {
	// Remove position from open positions
	delete(r.openPositions, position.ID)

	// Remove trailing stop
	delete(r.trailingStops, position.ID)

	// Remove take profit levels
	delete(r.takeProfitLevels, position.ID)

	// Update capital
	r.currentCapital += position.ExitPrice * position.Quantity

	// Update daily P&L
	dateStr := position.CloseTime.Format("2006-01-02")
	r.dailyPnL[dateStr] += position.ProfitLoss

	// Update high water mark and drawdown
	if r.currentCapital > r.highWaterMark {
		r.highWaterMark = r.currentCapital
	}
	r.currentDrawdown = (r.highWaterMark - r.currentCapital) / r.highWaterMark * 100.0
	if r.currentDrawdown > r.maxDrawdown {
		r.maxDrawdown = r.currentDrawdown
	}

	// Save closed position to database
	err := r.db.Create(position).Error
	if err != nil {
		r.logger.Error("Error saving closed position to database",
			zap.Error(err),
			zap.String("position_id", position.ID),
			zap.String("symbol", position.Symbol),
		)
	}

	// Clear position size cache for this symbol
	delete(r.positionSizeCache, position.Symbol)

	return nil
}

// OnPriceUpdate is called when a price update is received
func (r *EventDrivenRiskManager) OnPriceUpdate(ctx context.Context, symbol string, price float64, timestamp time.Time) ([]*Signal, error) {
	var signals []*Signal

	// Check trailing stops
	for id, position := range r.openPositions {
		if position.Symbol != symbol {
			continue
		}

		// Check if position has a trailing stop
		if stopPrice, ok := r.trailingStops[id]; ok {
			if position.Side == "BUY" && price <= stopPrice {
				// Trigger trailing stop for long position
				signals = append(signals, &Signal{
					Symbol:    symbol,
					Side:      "SELL",
					Quantity:  position.Quantity,
					Price:     price,
					Timestamp: timestamp,
					Reason:    "Trailing stop triggered",
				})
			} else if position.Side == "SELL" && price >= stopPrice {
				// Trigger trailing stop for short position
				signals = append(signals, &Signal{
					Symbol:    symbol,
					Side:      "BUY",
					Quantity:  position.Quantity,
					Price:     price,
					Timestamp: timestamp,
					Reason:    "Trailing stop triggered",
				})
			} else if position.Side == "BUY" && price > position.EntryPrice+(stopPrice-position.EntryPrice)/r.config.TrailingStopPercent*100.0 {
				// Update trailing stop for long position
				newStopPrice := price * (1.0 - r.config.TrailingStopPercent/100.0)
				if newStopPrice > stopPrice {
					r.trailingStops[id] = newStopPrice
					r.logger.Info("Updated trailing stop",
						zap.String("position_id", id),
						zap.String("symbol", symbol),
						zap.Float64("price", price),
						zap.Float64("new_stop_price", newStopPrice),
					)
				}
			} else if position.Side == "SELL" && price < position.EntryPrice-(position.EntryPrice-stopPrice)/r.config.TrailingStopPercent*100.0 {
				// Update trailing stop for short position
				newStopPrice := price * (1.0 + r.config.TrailingStopPercent/100.0)
				if newStopPrice < stopPrice {
					r.trailingStops[id] = newStopPrice
					r.logger.Info("Updated trailing stop",
						zap.String("position_id", id),
						zap.String("symbol", symbol),
						zap.Float64("price", price),
						zap.Float64("new_stop_price", newStopPrice),
					)
				}
			}
		}

		// Check take profit levels
		if levels, ok := r.takeProfitLevels[id]; ok {
			for i, level := range levels {
				if !level.Triggered {
					if position.Side == "BUY" && price >= level.Price {
						// Trigger take profit for long position
						quantity := level.Quantity
						signals = append(signals, &Signal{
							Symbol:    symbol,
							Side:      "SELL",
							Quantity:  quantity,
							Price:     price,
							Timestamp: timestamp,
							Reason:    fmt.Sprintf("Take profit level %d triggered", i+1),
						})

						// Mark level as triggered
						level.Triggered = true
						r.takeProfitLevels[id][i] = level

						// Update database
						err := r.db.Model(&models.TakeProfitLevel{}).Where("id = ?", level.ID).Update("triggered", true).Error
						if err != nil {
							r.logger.Error("Error updating take profit level in database",
								zap.Error(err),
								zap.String("level_id", level.ID),
							)
						}
					} else if position.Side == "SELL" && price <= level.Price {
						// Trigger take profit for short position
						quantity := level.Quantity
						signals = append(signals, &Signal{
							Symbol:    symbol,
							Side:      "BUY",
							Quantity:  quantity,
							Price:     price,
							Timestamp: timestamp,
							Reason:    fmt.Sprintf("Take profit level %d triggered", i+1),
						})

						// Mark level as triggered
						level.Triggered = true
						r.takeProfitLevels[id][i] = level

						// Update database
						err := r.db.Model(&models.TakeProfitLevel{}).Where("id = ?", level.ID).Update("triggered", true).Error
						if err != nil {
							r.logger.Error("Error updating take profit level in database",
								zap.Error(err),
								zap.String("level_id", level.ID),
							)
						}
					}
				}
			}
		}
	}

	return signals, nil
}

// GetCurrentDrawdown returns the current drawdown percentage
func (r *EventDrivenRiskManager) GetCurrentDrawdown() float64 {
	return r.currentDrawdown
}

// GetMaxDrawdown returns the maximum drawdown percentage
func (r *EventDrivenRiskManager) GetMaxDrawdown() float64 {
	return r.maxDrawdown
}

// GetDailyPnL returns the P&L for a specific date
func (r *EventDrivenRiskManager) GetDailyPnL(date time.Time) float64 {
	dateStr := date.Format("2006-01-02")
	return r.dailyPnL[dateStr]
}

// GetCurrentCapital returns the current capital
func (r *EventDrivenRiskManager) GetCurrentCapital() float64 {
	return r.currentCapital
}

// GetOpenPositions returns all open positions
func (r *EventDrivenRiskManager) GetOpenPositions() []*models.Position {
	positions := make([]*models.Position, 0, len(r.openPositions))
	for _, position := range r.openPositions {
		positions = append(positions, position)
	}
	return positions
}

// GetTrailingStops returns all trailing stops
func (r *EventDrivenRiskManager) GetTrailingStops() map[string]float64 {
	return r.trailingStops
}

// GetTakeProfitLevels returns all take profit levels
func (r *EventDrivenRiskManager) GetTakeProfitLevels() map[string][]models.TakeProfitLevel {
	return r.takeProfitLevels
}

package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"go-crypto-bot-clean/backend/internal/domain/interfaces"
	"go-crypto-bot-clean/backend/internal/domain/models"
)

// SQLitePositionRepository implements the interfaces.PositionRepository interface using SQLite
type SQLitePositionRepository struct {
	db *sqlx.DB
}

// NewSQLitePositionRepository creates a new SQLite position repository
func NewSQLitePositionRepository(db *sqlx.DB) interfaces.PositionRepository {
	return &SQLitePositionRepository{
		db: db,
	}
}

// positionRow represents a position in the database
type positionRow struct {
	ID            string          `db:"id"`
	Symbol        string          `db:"symbol"`
	Quantity      float64         `db:"quantity"`
	EntryPrice    float64         `db:"entry_price"`
	CurrentPrice  float64         `db:"current_price"`
	OpenTime      time.Time       `db:"open_time"`
	StopLoss      float64         `db:"stop_loss"`
	TakeProfit    float64         `db:"take_profit"`
	TrailingStop  sql.NullFloat64 `db:"trailing_stop"`
	CreatedAt     time.Time       `db:"created_at"`
	UpdatedAt     time.Time       `db:"updated_at"`
	PnL           float64         `db:"pnl"`
	PnLPercentage float64         `db:"pnl_percentage"`
	Status        string          `db:"status"`
	OrdersJSON    string          `db:"orders_json"`
}

// toModel converts a positionRow to a Position model
func (r *positionRow) toModel() (*models.Position, error) {
	position := &models.Position{
		ID:            r.ID,
		Symbol:        r.Symbol,
		Quantity:      r.Quantity,
		Amount:        r.Quantity, // For backward compatibility
		EntryPrice:    r.EntryPrice,
		CurrentPrice:  r.CurrentPrice,
		OpenTime:      r.OpenTime,
		OpenedAt:      r.OpenTime, // For backward compatibility
		StopLoss:      r.StopLoss,
		TakeProfit:    r.TakeProfit,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
		PnL:           r.PnL,
		PnLPercentage: r.PnLPercentage,
		Status:        models.PositionStatus(r.Status),
	}

	// Set trailing stop if present
	if r.TrailingStop.Valid {
		trailingStop := r.TrailingStop.Float64
		position.TrailingStop = &trailingStop
	}

	// Unmarshal orders JSON
	if r.OrdersJSON != "" {
		var orders []models.Order
		if err := json.Unmarshal([]byte(r.OrdersJSON), &orders); err != nil {
			return nil, fmt.Errorf("failed to unmarshal orders JSON: %w", err)
		}
		position.Orders = orders
	}

	return position, nil
}

// fromModel converts a Position model to a positionRow
func fromModel(position *models.Position) (*positionRow, error) {
	row := &positionRow{
		ID:            position.ID,
		Symbol:        position.Symbol,
		Quantity:      position.Quantity,
		EntryPrice:    position.EntryPrice,
		CurrentPrice:  position.CurrentPrice,
		OpenTime:      position.OpenTime,
		StopLoss:      position.StopLoss,
		TakeProfit:    position.TakeProfit,
		CreatedAt:     position.CreatedAt,
		UpdatedAt:     position.UpdatedAt,
		PnL:           position.PnL,
		PnLPercentage: position.PnLPercentage,
		Status:        string(position.Status),
	}

	// Set trailing stop if present
	if position.TrailingStop != nil {
		row.TrailingStop = sql.NullFloat64{
			Float64: *position.TrailingStop,
			Valid:   true,
		}
	}

	// Marshal orders to JSON
	if len(position.Orders) > 0 {
		ordersJSON, err := json.Marshal(position.Orders)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal orders to JSON: %w", err)
		}
		row.OrdersJSON = string(ordersJSON)
	}

	return row, nil
}

// FindAll returns all positions matching the filter
func (r *SQLitePositionRepository) FindAll(ctx context.Context, filter interfaces.PositionFilter) ([]*models.Position, error) {
	query := "SELECT * FROM positions WHERE 1=1"
	args := []interface{}{}

	// Apply filters
	if filter.Symbol != "" {
		query += " AND symbol = ?"
		args = append(args, filter.Symbol)
	}
	if filter.Status != "" {
		query += " AND status = ?"
		args = append(args, filter.Status)
	}
	if filter.MinPnL != nil {
		query += " AND pnl >= ?"
		args = append(args, *filter.MinPnL)
	}
	if filter.MaxPnL != nil {
		query += " AND pnl <= ?"
		args = append(args, *filter.MaxPnL)
	}
	if filter.FromDate != nil {
		query += " AND created_at >= ?"
		args = append(args, *filter.FromDate)
	}
	if filter.ToDate != nil {
		query += " AND created_at <= ?"
		args = append(args, *filter.ToDate)
	}

	// Execute query
	rows := []positionRow{}
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("failed to query positions: %w", err)
	}

	// Convert to models
	positions := make([]*models.Position, 0, len(rows))
	for _, row := range rows {
		position, err := row.toModel()
		if err != nil {
			return nil, err
		}
		positions = append(positions, position)
	}

	return positions, nil
}

// FindByID returns a specific position by ID
func (r *SQLitePositionRepository) FindByID(ctx context.Context, id string) (*models.Position, error) {
	query := "SELECT * FROM positions WHERE id = ?"
	row := positionRow{}
	if err := r.db.GetContext(ctx, &row, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("position not found: %s", id)
		}
		return nil, fmt.Errorf("failed to query position: %w", err)
	}

	return row.toModel()
}

// FindBySymbol returns positions for a specific symbol
func (r *SQLitePositionRepository) FindBySymbol(ctx context.Context, symbol string) ([]*models.Position, error) {
	query := "SELECT * FROM positions WHERE symbol = ?"
	rows := []positionRow{}
	if err := r.db.SelectContext(ctx, &rows, query, symbol); err != nil {
		return nil, fmt.Errorf("failed to query positions: %w", err)
	}

	// Convert to models
	positions := make([]*models.Position, 0, len(rows))
	for _, row := range rows {
		position, err := row.toModel()
		if err != nil {
			return nil, err
		}
		positions = append(positions, position)
	}

	return positions, nil
}

// Create adds a new position
func (r *SQLitePositionRepository) Create(ctx context.Context, position *models.Position) (string, error) {
	// Generate ID if not provided
	if position.ID == "" {
		position.ID = fmt.Sprintf("pos_%d", time.Now().UnixNano())
	}

	// Convert to row
	row, err := fromModel(position)
	if err != nil {
		return "", err
	}

	// Insert into database
	query := `
		INSERT INTO positions (
			id, symbol, quantity, entry_price, current_price, open_time,
			stop_loss, take_profit, trailing_stop, created_at, updated_at,
			pnl, pnl_percentage, status, orders_json
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = r.db.ExecContext(ctx, query,
		row.ID, row.Symbol, row.Quantity, row.EntryPrice, row.CurrentPrice, row.OpenTime,
		row.StopLoss, row.TakeProfit, row.TrailingStop, row.CreatedAt, row.UpdatedAt,
		row.PnL, row.PnLPercentage, row.Status, row.OrdersJSON,
	)
	if err != nil {
		return "", fmt.Errorf("failed to insert position: %w", err)
	}

	return position.ID, nil
}

// Update modifies an existing position
func (r *SQLitePositionRepository) Update(ctx context.Context, position *models.Position) error {
	// Convert to row
	row, err := fromModel(position)
	if err != nil {
		return err
	}

	// Update in database
	query := `
		UPDATE positions SET
			symbol = ?, quantity = ?, entry_price = ?, current_price = ?, open_time = ?,
			stop_loss = ?, take_profit = ?, trailing_stop = ?, updated_at = ?,
			pnl = ?, pnl_percentage = ?, status = ?, orders_json = ?
		WHERE id = ?
	`
	_, err = r.db.ExecContext(ctx, query,
		row.Symbol, row.Quantity, row.EntryPrice, row.CurrentPrice, row.OpenTime,
		row.StopLoss, row.TakeProfit, row.TrailingStop, row.UpdatedAt,
		row.PnL, row.PnLPercentage, row.Status, row.OrdersJSON,
		row.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update position: %w", err)
	}

	return nil
}

// Delete removes a position
func (r *SQLitePositionRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM positions WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete position: %w", err)
	}

	return nil
}

// AddOrder adds an order to a position
func (r *SQLitePositionRepository) AddOrder(ctx context.Context, positionID string, order *models.Order) error {
	// Get the position
	position, err := r.FindByID(ctx, positionID)
	if err != nil {
		return err
	}

	// Add the order
	position.Orders = append(position.Orders, *order)

	// Update the position
	return r.Update(ctx, position)
}

// UpdateOrder updates an order in a position
func (r *SQLitePositionRepository) UpdateOrder(ctx context.Context, positionID string, order *models.Order) error {
	// Get the position
	position, err := r.FindByID(ctx, positionID)
	if err != nil {
		return err
	}

	// Find and update the order
	found := false
	for i, o := range position.Orders {
		if o.ID == order.ID {
			position.Orders[i] = *order
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("order not found in position: %s", order.ID)
	}

	// Update the position
	return r.Update(ctx, position)
}

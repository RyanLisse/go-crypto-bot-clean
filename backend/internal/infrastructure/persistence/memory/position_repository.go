package memory

import (
	"context"
	"errors"
	"sync"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/ports"
)

// PositionRepository implements the ports.PositionRepository interface with in-memory storage
type PositionRepository struct {
	mu        sync.RWMutex
	positions map[string]*models.Position
}

// NewPositionRepository creates a new in-memory position repository
func NewPositionRepository() ports.PositionRepository {
	return &PositionRepository{
		positions: make(map[string]*models.Position),
	}
}

// Create stores a new position in memory
func (r *PositionRepository) Create(ctx context.Context, position *models.Position) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if position.ID == "" {
		return errors.New("position ID cannot be empty")
	}

	r.positions[position.ID] = position
	return nil
}

// GetByID retrieves a position by its ID
func (r *PositionRepository) GetByID(ctx context.Context, id string) (*models.Position, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	position, exists := r.positions[id]
	if !exists {
		return nil, errors.New("position not found")
	}

	return position, nil
}

// List retrieves positions based on status
func (r *PositionRepository) List(ctx context.Context, status models.PositionStatus) ([]*models.Position, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*models.Position
	for _, position := range r.positions {
		if status == "" || position.Status == status {
			result = append(result, position)
		}
	}

	return result, nil
}

// Update updates an existing position
func (r *PositionRepository) Update(ctx context.Context, position *models.Position) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if position.ID == "" {
		return errors.New("position ID cannot be empty")
	}

	if _, exists := r.positions[position.ID]; !exists {
		return errors.New("position not found")
	}

	r.positions[position.ID] = position
	return nil
}

// Delete removes a position
func (r *PositionRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.positions[id]; !exists {
		return errors.New("position not found")
	}

	delete(r.positions, id)
	return nil
}

// GetOpenPositionBySymbol retrieves an open position for a specific symbol
func (r *PositionRepository) GetOpenPositionBySymbol(ctx context.Context, symbol string) (*models.Position, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, position := range r.positions {
		if position.Symbol == symbol && position.Status == models.PositionStatusOpen {
			return position, nil
		}
	}

	return nil, errors.New("no open position found for symbol")
}

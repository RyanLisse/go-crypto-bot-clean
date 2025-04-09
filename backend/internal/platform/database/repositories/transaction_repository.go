package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// SQLiteTransactionRepository implements the TransactionRepository interface for SQLite
type SQLiteTransactionRepository struct {
	db *sqlx.DB
}

// NewSQLiteTransactionRepository creates a new SQLite transaction repository
func NewSQLiteTransactionRepository(db *sqlx.DB) repository.TransactionRepository {
	return &SQLiteTransactionRepository{
		db: db,
	}
}

// Create creates a new transaction record
func (r *SQLiteTransactionRepository) Create(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error) {
	// Generate a new UUID for the transaction
	transaction.ID = uuid.New().String()

	query := `
		INSERT INTO transactions (
			id, amount, balance, reason, timestamp,
			created_at, updated_at, wallet_id, position_id, order_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	transaction.CreatedAt = now
	transaction.UpdatedAt = now
	
	_, err := r.db.ExecContext(ctx, query,
		transaction.ID,
		transaction.Amount,
		transaction.Balance,
		transaction.Reason,
		transaction.Timestamp,
		transaction.CreatedAt,
		transaction.UpdatedAt,
		transaction.WalletID,
		transaction.PositionID,
		transaction.OrderID)
	
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

// FindByID retrieves a transaction by its ID
func (r *SQLiteTransactionRepository) FindByID(ctx context.Context, id string) (*models.Transaction, error) {
	query := `
		SELECT id, amount, balance, reason, timestamp,
			created_at, updated_at, wallet_id, position_id, order_id
		FROM transactions
		WHERE id = ?
	`
	
	var transaction models.Transaction
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.Amount,
		&transaction.Balance,
		&transaction.Reason,
		&transaction.Timestamp,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&transaction.WalletID,
		&transaction.PositionID,
		&transaction.OrderID)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	
	return &transaction, nil
}

// FindByTimeRange retrieves transactions within a time range
func (r *SQLiteTransactionRepository) FindByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*models.Transaction, error) {
	query := `
		SELECT id, amount, balance, reason, timestamp,
			created_at, updated_at, wallet_id, position_id, order_id
		FROM transactions
		WHERE timestamp BETWEEN ? AND ?
		ORDER BY timestamp DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var transactions []*models.Transaction
	for rows.Next() {
		var tx models.Transaction
		if err := rows.Scan(
			&tx.ID,
			&tx.Amount,
			&tx.Balance,
			&tx.Reason,
			&tx.Timestamp,
			&tx.CreatedAt,
			&tx.UpdatedAt,
			&tx.WalletID,
			&tx.PositionID,
			&tx.OrderID,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, &tx)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return transactions, nil
}

// FindAll retrieves all transactions
func (r *SQLiteTransactionRepository) FindAll(ctx context.Context) ([]*models.Transaction, error) {
	query := `
		SELECT id, amount, balance, reason, timestamp,
			created_at, updated_at, wallet_id, position_id, order_id
		FROM transactions
		ORDER BY timestamp DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var transactions []*models.Transaction
	for rows.Next() {
		var tx models.Transaction
		if err := rows.Scan(
			&tx.ID,
			&tx.Amount,
			&tx.Balance,
			&tx.Reason,
			&tx.Timestamp,
			&tx.CreatedAt,
			&tx.UpdatedAt,
			&tx.WalletID,
			&tx.PositionID,
			&tx.OrderID,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, &tx)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return transactions, nil
}

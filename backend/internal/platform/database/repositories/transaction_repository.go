package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/repository"
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
	query := `
		INSERT INTO transactions (amount, balance, reason, timestamp)
		VALUES (?, ?, ?, ?)
		RETURNING id
	`
	
	var id int64
	err := r.db.QueryRowContext(ctx, query, 
		transaction.Amount, 
		transaction.Balance, 
		transaction.Reason, 
		transaction.Timestamp).Scan(&id)
	
	if err != nil {
		return nil, err
	}
	
	transaction.ID = id
	return transaction, nil
}

// FindByID retrieves a transaction by its ID
func (r *SQLiteTransactionRepository) FindByID(ctx context.Context, id int64) (*models.Transaction, error) {
	query := `
		SELECT id, amount, balance, reason, timestamp
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
	)
	
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
		SELECT id, amount, balance, reason, timestamp
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
		SELECT id, amount, balance, reason, timestamp
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

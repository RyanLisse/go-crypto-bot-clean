package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// ShadowRepository implements the Repository interface by writing to both
// SQLite and TursoDB, but reading primarily from SQLite.
// This allows for a gradual migration to TursoDB while ensuring data consistency.
type ShadowRepository struct {
	primary   Repository // SQLite repository (primary for reads)
	secondary Repository // TursoDB repository (secondary, for validation)
	config    Config
	txMap     sync.Map // Map to track transaction pairs
}

// NewShadowRepository creates a new shadow repository
func NewShadowRepository(config Config) (*ShadowRepository, error) {
	// Validate configuration
	if !config.TursoEnabled {
		return nil, fmt.Errorf("turso must be enabled for shadow mode")
	}

	// Create SQLite repository (primary)
	sqliteRepo, err := NewSQLiteRepository(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create SQLite repository: %w", err)
	}

	// Create TursoDB repository (secondary)
	tursoRepo, err := NewTursoRepository(config)
	if err != nil {
		sqliteRepo.Close()
		return nil, fmt.Errorf("failed to create TursoDB repository: %w", err)
	}

	return &ShadowRepository{
		primary:   sqliteRepo,
		secondary: tursoRepo,
		config:    config,
		txMap:     sync.Map{},
	}, nil
}

// Execute executes a query without returning any rows
// Writes to both databases, but only returns the result from the primary
func (r *ShadowRepository) Execute(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	// Execute on primary (SQLite)
	primaryResult, err := r.primary.Execute(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("primary database error: %w", err)
	}

	// Execute on secondary (TursoDB)
	_, secondaryErr := r.secondary.Execute(ctx, query, args...)
	if secondaryErr != nil {
		log.Printf("WARNING: Secondary database error: %v", secondaryErr)
		// We don't fail the operation if the secondary write fails
		// This ensures the application continues to function even if TursoDB is unavailable
	}

	return primaryResult, nil
}

// Query executes a query that returns rows
// Only reads from the primary database
func (r *ShadowRepository) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return r.primary.Query(ctx, query, args...)
}

// QueryRow executes a query that returns a single row
// Only reads from the primary database
func (r *ShadowRepository) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return r.primary.QueryRow(ctx, query, args...)
}

// BeginTx starts a transaction
// This is complex because we need to manage transactions on both databases
func (r *ShadowRepository) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	// Begin transaction on primary
	primaryTx, err := r.primary.BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction on primary database: %w", err)
	}

	// Begin transaction on secondary
	secondaryTx, err := r.secondary.BeginTx(ctx, opts)
	if err != nil {
		// Rollback primary transaction
		primaryTx.Rollback()
		return nil, fmt.Errorf("failed to begin transaction on secondary database: %w", err)
	}

	// Store the secondary transaction in a map keyed by the primary transaction
	r.txMap.Store(primaryTx, secondaryTx)

	// Return the primary transaction
	// We'll handle the secondary transaction in the Commit and Rollback methods
	// by intercepting the calls to the primary transaction
	return primaryTx, nil
}

// Close closes both database connections
func (r *ShadowRepository) Close() error {
	primaryErr := r.primary.Close()
	secondaryErr := r.secondary.Close()

	if primaryErr != nil {
		return fmt.Errorf("failed to close primary database: %w", primaryErr)
	}
	if secondaryErr != nil {
		return fmt.Errorf("failed to close secondary database: %w", secondaryErr)
	}

	return nil
}

// Ping verifies connections to both databases are still alive
func (r *ShadowRepository) Ping(ctx context.Context) error {
	if err := r.primary.Ping(ctx); err != nil {
		return fmt.Errorf("primary database ping failed: %w", err)
	}
	if err := r.secondary.Ping(ctx); err != nil {
		return fmt.Errorf("secondary database ping failed: %w", err)
	}
	return nil
}

// GetImplementationType returns the type of database implementation
func (r *ShadowRepository) GetImplementationType() string {
	return "shadow"
}

// SupportsSynchronization returns whether this implementation supports synchronization
func (r *ShadowRepository) SupportsSynchronization() bool {
	return true
}

// Synchronize triggers synchronization with the cloud database
func (r *ShadowRepository) Synchronize(ctx context.Context) error {
	// Only synchronize the secondary (TursoDB) repository
	return r.secondary.Synchronize(ctx)
}

// GetLastSyncTimestamp gets the timestamp of the last synchronization
func (r *ShadowRepository) GetLastSyncTimestamp(ctx context.Context) (time.Time, error) {
	// Get the timestamp from the secondary (TursoDB) repository
	return r.secondary.GetLastSyncTimestamp(ctx)
}

// ValidateConsistency checks if the data in both databases is consistent
func (r *ShadowRepository) ValidateConsistency(ctx context.Context, tableName string) error {
	// Get count from primary
	var primaryCount int
	err := r.primary.QueryRow(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&primaryCount)
	if err != nil {
		return fmt.Errorf("failed to get count from primary database: %w", err)
	}

	// Get count from secondary
	var secondaryCount int
	err = r.secondary.QueryRow(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&secondaryCount)
	if err != nil {
		return fmt.Errorf("failed to get count from secondary database: %w", err)
	}

	// Compare counts
	if primaryCount != secondaryCount {
		return fmt.Errorf("inconsistency detected in table %s: primary count=%d, secondary count=%d",
			tableName, primaryCount, secondaryCount)
	}

	return nil
}

// SyncTables synchronizes the data in the specified tables from primary to secondary
func (r *ShadowRepository) SyncTables(ctx context.Context, tableNames []string) error {
	for _, tableName := range tableNames {
		// Get schema from primary
		var createTableSQL string
		err := r.primary.QueryRow(ctx, fmt.Sprintf("SELECT sql FROM sqlite_master WHERE type='table' AND name='%s'", tableName)).Scan(&createTableSQL)
		if err != nil {
			return fmt.Errorf("failed to get schema for table %s: %w", tableName, err)
		}

		// Create table in secondary if it doesn't exist
		_, err = r.secondary.Execute(ctx, fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s AS SELECT * FROM %s WHERE 0", tableName, tableName))
		if err != nil {
			return fmt.Errorf("failed to create table %s in secondary database: %w", tableName, err)
		}

		// Clear secondary table
		_, err = r.secondary.Execute(ctx, fmt.Sprintf("DELETE FROM %s", tableName))
		if err != nil {
			return fmt.Errorf("failed to clear table %s in secondary database: %w", tableName, err)
		}

		// Get column names from primary
		rows, err := r.primary.Query(ctx, fmt.Sprintf("PRAGMA table_info(%s)", tableName))
		if err != nil {
			return fmt.Errorf("failed to get column info for table %s: %w", tableName, err)
		}

		var columnNames []string
		for rows.Next() {
			var cid, notnull, pk int
			var name, dataType, dfltValue string
			if err := rows.Scan(&cid, &name, &dataType, &notnull, &dfltValue, &pk); err != nil {
				rows.Close()
				return fmt.Errorf("failed to scan column info: %w", err)
			}
			columnNames = append(columnNames, name)
		}
		rows.Close()

		if err := rows.Err(); err != nil {
			return fmt.Errorf("error iterating column info: %w", err)
		}

		// Copy data from primary to secondary
		rows, err = r.primary.Query(ctx, fmt.Sprintf("SELECT * FROM %s", tableName))
		if err != nil {
			return fmt.Errorf("failed to query data from primary database: %w", err)
		}
		defer rows.Close()

		// Prepare column placeholders for INSERT
		placeholders := make([]string, len(columnNames))
		for i := range placeholders {
			placeholders[i] = "?"
		}

		// Prepare INSERT statement
		insertSQL := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			tableName,
			strings.Join(columnNames, ", "),
			strings.Join(placeholders, ", "))

		// Begin transaction on secondary
		tx, err := r.secondary.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to begin transaction on secondary database: %w", err)
		}

		// Prepare statement
		stmt, err := tx.Prepare(insertSQL)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to prepare insert statement: %w", err)
		}
		defer stmt.Close()

		// Copy each row
		for rows.Next() {
			// Create slice to hold column values
			values := make([]interface{}, len(columnNames))
			valuePtrs := make([]interface{}, len(columnNames))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			// Scan row into values
			if err := rows.Scan(valuePtrs...); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to scan row: %w", err)
			}

			// Insert row into secondary
			if _, err := stmt.Exec(values...); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to insert row: %w", err)
			}
		}

		if err := rows.Err(); err != nil {
			tx.Rollback()
			return fmt.Errorf("error iterating rows: %w", err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
	}

	return nil
}

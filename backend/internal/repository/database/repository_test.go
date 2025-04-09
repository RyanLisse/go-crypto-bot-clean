package database

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSQLiteRepository tests the SQLite implementation of the Repository interface
func TestSQLiteRepository(t *testing.T) {
	// Create a temporary database file
	dbPath := "./test_sqlite.db"
	defer os.Remove(dbPath)

	// Create configuration
	config := Config{
		DatabasePath:    dbPath,
		MaxOpenConns:    5,
		MaxIdleConns:    2,
		ConnMaxLifetime: 1 * time.Minute,
	}

	// Create repository
	repo, err := NewSQLiteRepository(config)
	require.NoError(t, err)
	defer repo.Close()

	// Run common tests
	runRepositoryTests(t, repo)
}

// TestTursoRepository tests the TursoDB implementation of the Repository interface
// This test is skipped by default because it requires a TursoDB instance
func TestTursoRepository(t *testing.T) {
	// Skip if TursoDB URL or auth token is not set
	tursoURL := os.Getenv("TURSO_URL")
	tursoAuthToken := os.Getenv("TURSO_AUTH_TOKEN")
	if tursoURL == "" || tursoAuthToken == "" {
		t.Skip("TURSO_URL or TURSO_AUTH_TOKEN not set, skipping TursoDB tests")
	}

	// Create configuration for local database with sync
	dbPath := "./test_turso.db"
	defer os.Remove(dbPath)

	config := Config{
		DatabasePath:    dbPath,
		MaxOpenConns:    5,
		MaxIdleConns:    2,
		ConnMaxLifetime: 1 * time.Minute,
		TursoEnabled:    true,
		TursoURL:        tursoURL,
		TursoAuthToken:  tursoAuthToken,
		SyncEnabled:     true,
		SyncInterval:    1 * time.Minute,
	}

	// Create repository
	repo, err := NewTursoRepository(config)
	require.NoError(t, err)
	defer repo.Close()

	// Run common tests
	runRepositoryTests(t, repo)

	// Run TursoDB-specific tests
	runTursoSpecificTests(t, repo)
}

// runRepositoryTests runs common tests for any Repository implementation
func runRepositoryTests(t *testing.T, repo Repository) {
	ctx := context.Background()

	// Test basic operations
	testBasicOperations(t, ctx, repo)

	// Test transactions
	testTransactions(t, ctx, repo)
}

// testBasicOperations tests basic CRUD operations
func testBasicOperations(t *testing.T, ctx context.Context, repo Repository) {
	// Create test table
	_, err := repo.Execute(ctx, `
		CREATE TABLE IF NOT EXISTS test_table (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			value INTEGER NOT NULL
		)
	`)
	require.NoError(t, err)

	// Insert data
	result, err := repo.Execute(ctx, `
		INSERT INTO test_table (name, value) VALUES (?, ?)
	`, "test1", 42)
	require.NoError(t, err)

	// Check rows affected
	rowsAffected, err := result.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	// Query data
	rows, err := repo.Query(ctx, `
		SELECT id, name, value FROM test_table WHERE name = ?
	`, "test1")
	require.NoError(t, err)
	defer rows.Close()

	// Check query results
	assert.True(t, rows.Next())
	var id int
	var name string
	var value int
	err = rows.Scan(&id, &name, &value)
	require.NoError(t, err)
	assert.Equal(t, "test1", name)
	assert.Equal(t, 42, value)
	assert.False(t, rows.Next())

	// QueryRow
	var count int
	err = repo.QueryRow(ctx, `
		SELECT COUNT(*) FROM test_table
	`).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	// Update data
	result, err = repo.Execute(ctx, `
		UPDATE test_table SET value = ? WHERE name = ?
	`, 84, "test1")
	require.NoError(t, err)
	rowsAffected, err = result.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	// Verify update
	err = repo.QueryRow(ctx, `
		SELECT value FROM test_table WHERE name = ?
	`, "test1").Scan(&value)
	require.NoError(t, err)
	assert.Equal(t, 84, value)

	// Delete data
	result, err = repo.Execute(ctx, `
		DELETE FROM test_table WHERE name = ?
	`, "test1")
	require.NoError(t, err)
	rowsAffected, err = result.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	// Verify delete
	err = repo.QueryRow(ctx, `
		SELECT COUNT(*) FROM test_table
	`).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

// testTransactions tests transaction support
func testTransactions(t *testing.T, ctx context.Context, repo Repository) {
	// Create test table
	_, err := repo.Execute(ctx, `
		CREATE TABLE IF NOT EXISTS test_transactions (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL
		)
	`)
	require.NoError(t, err)

	// Test successful transaction
	tx, err := repo.BeginTx(ctx, nil)
	require.NoError(t, err)

	_, err = tx.Exec("INSERT INTO test_transactions (name) VALUES (?)", "tx1")
	require.NoError(t, err)

	_, err = tx.Exec("INSERT INTO test_transactions (name) VALUES (?)", "tx2")
	require.NoError(t, err)

	err = tx.Commit()
	require.NoError(t, err)

	// Verify transaction committed
	var count int
	err = repo.QueryRow(ctx, `
		SELECT COUNT(*) FROM test_transactions
	`).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 2, count)

	// Test rollback
	tx, err = repo.BeginTx(ctx, nil)
	require.NoError(t, err)

	_, err = tx.Exec("INSERT INTO test_transactions (name) VALUES (?)", "tx3")
	require.NoError(t, err)

	err = tx.Rollback()
	require.NoError(t, err)

	// Verify transaction rolled back
	err = repo.QueryRow(ctx, `
		SELECT COUNT(*) FROM test_transactions
	`).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 2, count)

	// Test WithTransaction helper
	err = WithTransaction(ctx, repo, func(tx *sql.Tx) error {
		// Execute a simple SQL statement
		_, err := tx.Exec("INSERT INTO test_transactions (name) VALUES (?)", "tx3")
		return err
	})
	require.NoError(t, err)

	// Verify helper transaction committed
	err = repo.QueryRow(ctx, `
		SELECT COUNT(*) FROM test_transactions
	`).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 3, count)

	// Test WithTransaction rollback on error
	err = WithTransaction(ctx, repo, func(tx *sql.Tx) error {
		// Execute a simple SQL statement
		_, err := tx.Exec("INSERT INTO test_transactions (name) VALUES (?)", "tx4")
		if err != nil {
			return err
		}
		return errors.New("intentional error")
	})
	require.Error(t, err)

	// Verify helper transaction rolled back
	err = repo.QueryRow(ctx, `
		SELECT COUNT(*) FROM test_transactions
	`).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 3, count)
}

// runTursoSpecificTests runs tests specific to TursoDB
func runTursoSpecificTests(t *testing.T, repo Repository) {
	// Skip if repository doesn't support synchronization
	if !repo.SupportsSynchronization() {
		t.Skip("Repository does not support synchronization")
	}

	ctx := context.Background()

	// Test synchronization
	err := repo.Synchronize(ctx)
	require.NoError(t, err)

	// Test getting last sync timestamp
	timestamp, err := repo.GetLastSyncTimestamp(ctx)
	require.NoError(t, err)
	assert.False(t, timestamp.IsZero())

	// Verify implementation type
	assert.Equal(t, "turso", repo.GetImplementationType())
}

package database

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShadowRepository(t *testing.T) {
	// Skip this test if we're not in a development environment
	// or if the TursoDB URL is not set
	if os.Getenv("TURSO_URL") == "" {
		t.Skip("Skipping TursoDB test because TURSO_URL is not set")
	}

	// Create temporary SQLite database
	sqliteDbPath := "test_shadow_sqlite.db"
	defer os.Remove(sqliteDbPath)

	// Create configuration
	config := Config{
		DatabasePath:    sqliteDbPath,
		MaxOpenConns:    5,
		MaxIdleConns:    2,
		ConnMaxLifetime: 5 * time.Minute,
		TursoEnabled:    true,
		TursoURL:        os.Getenv("TURSO_URL"),
		TursoAuthToken:  os.Getenv("TURSO_AUTH_TOKEN"),
		SyncEnabled:     true,
		SyncInterval:    5 * time.Minute,
		ShadowMode:      true,
	}

	// Create shadow repository
	repo, err := NewShadowRepository(config)
	require.NoError(t, err)
	defer repo.Close()

	// Test basic operations
	t.Run("Execute", func(t *testing.T) {
		ctx := context.Background()

		// Create test table
		_, err := repo.Execute(ctx, `
			CREATE TABLE IF NOT EXISTS test_shadow (
				id INTEGER PRIMARY KEY,
				name TEXT NOT NULL,
				value INTEGER NOT NULL
			)
		`)
		require.NoError(t, err)

		// Insert data
		result, err := repo.Execute(ctx, "INSERT INTO test_shadow (name, value) VALUES (?, ?)", "test1", 42)
		require.NoError(t, err)

		// Check result
		id, err := result.LastInsertId()
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		// Verify data was written to both databases
		var sqliteCount, tursoCount int
		err = repo.primary.QueryRow(ctx, "SELECT COUNT(*) FROM test_shadow WHERE name = ?", "test1").Scan(&sqliteCount)
		require.NoError(t, err)
		err = repo.secondary.QueryRow(ctx, "SELECT COUNT(*) FROM test_shadow WHERE name = ?", "test1").Scan(&tursoCount)
		require.NoError(t, err)

		assert.Equal(t, 1, sqliteCount)
		assert.Equal(t, 1, tursoCount)
	})

	t.Run("Query", func(t *testing.T) {
		ctx := context.Background()

		// Query data
		rows, err := repo.Query(ctx, "SELECT id, name, value FROM test_shadow WHERE name = ?", "test1")
		require.NoError(t, err)
		defer rows.Close()

		// Check result
		assert.True(t, rows.Next())
		var id int
		var name string
		var value int
		err = rows.Scan(&id, &name, &value)
		require.NoError(t, err)
		assert.Equal(t, "test1", name)
		assert.Equal(t, 42, value)
		assert.False(t, rows.Next())
	})

	t.Run("QueryRow", func(t *testing.T) {
		ctx := context.Background()

		// Query single row
		var id int
		var name string
		var value int
		err := repo.QueryRow(ctx, "SELECT id, name, value FROM test_shadow WHERE name = ?", "test1").Scan(&id, &name, &value)
		require.NoError(t, err)
		assert.Equal(t, "test1", name)
		assert.Equal(t, 42, value)
	})

	t.Run("Transaction", func(t *testing.T) {
		ctx := context.Background()

		// Begin transaction
		tx, err := repo.BeginTx(ctx, nil)
		require.NoError(t, err)

		// Insert data in transaction
		_, err = tx.Exec("INSERT INTO test_shadow (name, value) VALUES (?, ?)", "test2", 84)
		require.NoError(t, err)

		// Commit transaction
		err = tx.Commit()
		require.NoError(t, err)

		// Verify data was written to both databases
		var sqliteCount, tursoCount int
		err = repo.primary.QueryRow(ctx, "SELECT COUNT(*) FROM test_shadow WHERE name = ?", "test2").Scan(&sqliteCount)
		require.NoError(t, err)
		err = repo.secondary.QueryRow(ctx, "SELECT COUNT(*) FROM test_shadow WHERE name = ?", "test2").Scan(&tursoCount)
		require.NoError(t, err)

		assert.Equal(t, 1, sqliteCount)
		assert.Equal(t, 1, tursoCount)
	})

	t.Run("Validation", func(t *testing.T) {
		ctx := context.Background()

		// Run validation
		err := repo.ValidateConsistency(ctx, "test_shadow")
		require.NoError(t, err)

		// Insert data directly to primary (SQLite) to create inconsistency
		_, err = repo.primary.Execute(ctx, "INSERT INTO test_shadow (name, value) VALUES (?, ?)", "inconsistent", 100)
		require.NoError(t, err)

		// Validation should fail
		err = repo.ValidateConsistency(ctx, "test_shadow")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "inconsistency detected")

		// Fix inconsistency by syncing
		err = repo.SyncTables(ctx, []string{"test_shadow"})
		require.NoError(t, err)

		// Validation should pass now
		err = repo.ValidateConsistency(ctx, "test_shadow")
		require.NoError(t, err)
	})
}

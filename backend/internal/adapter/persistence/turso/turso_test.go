package turso

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/tursodatabase/go-libsql"
)

func TestTursoIntegration(t *testing.T) {
	// Initialize logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Load environment variables
	err := godotenv.Load()
	require.NoError(t, err, "Failed to load .env file")

	// Get Turso database configuration
	primaryUrl := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	if primaryUrl == "" || authToken == "" {
		t.Skip("Skipping test: TURSO_DATABASE_URL or TURSO_AUTH_TOKEN not set")
	}

	// Create a temporary directory for the embedded database
	dir, err := os.MkdirTemp("", "libsql-test-*")
	if err != nil {
		t.Fatal(apperror.NewExternalService("Turso", "Failed to create temporary directory", err))
	}
	defer os.RemoveAll(dir)

	// Define local database path
	dbPath := filepath.Join(dir, "local.db")

	// Create embedded replica connector
	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl,
		libsql.WithAuthToken(authToken),
	)
	if err != nil {
		t.Fatal(apperror.NewExternalService("Turso", "Failed to create connector", err))
	}
	defer connector.Close()

	db := sql.OpenDB(connector)
	defer db.Close()

	// Test database connection
	err = db.Ping()
	if err != nil {
		t.Fatal(apperror.NewExternalService("Turso", "Failed to ping database", err))
	}

	// Create test table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS test_table (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatal(apperror.NewExternalService("Turso", "Failed to create test table", err))
	}

	// Insert a record
	result, err := db.Exec("INSERT INTO test_table (name) VALUES (?)", "test_name")
	if err != nil {
		t.Fatal(apperror.NewExternalService("Turso", "Failed to insert record", err))
	}

	// Get the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(apperror.NewExternalService("Turso", "Failed to get last insert ID", err))
	}

	// Query the record
	var name string
	var createdAt time.Time
	err = db.QueryRow("SELECT name, created_at FROM test_table WHERE id = ?", id).Scan(&name, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			t.Fatal(apperror.NewNotFound("Record", id, err))
		}
		t.Fatal(apperror.NewExternalService("Turso", "Failed to query record", err))
	}

	// Verify the record
	require.Equal(t, "test_name", name)
	require.WithinDuration(t, time.Now(), createdAt, 5*time.Second)

	// Manually trigger sync with primary database
	// The Sync method returns a Replicated instance and error
	_, err = connector.Sync()
	if err != nil {
		t.Fatal(apperror.NewExternalService("Turso", "Failed to sync with primary database", err))
	}

	// Clean up - drop the test table
	_, err = db.Exec("DROP TABLE test_table")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to drop test table")
		// Don't fail the test on cleanup error
	}
}

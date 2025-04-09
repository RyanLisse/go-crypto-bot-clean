# SQLite Migrations for Go Crypto Bot

This guide covers a simple, dependency-free approach to managing SQLite database migrations for your Go crypto trading bot.

## 1. Migration Structure

We'll use a simple embedded approach with plain SQL files and a version tracking table.

### Directory Structure

```
internal/platform/database/
├── migrations/
│   ├── 001_initial_schema.sql
│   ├── 002_add_take_profit_levels.sql
│   └── 003_add_indexes.sql
├── migrate.go
└── sqlite.go
```

## 2. Creating the Migration System

Let's create a simple migration system:

```go
// internal/platform/database/migrate.go
package database

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Migrate applies all pending migrations to the database
func Migrate(db *sql.DB) error {
	// Create migrations table if it doesn't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("creating migrations table: %w", err)
	}

	// Get already applied migrations
	appliedMigrations := make(map[string]bool)
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return fmt.Errorf("querying applied migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return fmt.Errorf("scanning migration version: %w", err)
		}
		appliedMigrations[version] = true
	}

	// Get all migration files
	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("reading migrations directory: %w", err)
	}

	// Sort migration files by name
	var migrationFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			migrationFiles = append(migrationFiles, entry.Name())
		}
	}
	sort.Strings(migrationFiles)

	// Apply pending migrations
	for _, filename := range migrationFiles {
		// Extract version from filename (e.g., "001_initial_schema.sql" -> "001")
		version := strings.Split(filename, "_")[0]
		
		// Skip already applied migrations
		if appliedMigrations[version] {
			log.Printf("Migration %s already applied, skipping", version)
			continue
		}

		log.Printf("Applying migration %s...", filename)

		// Read migration file
		content, err := fs.ReadFile(migrationsFS, "migrations/"+filename)
		if err != nil {
			return fmt.Errorf("reading migration file %s: %w", filename, err)
		}

		// Start a transaction for atomicity
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("starting transaction for migration %s: %w", filename, err)
		}

		// Apply migration
		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("executing migration %s: %w", filename, err)
		}

		// Record migration as applied
		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version); err != nil {
			tx.Rollback()
			return fmt.Errorf("recording migration %s: %w", filename, err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("committing migration %s: %w", filename, err)
		}

		log.Printf("Migration %s applied successfully", filename)
	}

	return nil
}
```

## 3. Creating Migration Files

Create SQL migration files in the `migrations` directory. Here's an example of the initial schema:

```sql
-- migrations/001_initial_schema.sql
-- Initial database schema for crypto trading bot

-- Store bought coins that we're tracking
CREATE TABLE IF NOT EXISTS bought_coins (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    symbol TEXT NOT NULL UNIQUE,
    purchase_price REAL NOT NULL,
    quantity REAL NOT NULL,
    purchase_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted INTEGER DEFAULT 0,
    stop_loss_price REAL NOT NULL
);

-- Take profit levels for each bought coin
CREATE TABLE IF NOT EXISTS take_profit_levels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    bought_coin_id INTEGER NOT NULL,
    percentage REAL NOT NULL,
    sell_quantity REAL NOT NULL,
    is_reached INTEGER DEFAULT 0,
    FOREIGN KEY (bought_coin_id) REFERENCES bought_coins(id) ON DELETE CASCADE
);

-- Store newly detected coins
CREATE TABLE IF NOT EXISTS new_coins (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    symbol TEXT NOT NULL UNIQUE,
    detected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_checked TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active INTEGER DEFAULT 1
);

-- Track purchase decisions (both buys and rejects)
CREATE TABLE IF NOT EXISTS purchase_decisions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    symbol TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status TEXT NOT NULL,
    reason TEXT,
    price REAL
);

-- Store log events
CREATE TABLE IF NOT EXISTS log_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    level TEXT NOT NULL,
    message TEXT NOT NULL,
    context TEXT
);
```

Here's a second migration that adds indexes for performance:

```sql
-- migrations/002_add_indexes.sql
-- Add indexes for better query performance

-- Indexes for bought_coins
CREATE INDEX IF NOT EXISTS idx_bought_coins_symbol ON bought_coins(symbol);
CREATE INDEX IF NOT EXISTS idx_bought_coins_is_deleted ON bought_coins(is_deleted);

-- Indexes for take_profit_levels
CREATE INDEX IF NOT EXISTS idx_tp_levels_bought_coin_id ON take_profit_levels(bought_coin_id);

-- Indexes for new_coins
CREATE INDEX IF NOT EXISTS idx_new_coins_symbol ON new_coins(symbol);
CREATE INDEX IF NOT EXISTS idx_new_coins_is_active ON new_coins(is_active);

-- Indexes for purchase_decisions
CREATE INDEX IF NOT EXISTS idx_purchase_decisions_symbol ON purchase_decisions(symbol);
CREATE INDEX IF NOT EXISTS idx_purchase_decisions_status ON purchase_decisions(status);

-- Indexes for log_events
CREATE INDEX IF NOT EXISTS idx_log_events_timestamp ON log_events(timestamp);
CREATE INDEX IF NOT EXISTS idx_log_events_level ON log_events(level);
```

## 4. Migration Best Practices

### Idempotent Migrations

Always make your migrations idempotent (can be run multiple times safely) by using `IF NOT EXISTS` or similar guards.

### Rollback Support

For more sophisticated needs, consider adding rollback support:

```sql
-- with rollback section
-- @up
CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT);

-- @down
DROP TABLE users;
```

Then modify the migration system to parse and handle these sections.

### Testing Migrations

Create a test that applies all migrations to an in-memory database:

```go
// internal/platform/database/migrate_test.go
package database_test

import (
	"testing"

	"github.com/ryanlisse/cryptobot/internal/platform/database"
)

func TestMigrate(t *testing.T) {
	// Create in-memory database
	db, err := database.Connect(database.Config{InMemory: true})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer database.Close(db)

	// Apply migrations
	if err := database.Migrate(db); err != nil {
		t.Fatalf("Failed to apply migrations: %v", err)
	}

	// Verify that migrations were applied
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query migration count: %v", err)
	}

	// Check if we have the expected number of migrations
	entriesCount, err := countMigrationFiles()
	if err != nil {
		t.Fatalf("Failed to count migration files: %v", err)
	}

	if count != entriesCount {
		t.Errorf("Expected %d migrations, but found %d", entriesCount, count)
	}

	// Verify key tables exist
	tables := []string{
		"bought_coins",
		"take_profit_levels",
		"new_coins",
		"purchase_decisions",
		"log_events",
	}

	for _, table := range tables {
		var exists int
		query := `SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`
		if err := db.QueryRow(query, table).Scan(&exists); err != nil {
			t.Fatalf("Failed to check if table %s exists: %v", table, err)
		}
		if exists != 1 {
			t.Errorf("Table %s was not created by migrations", table)
		}
	}
}

func countMigrationFiles() (int, error) {
	// This is a helper function to count the number of migration files
	// Implementation depends on how you access the embedded filesystem in tests
	// For simplicity, you might hard-code the expected count for now
	return 2, nil // Assuming we have 2 migration files
}
```

## 5. Schema Documentation

Consider generating schema documentation automatically from your migrations:

```go
// internal/platform/database/schema_doc.go
package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// GenerateSchemaDoc writes a markdown document describing the database schema
func GenerateSchemaDoc(db *sql.DB, outputPath string) error {
	// Get all tables
	rows, err := db.Query(`
		SELECT name FROM sqlite_master 
		WHERE type='table' AND name NOT LIKE 'sqlite_%' AND name != 'schema_migrations'
	`)
	if err != nil {
		return fmt.Errorf("querying tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scanning table name: %w", err)
		}
		tables = append(tables, name)
	}

	// Create output file
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer f.Close()

	// Write header
	fmt.Fprintln(f, "# Database Schema")
	fmt.Fprintln(f, "\nGenerated from SQLite schema on", "")
	fmt.Fprintln(f)

	// For each table, get its structure
	for _, table := range tables {
		fmt.Fprintf(f, "## %s\n\n", table)

		// Get table info
		rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
		if err != nil {
			return fmt.Errorf("getting table info for %s: %w", table, err)
		}

		// Use tabwriter for nice formatting
		w := tabwriter.NewWriter(f, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "Name\tType\tNot Null\tDefault\tPrimary Key\t")
		fmt.Fprintln(w, "----\t----\t--------\t-------\t-----------\t")

		for rows.Next() {
			var cid, notnull, pk int
			var name, typ, dflt_value sql.NullString
			if err := rows.Scan(&cid, &name, &typ, &notnull, &dflt_value, &pk); err != nil {
				rows.Close()
				return fmt.Errorf("scanning column info: %w", err)
			}

			defaultVal := "NULL"
			if dflt_value.Valid {
				defaultVal = dflt_value.String
			}

			fmt.Fprintf(w, "%s\t%s\t%v\t%s\t%v\t\n",
				name.String, typ.String, notnull == 1, defaultVal, pk == 1)
		}
		rows.Close()
		w.Flush()
		fmt.Fprintln(f)

		// Get indexes
		rows, err = db.Query(fmt.Sprintf("PRAGMA index_list(%s)", table))
		if err != nil {
			return fmt.Errorf("getting indexes for %s: %w", table, err)
		}

		var indexes []string
		for rows.Next() {
			var seq, unique int
			var name, origin string
			var partial sql.NullString
			if err := rows.Scan(&seq, &name, &unique, &origin, &partial); err != nil {
				rows.Close()
				return fmt.Errorf("scanning index info: %w", err)
			}
			
			// Skip automatically generated indexes
			if strings.HasPrefix(name, "sqlite_") {
				continue
			}
			
			indexes = append(indexes, name)
		}
		rows.Close()

		if len(indexes) > 0 {
			fmt.Fprintf(f, "### Indexes\n\n")
			for _, idx := range indexes {
				fmt.Fprintf(f, "- %s\n", idx)
			}
			fmt.Fprintln(f)
		}
	}

	return nil
}
```

## 6. Handling Production Deployments

For production, consider more robust migration strategies:

1. Always test migrations on a copy of the production database
2. Consider taking a backup before applying migrations
3. For critical systems, implement a rollback plan

Here's a simple function to help with this:

```go
// SafeMigrate performs a backup before migration
func SafeMigrate(db *sql.DB, backupDir string) error {
	// Create backup
	backupPath, err := Backup(db, backupDir)
	if err != nil {
		return fmt.Errorf("creating backup before migration: %w", err)
	}
	log.Printf("Created database backup at %s", backupPath)

	// Apply migrations
	if err := Migrate(db); err != nil {
		return fmt.Errorf("applying migrations: %w", err)
	}

	return nil
}
```

## Next Steps

With migrations set up, you're ready to implement the repositories that will interact with your SQLite database. See `05c-sqlite-repositories.md` for details.

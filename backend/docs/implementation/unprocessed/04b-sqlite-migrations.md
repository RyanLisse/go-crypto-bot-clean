# SQLite Migrations

This document covers the implementation of a simple, dependency-free approach to managing SQLite database migrations for the Go crypto trading bot.

## 1. Migration System Overview

The migration system uses embedded SQL scripts and a version tracking table to manage database schema changes in a reliable, version-controlled manner.

### 1.1 Directory Structure

```
internal/platform/database/
├── migrations/
│   ├── 001_initial_schema.sql
│   ├── 002_add_take_profit_levels.sql
│   └── 003_add_indexes.sql
├── migrate.go
└── sqlite.go
```

## 2. Migration System Implementation

First, implement the migration system in a dedicated file:

```go
// internal/platform/database/migrate.go
package database

import (
    "database/sql"
    "embed"
    "fmt"
    "io/fs"
    "log"
    "path"
    "sort"
    "strings"

    "github.com/jmoiron/sqlx"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations runs all pending database migrations
func RunMigrations(db *sqlx.DB) error {
    // Ensure migrations table exists
    if err := createMigrationsTable(db); err != nil {
        return fmt.Errorf("failed to create migrations table: %w", err)
    }

    // Get current migration version
    currentVersion, err := getCurrentVersion(db)
    if err != nil {
        return fmt.Errorf("failed to get current migration version: %w", err)
    }

    // Get available migrations
    migrations, err := listMigrations()
    if err != nil {
        return fmt.Errorf("failed to list migrations: %w", err)
    }

    // Sort migrations by version
    sort.Strings(migrations)

    // Apply pending migrations
    for _, migrationFile := range migrations {
        version := extractVersionFromFilename(migrationFile)
        
        if version <= currentVersion {
            log.Printf("Skipping migration %s (already applied)", migrationFile)
            continue
        }

        log.Printf("Applying migration %s", migrationFile)
        if err := applyMigration(db, migrationFile); err != nil {
            return fmt.Errorf("failed to apply migration %s: %w", migrationFile, err)
        }

        if err := updateVersion(db, version); err != nil {
            return fmt.Errorf("failed to update migration version: %w", err)
        }
    }

    log.Printf("Database is now at version %d", currentVersion)
    return nil
}

// createMigrationsTable creates the schema_migrations table if it doesn't exist
func createMigrationsTable(db *sqlx.DB) error {
    query := `
    CREATE TABLE IF NOT EXISTS schema_migrations (
        version INTEGER PRIMARY KEY,
        applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    `
    _, err := db.Exec(query)
    return err
}

// getCurrentVersion gets the current migration version from the database
func getCurrentVersion(db *sqlx.DB) (int, error) {
    var version int
    err := db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&version)
    return version, err
}

// listMigrations returns a list of migration files from the embedded filesystem
func listMigrations() ([]string, error) {
    var migrations []string
    
    entries, err := migrationsFS.ReadDir("migrations")
    if err != nil {
        return nil, err
    }
    
    for _, entry := range entries {
        if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
            migrations = append(migrations, entry.Name())
        }
    }
    
    return migrations, nil
}

// extractVersionFromFilename extracts the version number from a migration filename
func extractVersionFromFilename(filename string) int {
    var version int
    fmt.Sscanf(filename, "%d_", &version)
    return version
}

// applyMigration applies a single migration file
func applyMigration(db *sqlx.DB, filename string) error {
    // Read migration content
    content, err := migrationsFS.ReadFile(path.Join("migrations", filename))
    if err != nil {
        return err
    }
    
    // Execute migration in a transaction
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()
    
    if _, err = tx.Exec(string(content)); err != nil {
        return err
    }
    
    return tx.Commit()
}

// updateVersion updates the migration version in the database
func updateVersion(db *sqlx.DB, version int) error {
    _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version)
    return err
}
```

## 3. Migration SQL Scripts

Let's examine some example migration scripts:

### 3.1 Initial Schema Migration

```sql
-- migrations/001_initial_schema.sql

-- Create bought_coins table
CREATE TABLE IF NOT EXISTS bought_coins (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    symbol TEXT NOT NULL,
    purchase_price REAL NOT NULL,
    quantity REAL NOT NULL,
    purchased_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted INTEGER DEFAULT 0,
    sold_at TIMESTAMP
);

-- Create new_coins table
CREATE TABLE IF NOT EXISTS new_coins (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    symbol TEXT NOT NULL UNIQUE,
    discovered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create purchase_decisions table
CREATE TABLE IF NOT EXISTS purchase_decisions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    symbol TEXT NOT NULL,
    decision TEXT NOT NULL,
    reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create log_events table
CREATE TABLE IF NOT EXISTS log_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    level TEXT NOT NULL,
    message TEXT NOT NULL,
    context TEXT
);
```

### 3.2 Adding Features Migration

```sql
-- migrations/002_add_take_profit_levels.sql

-- Add positions table for advanced position management
CREATE TABLE IF NOT EXISTS positions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    symbol TEXT NOT NULL,
    entry_price REAL NOT NULL,
    quantity REAL NOT NULL,
    status TEXT NOT NULL DEFAULT 'open',
    stop_loss REAL,
    take_profit REAL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    closed_at TIMESTAMP,
    profit_loss REAL,
    profit_loss_percent REAL,
    notes TEXT
);

-- Add additional fields to bought_coins table
ALTER TABLE bought_coins ADD COLUMN stop_loss REAL;
ALTER TABLE bought_coins ADD COLUMN take_profit REAL;
ALTER TABLE bought_coins ADD COLUMN strategy TEXT;
```

### 3.3 Adding Indexes Migration

```sql
-- migrations/003_add_indexes.sql

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_bought_coins_symbol ON bought_coins(symbol);
CREATE INDEX IF NOT EXISTS idx_bought_coins_purchased_at ON bought_coins(purchased_at);
CREATE INDEX IF NOT EXISTS idx_new_coins_discovered_at ON new_coins(discovered_at);
CREATE INDEX IF NOT EXISTS idx_purchase_decisions_symbol ON purchase_decisions(symbol);
CREATE INDEX IF NOT EXISTS idx_purchase_decisions_created_at ON purchase_decisions(created_at);
CREATE INDEX IF NOT EXISTS idx_positions_symbol ON positions(symbol);
CREATE INDEX IF NOT EXISTS idx_positions_status ON positions(status);
```

## 4. Database Schema Documentation Generator

To automatically generate documentation about the database schema, implement a tool that reads the database structure:

```go
// cmd/tools/genschema/main.go
package main

import (
    "database/sql"
    "flag"
    "fmt"
    "log"
    "os"
    "strings"
    "text/template"

    _ "github.com/mattn/go-sqlite3"
)

func main() {
    dbPath := flag.String("db", "data/cryptobot.db", "Path to SQLite database")
    outputPath := flag.String("output", "docs/schema.md", "Output path for schema documentation")
    flag.Parse()

    db, err := sql.Open("sqlite3", *dbPath)
    if err != nil {
        log.Fatalf("Failed to open database: %v", err)
    }
    defer db.Close()

    if err := GenerateSchemaDoc(db, *outputPath); err != nil {
        log.Fatalf("Failed to generate schema documentation: %v", err)
    }

    log.Printf("Schema documentation generated at %s", *outputPath)
}

func GenerateSchemaDoc(db *sql.DB, outputPath string) error {
    // Get all tables
    rows, err := db.Query(`
        SELECT name FROM sqlite_master 
        WHERE type='table' AND name NOT LIKE 'sqlite_%' AND name != 'schema_migrations'
    `)
    if err != nil {
        return err
    }
    defer rows.Close()

    var tables []string
    for rows.Next() {
        var tableName string
        if err := rows.Scan(&tableName); err != nil {
            return err
        }
        tables = append(tables, tableName)
    }

    // Create template data
    type Column struct {
        Name    string
        Type    string
        NotNull bool
        Default string
        PK      bool
    }

    type Table struct {
        Name    string
        Columns []Column
        Indexes []string
    }

    var tableData []Table

    // Get column info for each table
    for _, tableName := range tables {
        table := Table{Name: tableName}

        // Get columns
        colRows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))
        if err != nil {
            return err
        }

        for colRows.Next() {
            var cid int
            var name, colType, defaultValue string
            var notNull, pk int

            if err := colRows.Scan(&cid, &name, &colType, &notNull, &defaultValue, &pk); err != nil {
                colRows.Close()
                return err
            }

            table.Columns = append(table.Columns, Column{
                Name:    name,
                Type:    colType,
                NotNull: notNull > 0,
                Default: defaultValue,
                PK:      pk > 0,
            })
        }
        colRows.Close()

        // Get indexes
        idxRows, err := db.Query(fmt.Sprintf("PRAGMA index_list(%s)", tableName))
        if err != nil {
            return err
        }

        for idxRows.Next() {
            var seq int
            var name, origin string
            var unique, partial int

            if err := idxRows.Scan(&seq, &name, &unique, &origin, &partial); err != nil {
                idxRows.Close()
                return err
            }

            // Skip auto-generated indexes
            if strings.HasPrefix(name, "sqlite_") {
                continue
            }

            table.Indexes = append(table.Indexes, name)
        }
        idxRows.Close()

        tableData = append(tableData, table)
    }

    // Create template
    tmpl := template.Must(template.New("schema").Parse(`# Database Schema

This document describes the database schema for the cryptocurrency trading bot.

{{range .}}
## {{.Name}}

{{if .Columns}}
| Column | Type | Constraints | Default |
|--------|------|-------------|---------|
{{range .Columns}}| {{.Name}} | {{.Type}} | {{if .PK}}PRIMARY KEY{{end}}{{if .NotNull}}{{if .PK}}, {{end}}NOT NULL{{end}} | {{.Default}} |
{{end}}{{end}}

{{if .Indexes}}
### Indexes
{{range .Indexes}}
- {{.}}
{{end}}
{{end}}

{{end}}
`))

    // Write to file
    file, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer file.Close()

    return tmpl.Execute(file, tableData)
}
```

## 5. Best Practices for Database Migrations

1. **Version Numbers**: Use sequential numbers for migration files (001, 002, etc.)
2. **One-way Migrations**: Design migrations to be forward-only; do not rely on rollbacks
3. **Small, Incremental Changes**: Keep each migration focused on a specific change
4. **Transaction Safety**: Ensure migrations run in transactions when possible
5. **Testing**: Test migrations on a copy of production data before deploying

For more detailed repository implementations, see [04c-repository-implementations.md](04c-repository-implementations.md).

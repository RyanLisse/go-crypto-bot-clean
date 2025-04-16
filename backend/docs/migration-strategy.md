# Database Migration Strategy

This document outlines the standardized approach to database migrations for the application.

## Overview

After reviewing the current codebase and considering maintainability, developer experience, and operational requirements, we have standardized on **GORM AutoMigrate** as our primary migration strategy.

## Migration Strategy: GORM AutoMigrate

### Why GORM AutoMigrate?

1. **Simplicity**: AutoMigrate provides a straightforward approach to schema evolution without requiring manual SQL script management
2. **Code-First Design**: Entity models in Go code serve as the source of truth for database schema
3. **Automatic Execution**: Migrations run automatically during application startup
4. **Incremental Updates**: AutoMigrate only adds new columns/tables and won't delete or modify existing data
5. **Integration**: Already well-integrated with our persistence layer (GORM)

### Implementation Guidelines

#### Core Migration Function

The primary migration function is defined in `internal/adapter/persistence/gorm/auto_migrate.go`:

```go
func AutoMigrate(db *gorm.DB) error {
    // Define entities to migrate in order of dependencies
    entities := []interface{}{
        &entity.User{},
        &entity.Wallet{},
        // Add other entities here in dependency order
    }
    
    // Run migrations
    for _, entity := range entities {
        if err := db.AutoMigrate(entity); err != nil {
            return fmt.Errorf("failed to migrate %T: %w", entity, err)
        }
        log.Printf("Migrated %T", entity)
    }
    
    return nil
}
```

#### Integration with Application Startup

Migrations should be executed during application startup in the following manner:

```go
func main() {
    // Initialize database connection
    db, err := database.Connect()
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    
    // Run migrations
    if err := gorm.AutoMigrate(db); err != nil {
        log.Fatalf("Database migration failed: %v", err)
    }
    
    // Continue with application startup...
}
```

### Version Control and Entity Updates

When modifying entity structures:

1. **Add new fields as nullable**: When adding new fields to entities, ensure they're nullable to avoid migration issues with existing data
2. **Document changes**: Add comments to the entity struct indicating when fields were added or modified
3. **Use tags appropriately**: Utilize GORM tags to control database behavior:
   ```go
   type User struct {
       ID        uint      `gorm:"primaryKey"`
       CreatedAt time.Time
       UpdatedAt time.Time
       DeletedAt gorm.DeletedAt `gorm:"index"`
       Email     string    `gorm:"uniqueIndex;size:255"`
       Username  string    `gorm:"size:50;not null"`
       // Added in v1.2.0
       AvatarURL string    `gorm:"size:255;default:null"`
   }
   ```

## Testing Migrations

### Development Environment

In the development environment, migrations should be run automatically with each application restart. This ensures that developers always work with the latest schema.

### Test Environment

For tests, consider using:

1. **Test Database**: Use a separate test database for running tests
2. **Transaction Rollback**: Run tests within transactions that are rolled back after completion
3. **Schema Reset**: Reset the schema before running the test suite

Example test helper:

```go
func SetupTestDB() (*gorm.DB, func()) {
    db, err := database.ConnectTest()
    if err != nil {
        log.Fatalf("Failed to connect to test database: %v", err)
    }
    
    // Run migrations
    if err := gorm.AutoMigrate(db); err != nil {
        log.Fatalf("Test database migration failed: %v", err)
    }
    
    return db, func() {
        // Cleanup function
        sqlDB, _ := db.DB()
        sqlDB.Close()
    }
}
```

## Production Deployment

For production environments:

1. **Safe Mode**: Consider enabling GORM's DryRun mode first to validate migrations without applying them
2. **Scheduled Migration**: Run migrations during maintenance windows when possible
3. **Backup Before Migration**: Always back up the database before running migrations
4. **Monitoring**: Log and monitor migration execution and duration

## Limitations and Considerations

1. **Data Loss Prevention**: AutoMigrate will not drop columns or tables. For such operations, manual SQL migrations are required.
2. **Complex Changes**: For complex schema changes (renaming columns, changing types), manual migrations may be necessary.
3. **Large Databases**: For very large databases, consider performance impact and potentially use manual migrations with optimized approaches.

## Advanced Migration Needs

If your schema changes require operations not supported by AutoMigrate:

1. **Create a manual migration script** in `internal/adapter/persistence/gorm/migrations/`
2. **Register the migration** in the AutoMigrate function
3. **Document the manual migration** in code comments and update this document

Example manual migration:

```go
func MigrateRenameColumn(db *gorm.DB) error {
    return db.Exec("ALTER TABLE users RENAME COLUMN old_name TO new_name").Error
}
```

## Conclusion

By standardizing on GORM AutoMigrate, we achieve a balance between simplicity and functionality. This approach provides a straightforward, code-first migration strategy that meets our current needs while allowing for advanced customization when necessary. 
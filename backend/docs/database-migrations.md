# Database Migration Strategy

## Chosen Approach: GORM AutoMigrate

This project standardizes on using **GORM's AutoMigrate** feature for all database schema migrations. We have removed support for SQL migration files and external migration tools (such as golang-migrate) to simplify development and onboarding.

### Why AutoMigrate?
- **Simplicity:** Schema changes are handled directly via Go structs and GORM tags, reducing context switching.
- **Consistency:** All migrations are versioned and executed in Go, ensuring that the schema always matches the code.
- **Developer Experience:** No need to maintain separate SQL migration files or learn additional migration tools.

## Migration Workflow
1. **Define or Update Models:**
   - Add or modify structs in the `internal/adapter/persistence/gorm/entity/` directory.
   - Use GORM struct tags to specify column types, constraints, and indexes.
2. **Run AutoMigrate:**
   - Migrations are triggered automatically on application startup or via a dedicated migration command.
   - All models are registered in the `AutoMigrateModels` function in `internal/adapter/persistence/gorm/db.go`.
3. **Testing:**
   - All migrations must be tested in a staging environment before production deployment.
   - Use the test suite and/or spin up a local database to verify schema changes.
4. **Onboarding:**
   - New developers only need to run the application or migration command to ensure their database is up to date.

## Removing Redundant Migration Methods
- All SQL migration files and references to tools like `golang-migrate` have been deprecated and removed.
- The only supported migration entrypoint is via GORM's AutoMigrate.

## Onboarding Instructions
1. **Clone the repository.**
2. **Configure your database connection** (see `.env.example`).
3. **Run the backend application** or the migration command (if available):
   ```sh
   go run main.go
   # or, if available
   go run cmd/migrate/main.go
   ```
   This will automatically apply all pending migrations using GORM AutoMigrate.
4. **Verify your database schema** using your database client of choice.

## Notes
- If you encounter issues with AutoMigrate, check the logs for errors and ensure your models are registered in `AutoMigrateModels`.
- For advanced schema changes (e.g., data migrations, complex index operations), consider writing custom Go migration functions.

---

_Last updated: 2025-04-15_

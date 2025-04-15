# Deprecated SQL Migrations

This directory previously contained raw SQL migration files for use with external migration tools (e.g., golang-migrate). As of 2025-04-15, all database migrations are managed exclusively via GORM's AutoMigrate feature.

## Migration Policy
- **Do not add new SQL migration files.**
- **Do not use external migration tools.**
- All schema changes must be made via Go structs and GORM tags, and applied with AutoMigrate.

See `../docs/database-migrations.md` for the current migration workflow and onboarding instructions.

---

_This directory is retained for historical reference only. All files in this directory may be deleted in a future cleanup._

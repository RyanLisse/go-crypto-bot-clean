-- DEPRECATED: This SQL migration file is no longer used. All database migrations are now handled via GORM AutoMigrate.
-- See docs/database-migrations.md for details.
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indices
CREATE INDEX IF NOT EXISTS idx_status_records_type ON status_records(type);
CREATE INDEX IF NOT EXISTS idx_status_records_component_name ON status_records(component_name);
CREATE INDEX IF NOT EXISTS idx_status_records_created_at ON status_records(created_at);

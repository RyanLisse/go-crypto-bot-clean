CREATE TABLE log_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    level TEXT NOT NULL,
    component TEXT NOT NULL,
    message TEXT NOT NULL,
    metadata TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_log_events_level ON log_events(level);
CREATE INDEX idx_log_events_component ON log_events(component);
CREATE INDEX idx_log_events_created_at ON log_events(created_at);
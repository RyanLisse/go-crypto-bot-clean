-- +goose Up
-- +goose StatementBegin
CREATE TABLE balance_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    balance REAL NOT NULL,
    timestamp DATETIME NOT NULL
);

CREATE INDEX idx_balance_history_timestamp ON balance_history(timestamp);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE balance_history;
-- +goose StatementEnd

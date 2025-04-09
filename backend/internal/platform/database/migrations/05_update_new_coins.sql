-- Add new columns to the new_coins table
ALTER TABLE new_coins ADD COLUMN first_open_time TIMESTAMP;
ALTER TABLE new_coins ADD COLUMN is_upcoming BOOLEAN NOT NULL DEFAULT 0;

-- Create index for first_open_time
CREATE INDEX idx_new_coins_first_open_time ON new_coins(first_open_time);

-- Create index for upcoming coins
CREATE INDEX idx_new_coins_upcoming ON new_coins(is_upcoming);

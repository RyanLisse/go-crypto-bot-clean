-- Add status and became_tradable_at columns to the new_coins table
ALTER TABLE new_coins ADD COLUMN status TEXT;
ALTER TABLE new_coins ADD COLUMN became_tradable_at TIMESTAMP;

-- Create index for status
CREATE INDEX idx_new_coins_status ON new_coins(status);

-- Create index for became_tradable_at
CREATE INDEX idx_new_coins_became_tradable_at ON new_coins(became_tradable_at);

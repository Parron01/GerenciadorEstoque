-- Add batch_id column if not exists
ALTER TABLE history
ADD COLUMN IF NOT EXISTS batch_id VARCHAR(100);

-- Backfill batch_id for existing records (use their own ID as batch_id)
UPDATE history
SET batch_id = id
WHERE batch_id IS NULL OR batch_id = '';

-- Create an index for faster lookups
CREATE INDEX IF NOT EXISTS idx_history_batch_id ON history(batch_id);

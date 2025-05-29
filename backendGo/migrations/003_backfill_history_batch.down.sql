-- Remove the index
DROP INDEX IF EXISTS idx_history_batch_id;

-- Remove the column
ALTER TABLE history
DROP COLUMN IF EXISTS batch_id;

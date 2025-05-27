DROP INDEX IF EXISTS idx_history_entity_type_entity_id;

ALTER TABLE history
DROP COLUMN IF EXISTS entity_id,
DROP COLUMN IF EXISTS entity_type;

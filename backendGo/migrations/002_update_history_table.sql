ALTER TABLE history
ADD COLUMN IF NOT EXISTS entity_type VARCHAR(50),
ADD COLUMN IF NOT EXISTS entity_id VARCHAR(100);

-- Attempt to backfill existing history records
-- This assumes the 'changes' JSON for products contains a 'productId' field.
DO $$
DECLARE
    rec RECORD;
    json_changes JSONB;
    p_id TEXT;
BEGIN
    FOR rec IN SELECT id, changes FROM history WHERE entity_type IS NULL LOOP
        BEGIN
            json_changes := rec.changes::JSONB;
            -- Try to extract productId, specific to ProductChange structure
            p_id := json_changes->>'productId';

            IF p_id IS NOT NULL THEN
                UPDATE history
                SET entity_type = 'product', entity_id = p_id
                WHERE id = rec.id;
            ELSE
                -- Fallback if productId is not found or structure is different
                UPDATE history
                SET entity_type = 'unknown' -- Or handle as per specific needs
                WHERE id = rec.id;
            END IF;
        EXCEPTION WHEN others THEN
            -- Log error or set a default if JSON parsing fails or key is missing
            RAISE NOTICE 'Could not parse changes for history ID %: %', rec.id, SQLERRM;
            UPDATE history
            SET entity_type = 'error_parsing'
            WHERE id = rec.id;
        END;
    END LOOP;
END $$;

-- Add NOT NULL constraints if desired after backfilling, e.g.:
-- ALTER TABLE history ALTER COLUMN entity_type SET NOT NULL;
-- ALTER TABLE history ALTER COLUMN entity_id SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_history_entity_type_entity_id ON history(entity_type, entity_id);

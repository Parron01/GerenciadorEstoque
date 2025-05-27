DROP TRIGGER IF EXISTS trg_update_product_quantity_after_lot_change ON product_lots;
DROP FUNCTION IF EXISTS update_product_quantity_from_lots();

DROP TRIGGER IF EXISTS set_timestamp_product_lots ON product_lots;
-- Note: The trigger_set_timestamp() function might be used by other tables.
-- Only drop it if you are sure it's exclusively for product_lots.
-- For safety, it's often better to leave shared functions unless explicitly managing their lifecycle.
-- DROP FUNCTION IF EXISTS trigger_set_timestamp();

DROP TABLE IF EXISTS product_lots;

-- Optionally, drop the extension if it was created solely for this and no other table uses it.
-- DROP EXTENSION IF EXISTS "uuid-ossp";

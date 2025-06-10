-- Remove user_id columns and constraints
ALTER TABLE products DROP CONSTRAINT IF EXISTS fk_products_user_id;
ALTER TABLE products DROP COLUMN IF EXISTS user_id;

ALTER TABLE product_lots DROP CONSTRAINT IF EXISTS fk_product_lots_user_id;
ALTER TABLE product_lots DROP COLUMN IF EXISTS user_id;

ALTER TABLE history DROP CONSTRAINT IF EXISTS fk_history_user_id;
ALTER TABLE history DROP COLUMN IF EXISTS user_id;

-- Drop indexes
DROP INDEX IF EXISTS idx_products_user_id;
DROP INDEX IF EXISTS idx_product_lots_user_id;
DROP INDEX IF EXISTS idx_history_user_id;

-- Restore original trigger function
CREATE OR REPLACE FUNCTION update_product_quantity_from_lots()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        UPDATE products
        SET quantity = (SELECT COALESCE(SUM(quantity), 0) FROM product_lots WHERE product_id = OLD.product_id)
        WHERE id = OLD.product_id;
        RETURN OLD;
    ELSE
        UPDATE products
        SET quantity = (SELECT COALESCE(SUM(quantity), 0) FROM product_lots WHERE product_id = NEW.product_id)
        WHERE id = NEW.product_id;
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

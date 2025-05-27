CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS product_lots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id VARCHAR(100) NOT NULL,
    quantity NUMERIC NOT NULL CHECK (quantity >= 0),
    data_validade DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_product
        FOREIGN KEY(product_id)
        REFERENCES products(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_product_lots_product_id ON product_lots(product_id);

-- Function to update product quantity based on its lots
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

-- Trigger to update product quantity when a lot is inserted, updated, or deleted
DROP TRIGGER IF EXISTS trg_update_product_quantity_after_lot_change ON product_lots;
CREATE TRIGGER trg_update_product_quantity_after_lot_change
AFTER INSERT OR UPDATE OR DELETE ON product_lots
FOR EACH ROW EXECUTE FUNCTION update_product_quantity_from_lots();

-- Function to update 'updated_at' timestamp
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for 'updated_at' on product_lots
DROP TRIGGER IF EXISTS set_timestamp_product_lots ON product_lots;
CREATE TRIGGER set_timestamp_product_lots
BEFORE UPDATE ON product_lots
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

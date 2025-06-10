-- Add user_id column to products table
ALTER TABLE products 
ADD COLUMN IF NOT EXISTS user_id INTEGER NOT NULL DEFAULT 1;

-- Add foreign key constraint
ALTER TABLE products 
ADD CONSTRAINT fk_products_user_id 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Add user_id column to product_lots table (inherits from product)
ALTER TABLE product_lots 
ADD COLUMN IF NOT EXISTS user_id INTEGER;

-- Update existing product_lots to have the same user_id as their product
UPDATE product_lots 
SET user_id = (
    SELECT user_id 
    FROM products 
    WHERE products.id = product_lots.product_id
)
WHERE user_id IS NULL;

-- Make user_id NOT NULL after updating existing records
ALTER TABLE product_lots 
ALTER COLUMN user_id SET NOT NULL;

-- Add foreign key constraint to product_lots
ALTER TABLE product_lots 
ADD CONSTRAINT fk_product_lots_user_id 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Add user_id column to history table
ALTER TABLE history 
ADD COLUMN IF NOT EXISTS user_id INTEGER NOT NULL DEFAULT 1;

-- Add foreign key constraint to history
ALTER TABLE history 
ADD CONSTRAINT fk_history_user_id 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_products_user_id ON products(user_id);
CREATE INDEX IF NOT EXISTS idx_product_lots_user_id ON product_lots(user_id);
CREATE INDEX IF NOT EXISTS idx_history_user_id ON history(user_id);

-- Update the trigger function to maintain user_id consistency
CREATE OR REPLACE FUNCTION update_product_quantity_from_lots()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        UPDATE products
        SET quantity = (SELECT COALESCE(SUM(quantity), 0) FROM product_lots WHERE product_id = OLD.product_id AND user_id = OLD.user_id)
        WHERE id = OLD.product_id AND user_id = OLD.user_id;
        RETURN OLD;
    ELSE
        -- Ensure new lots inherit user_id from their product
        IF NEW.user_id IS NULL THEN
            SELECT user_id INTO NEW.user_id FROM products WHERE id = NEW.product_id;
        END IF;
        
        UPDATE products
        SET quantity = (SELECT COALESCE(SUM(quantity), 0) FROM product_lots WHERE product_id = NEW.product_id AND user_id = NEW.user_id)
        WHERE id = NEW.product_id AND user_id = NEW.user_id;
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

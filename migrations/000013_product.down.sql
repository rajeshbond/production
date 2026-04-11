-- Drop table

DROP TRIGGER IF EXISTS trigger_update_product_updated_at ON product;

DROP FUNCTION IF EXISTS update_updated_at_column;

DROP TABLE IF EXISTS product;
-- =========================================
-- DROP TRIGGER
-- =========================================
DROP TRIGGER IF EXISTS trigger_update_pos_updated_at ON product_operation_sequence;

-- =========================================
-- DROP UNIQUE CONSTRAINT
-- =========================================
ALTER TABLE product_operation_sequence
DROP CONSTRAINT IF EXISTS uix_product_sequence;

-- =========================================
-- DROP INDEX
-- =========================================
DROP INDEX IF EXISTS idx_pos_tenant_product;

-- =========================================
-- DROP TABLE
-- =========================================
DROP TABLE IF EXISTS product_operation_sequence;

-- =========================================
-- OPTIONAL: DROP UNIQUE CONSTRAINTS
-- ⚠️ Only if not used elsewhere
-- =========================================
-- ALTER TABLE product DROP CONSTRAINT IF EXISTS uix_product_id_tenant;
-- ALTER TABLE operation_master DROP CONSTRAINT IF EXISTS uix_operation_id_tenant;

-- =========================================
-- OPTIONAL: DROP FUNCTION
-- ⚠️ Only if not used globally
-- =========================================
-- DROP FUNCTION IF EXISTS update_updated_at_column;
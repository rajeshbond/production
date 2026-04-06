-- =========================================
-- Drop trigger
-- =========================================
DROP TRIGGER IF EXISTS trg_update_defect_updated_at ON defect;

-- =========================================
-- Drop function
-- =========================================
DROP FUNCTION IF EXISTS update_updated_at_column;

-- =========================================
-- Drop indexes
-- =========================================
DROP INDEX IF EXISTS uix_defect_active;

DROP INDEX IF EXISTS idx_defect_active;

DROP INDEX IF EXISTS idx_defect_tenant;

-- =========================================
-- Drop table
-- =========================================
DROP TABLE IF EXISTS defect;
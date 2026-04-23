-- =========================================
-- DROP TRIGGER
-- =========================================
DROP TRIGGER IF EXISTS trigger_mold_updated_at ON mold;

-- =========================================
-- DROP FUNCTION
-- =========================================
DROP FUNCTION IF EXISTS update_mold_updated_at;

-- =========================================
-- DROP INDEXES (important before table drop in some tools)
-- =========================================
DROP INDEX IF EXISTS uix_mold_active;

DROP INDEX IF EXISTS idx_mold_tenant;

DROP INDEX IF EXISTS idx_mold_active;

DROP INDEX IF EXISTS idx_mold_no_search;

-- =========================================
-- DROP TABLE
-- =========================================
DROP TABLE IF EXISTS mold;
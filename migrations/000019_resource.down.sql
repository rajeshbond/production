-- =========================================
-- DROP TRIGGER
-- =========================================

DROP TRIGGER IF EXISTS trigger_resource_updated_at ON resource;

-- =========================================
-- DROP FUNCTION
-- =========================================

DROP FUNCTION IF EXISTS update_resource_updated_at;

-- =========================================
-- DROP INDEXES
-- =========================================

DROP INDEX IF EXISTS uix_resource_active;

DROP INDEX IF EXISTS idx_resource_tenant;

DROP INDEX IF EXISTS idx_resource_type;

DROP INDEX IF EXISTS idx_resource_active;

-- =========================================
-- DROP TABLE
-- =========================================

DROP TABLE IF EXISTS resource;
-- =========================================
-- DROP TRIGGER
-- =========================================

DROP TRIGGER IF EXISTS trigger_pl_updated_at ON production_log;

-- =========================================
-- DROP FUNCTION
-- =========================================

DROP FUNCTION IF EXISTS update_production_log_updated_at;

-- =========================================
-- DROP INDEXES
-- =========================================

DROP INDEX IF EXISTS idx_pl_setup;

DROP INDEX IF EXISTS idx_pl_machine;

DROP INDEX IF EXISTS idx_pl_date;

DROP INDEX IF EXISTS idx_pl_tenant;

DROP INDEX IF EXISTS uix_pl_unique_slot;

-- =========================================
-- DROP TABLE
-- =========================================

DROP TABLE IF EXISTS production_log;
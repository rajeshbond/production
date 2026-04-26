-- =========================================
-- DROP INDEXES (downtime)
-- =========================================

DROP INDEX IF EXISTS idx_pdt_tenant;

DROP INDEX IF EXISTS idx_pdt_log;

DROP INDEX IF EXISTS uix_pdt_unique;

-- =========================================
-- DROP TABLE: production_downtime
-- =========================================

DROP TABLE IF EXISTS production_downtime;
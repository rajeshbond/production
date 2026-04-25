-- =========================================
-- DROP INDEXES (defect)
-- =========================================

DROP INDEX IF EXISTS idx_pd_tenant;

DROP INDEX IF EXISTS idx_pd_log;

DROP INDEX IF EXISTS uix_pd_unique;

-- =========================================
-- DROP TABLE: production_defect
-- =========================================

DROP TABLE IF EXISTS production_defect;
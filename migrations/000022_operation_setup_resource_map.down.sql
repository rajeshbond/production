-- =========================================
-- DROP INDEXES (optional but clean)
-- =========================================

DROP INDEX IF EXISTS idx_osrm_resource;

DROP INDEX IF EXISTS idx_osrm_setup;

DROP INDEX IF EXISTS uix_osrm_tenant_setup_resource;

-- =========================================
-- DROP TABLE
-- =========================================

DROP TABLE IF EXISTS operation_setup_resource_map;
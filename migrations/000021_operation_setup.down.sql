-- =========================
-- DOWN
-- =========================

DROP TRIGGER IF EXISTS trigger_op_setup_updated_at ON operation_setup;

DROP FUNCTION IF EXISTS update_op_setup_updated_at;

DROP INDEX IF EXISTS uix_setup_unique;

DROP INDEX IF EXISTS idx_osrm_setup;

DROP INDEX IF EXISTS uix_osrm;

DROP TABLE IF EXISTS operation_setup_resource_map;

DROP TABLE IF EXISTS operation_setup;
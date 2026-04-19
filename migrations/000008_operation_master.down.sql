DROP TRIGGER IF EXISTS trg_update_operation_updated_at ON operation_master;

DROP INDEX IF EXISTS uix_operation_active;

DROP INDEX IF EXISTS idx_operation_tenant;

DROP TABLE IF EXISTS operation_master;
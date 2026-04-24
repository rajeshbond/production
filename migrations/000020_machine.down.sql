DROP TRIGGER IF EXISTS trg_machine_updated_at ON machine;

DROP FUNCTION IF EXISTS update_updated_at_column;

DROP INDEX IF EXISTS idx_machine_tenant;

DROP INDEX IF EXISTS uix_tenant_machine_code;

DROP TABLE IF EXISTS machine;
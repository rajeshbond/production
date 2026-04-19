DROP TRIGGER IF EXISTS trg_update_tenant_shift_updated_at ON tenant_shift;

DROP FUNCTION IF EXISTS update_updated_at_column;

DROP TABLE IF EXISTS tenant_shift;
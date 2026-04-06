DROP TRIGGER IF EXISTS trg_update_downtime_updated_at ON downtime;

DROP INDEX IF EXISTS uix_downtime_active;

DROP INDEX IF EXISTS idx_downtime_tenant;

DROP TABLE IF EXISTS downtime;
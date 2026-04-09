-- Drop indexes first
DROP INDEX IF EXISTS uix_operation_downtime;

DROP INDEX IF EXISTS idx_odtm_tenant;

DROP INDEX IF EXISTS idx_odtm_operation;

DROP INDEX IF EXISTS idx_odtm_downtime;

-- Then drop table
DROP TABLE IF EXISTS operation_downtime_map;
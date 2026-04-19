-- =========================
-- DROP TRIGGER
-- =========================
DROP TRIGGER IF EXISTS trg_update_shift_timing_updated_at ON shift_timing;

-- =========================
-- DROP INDEX
-- =========================
DROP INDEX IF EXISTS idx_shift_tenant_weekday;

-- =========================
-- DROP TABLE
-- =========================
DROP TABLE IF EXISTS shift_timing;
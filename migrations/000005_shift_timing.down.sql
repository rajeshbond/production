-- =========================================
-- DROP TRIGGERS
-- =========================================

DROP TRIGGER IF EXISTS trg_check_shift_overlap ON shift_timing;

DROP TRIGGER IF EXISTS trg_update_shift_timing_updated_at ON shift_timing;

-- =========================================
-- DROP FUNCTIONS
-- =========================================

DROP FUNCTION IF EXISTS check_shift_overlap;
-- ⚠️ Only drop if not used elsewhere
-- DROP FUNCTION IF EXISTS update_updated_at_column;

-- =========================================
-- DROP INDEX
-- =========================================

DROP INDEX IF EXISTS idx_shift_tenant_weekday;

-- =========================================
-- DROP TABLE
-- =========================================

DROP TABLE IF EXISTS shift_timing;
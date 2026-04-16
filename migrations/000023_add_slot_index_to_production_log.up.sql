-- =========================================
-- ADD SLOT INDEX TO PRODUCTION LOG
-- =========================================

ALTER TABLE production_log
ADD COLUMN slot_index INT NOT NULL DEFAULT 1;

-- validation
ALTER TABLE production_log
ADD CONSTRAINT chk_pl_slot_index_positive CHECK (slot_index > 0);

-- =========================================
-- IMPORTANT UNIQUE UPDATE
-- =========================================
-- (prevents duplicate slot entry per machine)

DROP INDEX IF EXISTS uix_pl_unique_slot;

CREATE UNIQUE INDEX uix_pl_unique_slot ON production_log (
    tenant_id,
    machine_id,
    shift_hour_slot_id,
    slot_index
);
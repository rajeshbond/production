-- =========================================
-- UP MIGRATION
-- =========================================

-- =========================================
-- TABLE: production_log
-- =========================================

CREATE TABLE IF NOT EXISTS production_log (
    id BIGSERIAL PRIMARY KEY,

-- 🔗 references
tenant_id BIGINT NOT NULL,
setup_id BIGINT NOT NULL,
machine_id BIGINT NOT NULL,
product_id BIGINT NOT NULL,
operation_id BIGINT NOT NULL,

-- 📅 time tracking
production_date DATE NOT NULL,
shift_name_id BIGINT NOT NULL,
shift_timing_id BIGINT NOT NULL,
shift_hour_slot_id BIGINT NOT NULL,
slot_index INT NOT NULL CHECK (slot_index > 0),

-- 📊 production metrics
target_qty INT CHECK (target_qty >= 0),
actual_qty INT NOT NULL DEFAULT 0 CHECK (actual_qty >= 0),
ok_qty INT NOT NULL DEFAULT 0 CHECK (ok_qty >= 0),
rejected_qty INT NOT NULL DEFAULT 0 CHECK (rejected_qty >= 0),
scrap INT NOT NULL DEFAULT 0 CHECK (scrap >= 0),
remarks TEXT,

-- 👤 audit
created_by BIGINT,
updated_by BIGINT,
deleted_by BIGINT,
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
deleted_at TIMESTAMPTZ,
is_deleted BOOLEAN NOT NULL DEFAULT FALSE,

-- 🔗 foreign keys


CONSTRAINT fk_pl_tenant
        FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,

    CONSTRAINT fk_pl_setup
        FOREIGN KEY (setup_id) REFERENCES operation_setup(id) ON DELETE CASCADE,

    CONSTRAINT fk_pl_machine
        FOREIGN KEY (machine_id) REFERENCES machine(id) ON DELETE CASCADE,

    CONSTRAINT fk_pl_product
        FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE,

    CONSTRAINT fk_pl_operation
        FOREIGN KEY (operation_id) REFERENCES operation_master(id) ON DELETE CASCADE,

    CONSTRAINT fk_pl_shift_name
        FOREIGN KEY (shift_name_id) REFERENCES tenant_shift(id) ON DELETE CASCADE,

    CONSTRAINT fk_pl_shift_timing
        FOREIGN KEY (shift_timing_id) REFERENCES shift_timing(id) ON DELETE CASCADE,

    CONSTRAINT fk_pl_shift_hour_slot
        FOREIGN KEY (shift_hour_slot_id) REFERENCES shift_hour_slot(id) ON DELETE CASCADE
);

-- =========================================
-- UNIQUE INDEX (PREVENT DUPLICATE SLOT)
-- =========================================

CREATE UNIQUE INDEX IF NOT EXISTS uix_pl_unique_slot ON production_log (
    tenant_id,
    machine_id,
    production_date,
    shift_name_id,
    shift_hour_slot_id
)
WHERE
    is_deleted = FALSE;

-- =========================================
-- INDEXES
-- =========================================

CREATE INDEX IF NOT EXISTS idx_pl_tenant ON production_log (tenant_id);

CREATE INDEX IF NOT EXISTS idx_pl_date ON production_log (production_date);

CREATE INDEX IF NOT EXISTS idx_pl_machine ON production_log (machine_id);

CREATE INDEX IF NOT EXISTS idx_pl_setup ON production_log (setup_id);

-- =========================================
-- FUNCTION
-- =========================================

CREATE OR REPLACE FUNCTION update_production_log_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =========================================
-- TRIGGER
-- =========================================

DROP TRIGGER IF EXISTS trigger_pl_updated_at ON production_log;

CREATE TRIGGER trigger_pl_updated_at
BEFORE UPDATE ON production_log
FOR EACH ROW
EXECUTE FUNCTION update_production_log_updated_at();
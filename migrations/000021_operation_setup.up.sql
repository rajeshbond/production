-- =========================================
-- TABLE FIRST
-- =========================================

CREATE TABLE IF NOT EXISTS operation_setup (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    pos_id BIGINT NOT NULL,
    machine_id BIGINT NOT NULL,
    setup_name TEXT,
    target_qty INT NOT NULL CHECK (target_qty > 0),
    cycle_time_sec INT,
    setup_time_min INT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by BIGINT,
    updated_by BIGINT,
    deleted_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT fk_op_setup_pos FOREIGN KEY (pos_id) REFERENCES product_operation_sequence (id) ON DELETE CASCADE
);

-- =========================================
-- INDEXES
-- =========================================

CREATE INDEX IF NOT EXISTS idx_op_setup_tenant_pos ON operation_setup (tenant_id, pos_id);

CREATE INDEX IF NOT EXISTS idx_op_setup_machine ON operation_setup (machine_id);

CREATE INDEX IF NOT EXISTS idx_op_setup_active ON operation_setup (tenant_id)
WHERE
    is_deleted = FALSE;

-- =========================================
-- UNIQUE INDEX
-- =========================================

CREATE UNIQUE INDEX IF NOT EXISTS uix_op_setup_unique ON operation_setup (
    tenant_id,
    pos_id,
    machine_id,
    setup_name
)
WHERE
    is_deleted = FALSE;

-- =========================================
-- FUNCTION
-- =========================================

CREATE OR REPLACE FUNCTION update_op_setup_updated_at()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =========================================
-- TRIGGER (AFTER TABLE EXISTS)
-- =========================================

DROP TRIGGER IF EXISTS trigger_op_setup_updated_at ON operation_setup;

CREATE TRIGGER trigger_op_setup_updated_at
BEFORE UPDATE ON operation_setup
FOR EACH ROW
EXECUTE FUNCTION update_op_setup_updated_at();
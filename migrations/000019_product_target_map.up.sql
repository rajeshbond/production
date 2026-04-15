-- =========================================
-- PRODUCTION TARGET (CORE TABLE)
-- =========================================


CREATE TABLE IF NOT EXISTS production_target (
    id BIGSERIAL PRIMARY KEY,

    tenant_id BIGINT NOT NULL,

    product_id BIGINT NOT NULL,
    operation_id BIGINT NOT NULL,
    machine_id BIGINT NOT NULL,

    process_type TEXT NOT NULL, -- MANUAL / SEMI_AUTOMATIC / AUTOMATIC

    target_per_hour INT NOT NULL,
    expected_efficiency INT,

-- 🔥 IMPORTANT
resource_signature TEXT NOT NULL DEFAULT '',

-- ✅ Audit
created_by BIGINT, updated_by BIGINT,

-- ✅ Soft delete
is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
deleted_by BIGINT,
deleted_at TIMESTAMPTZ,

-- ✅ Timestamps
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

-- =====================================
-- 🔗 FOREIGN KEYS
-- =====================================

CONSTRAINT fk_pt_tenant FOREIGN KEY (tenant_id) REFERENCES tenant (id) ON DELETE CASCADE,
CONSTRAINT fk_pt_product FOREIGN KEY (product_id) REFERENCES product (id) ON DELETE CASCADE,
CONSTRAINT fk_pt_operation FOREIGN KEY (operation_id) REFERENCES operation_master (id) ON DELETE CASCADE,
CONSTRAINT fk_pt_machine FOREIGN KEY (machine_id) REFERENCES machine (id) ON DELETE CASCADE,

-- =====================================
-- ✅ VALIDATIONS
-- =====================================

CONSTRAINT chk_process_type
        CHECK (process_type IN ('MANUAL', 'SEMI_AUTOMATIC', 'AUTOMATIC')),

    CONSTRAINT chk_target_positive
        CHECK (target_per_hour > 0),

    CONSTRAINT chk_efficiency_range
        CHECK (
            expected_efficiency IS NULL OR
            (expected_efficiency BETWEEN 0 AND 100)
        )
);

-- =========================================
-- UNIQUE INDEX (COMBINATION LEVEL)
-- =========================================

CREATE UNIQUE INDEX IF NOT EXISTS uix_pt_unique_combination ON production_target (
    tenant_id,
    product_id,
    operation_id,
    machine_id,
    process_type,
    resource_signature
)
WHERE
    is_deleted = FALSE;

-- =========================================
-- INDEXES
-- =========================================

CREATE INDEX IF NOT EXISTS idx_pt_tenant ON production_target (tenant_id);

CREATE INDEX IF NOT EXISTS idx_pt_machine ON production_target (machine_id);

CREATE INDEX IF NOT EXISTS idx_pt_product_operation ON production_target (product_id, operation_id);
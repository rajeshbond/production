-- =========================================
-- PRODUCTION TARGET RESOURCE MAP
-- =========================================


CREATE TABLE IF NOT EXISTS production_target_resource_map (
    id BIGSERIAL PRIMARY KEY,

    tenant_id BIGINT NOT NULL,
    production_target_id BIGINT NOT NULL,
    resource_id BIGINT NOT NULL,

    created_at TIMESTAMPTZ DEFAULT NOW(),

-- =====================================
-- 🔗 FOREIGN KEYS
-- =====================================

CONSTRAINT fk_ptr_tenant FOREIGN KEY (tenant_id) REFERENCES tenant (id) ON DELETE CASCADE,
CONSTRAINT fk_ptr_target FOREIGN KEY (production_target_id) REFERENCES production_target (id) ON DELETE CASCADE,
CONSTRAINT fk_ptr_resource FOREIGN KEY (resource_id) REFERENCES resource_master (id) ON DELETE CASCADE,

-- =====================================
-- UNIQUE CONSTRAINT
-- =====================================

CONSTRAINT uix_ptr_unique
        UNIQUE (production_target_id, resource_id)
);

-- =========================================
-- INDEX
-- =========================================

CREATE INDEX IF NOT EXISTS idx_ptr_target ON production_target_resource_map (production_target_id);
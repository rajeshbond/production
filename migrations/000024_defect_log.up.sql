-- =========================================
-- TABLE: production_defect
-- =========================================


CREATE TABLE IF NOT EXISTS production_defect (
    id BIGSERIAL PRIMARY KEY,

    tenant_id BIGINT NOT NULL,
    production_log_id BIGINT NOT NULL,
    defect_id BIGINT NOT NULL,

    qty INT NOT NULL CHECK (qty > 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

-- =========================================
-- FK
-- =========================================


CONSTRAINT fk_pd_tenant
        FOREIGN KEY (tenant_id)
        REFERENCES tenant(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_pd_production_log
        FOREIGN KEY (production_log_id)
        REFERENCES production_log(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_pd_defect
        FOREIGN KEY (defect_id)
        REFERENCES defect(id)
        ON DELETE CASCADE
);

-- =========================================
-- UNIQUE (prevent duplicate defect entry)
-- =========================================

CREATE UNIQUE INDEX IF NOT EXISTS uix_pd_unique ON production_defect (
    tenant_id,
    production_log_id,
    defect_id
);

-- =========================================
-- INDEXES
-- =========================================

CREATE INDEX IF NOT EXISTS idx_pd_log ON production_defect (production_log_id);

CREATE INDEX IF NOT EXISTS idx_pd_tenant ON production_defect (tenant_id);
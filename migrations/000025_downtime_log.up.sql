-- =========================================
-- TABLE: production_downtime
-- =========================================


CREATE TABLE IF NOT EXISTS production_downtime (
    id BIGSERIAL PRIMARY KEY,

    tenant_id BIGINT NOT NULL,
    production_log_id BIGINT NOT NULL,
    downtime_id BIGINT NOT NULL,

    downtime_minutes INT NOT NULL CHECK (downtime_minutes > 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

-- =========================================
-- FK
-- =========================================


CONSTRAINT fk_pdt_tenant
        FOREIGN KEY (tenant_id)
        REFERENCES tenant(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_pdt_production_log
        FOREIGN KEY (production_log_id)
        REFERENCES production_log(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_pdt_downtime
        FOREIGN KEY (downtime_id)
        REFERENCES downtime(id)
        ON DELETE CASCADE
);

-- =========================================
-- UNIQUE (prevent duplicate downtime entry)
-- =========================================

CREATE UNIQUE INDEX IF NOT EXISTS uix_pdt_unique ON production_downtime (
    tenant_id,
    production_log_id,
    downtime_id
);

-- =========================================
-- INDEXES
-- =========================================

CREATE INDEX IF NOT EXISTS idx_pdt_log ON production_downtime (production_log_id);

CREATE INDEX IF NOT EXISTS idx_pdt_tenant ON production_downtime (tenant_id);
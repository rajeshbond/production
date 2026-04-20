-- =========================================
-- OPERATION ↔ DOWNTIME MAP
-- =========================================

CREATE TABLE IF NOT EXISTS operation_downtime_map (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    operation_id BIGINT NOT NULL,
    downtime_id BIGINT NOT NULL,
    created_by BIGINT ,
    updated_by BIGINT ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

-- Tenant FK
CONSTRAINT fk_odm_tenant FOREIGN KEY (tenant_id) REFERENCES tenant (id) ON DELETE CASCADE,

-- Composite FK (tenant-safe)


CONSTRAINT fk_odm_operation 
        FOREIGN KEY (tenant_id, operation_id) 
        REFERENCES operation_master (tenant_id, id) 
        ON DELETE CASCADE,

    CONSTRAINT fk_odm_downtime 
        FOREIGN KEY (tenant_id, downtime_id) 
        REFERENCES downtime (tenant_id, id) 
        ON DELETE CASCADE
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_odm_tenant ON operation_downtime_map (tenant_id);

CREATE INDEX IF NOT EXISTS idx_odm_operation ON operation_downtime_map (operation_id);

-- Unique (tenant-safe)
CREATE UNIQUE INDEX IF NOT EXISTS uix_operation_downtime ON operation_downtime_map (
    tenant_id,
    operation_id,
    downtime_id
);
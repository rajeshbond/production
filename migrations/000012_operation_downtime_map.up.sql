CREATE TABLE IF NOT EXISTS operation_downtime_map (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    operation_id BIGINT NOT NULL,
    downtime_id BIGINT NOT NULL,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_odfm_tenant FOREIGN KEY (tenant_id) REFERENCES tenant (id) ON DELETE CASCADE,
    CONSTRAINT fk_odtm_operation FOREIGN KEY (operation_id) REFERENCES operation_master (id) ON DELETE CASCADE,
    CONSTRAINT fk_odtm_downtime FOREIGN KEY (downtime_id) REFERENCES downtime (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_odtm_tenant ON operation_downtime_map (tenant_id);

CREATE INDEX IF NOT EXISTS idx_odtm_operation ON operation_downtime_map (operation_id);

CREATE INDEX IF NOT EXISTS idx_odtm_downtime ON operation_downtime_map (downtime_id);

CREATE UNIQUE INDEX IF NOT EXISTS uix_operation_downtime ON operation_downtime_map (
    tenant_id,
    operation_id,
    downtime_id
);
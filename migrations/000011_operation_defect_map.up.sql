CREATE TABLE IF NOT EXISTS operation_defect_map (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    operation_id BIGINT NOT NULL,
    defect_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_odfm_tenant FOREIGN KEY (tenant_id) REFERENCES tenant (id) ON DELETE CASCADE,
    CONSTRAINT fk_odfm_operation FOREIGN KEY (operation_id) REFERENCES operation_master (id) ON DELETE CASCADE,
    CONSTRAINT fk_odfm_defect FOREIGN KEY (defect_id) REFERENCES defect (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_odfm_tenant ON operation_defect_map (tenant_id);

CREATE INDEX IF NOT EXISTS idx_odfm_operation ON operation_defect_map (defect_id);

CREATE UNIQUE INDEX IF NOT EXISTS uix_operation_defect ON operation_defect_map (
    tenant_id,
    operation_id,
    defect_id
);
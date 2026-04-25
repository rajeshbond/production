-- =========================
-- UP
-- =========================

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

CREATE TABLE IF NOT EXISTS operation_setup_resource_map (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    setup_id BIGINT NOT NULL,
    resource_id BIGINT NOT NULL,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_osrm_setup FOREIGN KEY (setup_id) REFERENCES operation_setup (id) ON DELETE CASCADE,
    CONSTRAINT fk_osrm_resource FOREIGN KEY (resource_id) REFERENCES resource (id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS uix_osrm ON operation_setup_resource_map (
    tenant_id,
    setup_id,
    resource_id
);

CREATE INDEX IF NOT EXISTS idx_osrm_setup ON operation_setup_resource_map (setup_id);

CREATE UNIQUE INDEX IF NOT EXISTS uix_setup_unique ON operation_setup (tenant_id, pos_id, machine_id)
WHERE
    is_deleted = FALSE;

CREATE OR REPLACE FUNCTION update_op_setup_updated_at()
RETURNS TRIGGER AS $$
BEGIN
 NEW.updated_at = NOW();
 RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_op_setup_updated_at
BEFORE UPDATE ON operation_setup
FOR EACH ROW
EXECUTE FUNCTION update_op_setup_updated_at();
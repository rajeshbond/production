-- =========================================
-- OPERATION MASTER
-- =========================================

CREATE TABLE IF NOT EXISTS operation_master (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    operation_name TEXT NOT NULL,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_at TIMESTAMPTZ,
    deleted_by BIGINT,
    CONSTRAINT fk_operation_tenant FOREIGN KEY (tenant_id) REFERENCES tenant (id) ON DELETE CASCADE
);

-- Index
CREATE INDEX IF NOT EXISTS idx_operation_tenant ON operation_master (tenant_id);

-- Unique (active only)
CREATE UNIQUE INDEX IF NOT EXISTS uix_operation_active ON operation_master (tenant_id, operation_name)
WHERE
    is_deleted = FALSE;

CONSTRAINT uix_operation_id_tenant UNIQUE (id, tenant_id);

-- Trigger function (shared)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger
DROP TRIGGER IF EXISTS trg_update_operation_updated_at ON operation_master;

CREATE TRIGGER trg_update_operation_updated_at
BEFORE UPDATE ON operation_master
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
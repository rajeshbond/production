-- =========================================
-- MACHINE TABLE
-- =========================================

CREATE TABLE IF NOT EXISTS machine (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    machine_code TEXT NOT NULL,
    machine_name TEXT NOT NULL,
    description TEXT,
    capacity TEXT,
    special_notes JSONB,
    created_by BIGINT,
    updated_by BIGINT,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_by BIGINT,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_machine_tenant FOREIGN KEY (tenant_id) REFERENCES tenant (id) ON DELETE CASCADE
);

-- UNIQUE INDEX
CREATE UNIQUE INDEX IF NOT EXISTS uix_tenant_machine_code ON machine (tenant_id, machine_code)
WHERE
    is_deleted = FALSE;

-- INDEX
CREATE INDEX IF NOT EXISTS idx_machine_tenant ON machine (tenant_id);

-- FUNCTION
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- TRIGGER
DROP TRIGGER IF EXISTS trg_machine_updated_at ON machine;

CREATE TRIGGER trg_machine_updated_at
BEFORE UPDATE ON machine
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
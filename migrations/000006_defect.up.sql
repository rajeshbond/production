-- =========================================
-- Create defect table (SaaS + Audit + Soft Delete)
-- =========================================


CREATE TABLE IF NOT EXISTS defect (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    defect_name TEXT NOT NULL,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_at TIMESTAMPTZ,
    deleted_by BIGINT,

    CONSTRAINT fk_defect_tenant 
        FOREIGN KEY (tenant_id) 
        REFERENCES tenant (id) 
        ON DELETE CASCADE,

-- ✅ REQUIRED for composite FK usage
CONSTRAINT uix_defect_tenant_id UNIQUE (tenant_id, id) );

-- =========================================
-- Indexes (SaaS + Performance)
-- =========================================

-- Fast tenant filtering
CREATE INDEX IF NOT EXISTS idx_defect_tenant ON defect (tenant_id);

-- Fast active record filtering
CREATE INDEX IF NOT EXISTS idx_defect_active ON defect (tenant_id, is_deleted);

-- =========================================
-- Unique constraint (active records only, case-insensitive)
-- =========================================

CREATE UNIQUE INDEX IF NOT EXISTS uix_defect_active ON defect (tenant_id, LOWER(defect_name))
WHERE
    is_deleted = FALSE;

-- =========================================
-- Function: auto-update updated_at
-- =========================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =========================================
-- Trigger
-- =========================================

DROP TRIGGER IF EXISTS trg_update_defect_updated_at ON defect;

CREATE TRIGGER trg_update_defect_updated_at
BEFORE UPDATE ON defect
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
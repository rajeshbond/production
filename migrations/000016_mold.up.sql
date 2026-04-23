-- =========================================
-- MOLD TABLE
-- =========================================
CREATE TABLE IF NOT EXISTS mold (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    type TEXT NOT NULL,
    mold_name TEXT NOT NULL,
    mold_no TEXT NOT NULL,
    description TEXT,
    cavities INT NOT NULL CHECK (cavities > 0),
    target_shots INT NOT NULL CHECK (target_shots > 0),
    special_notes JSONB,
    created_by BIGINT,
    updated_by BIGINT,
    deleted_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

-- =========================================
-- PARTIAL UNIQUE INDEX (CASE-INSENSITIVE)
-- =========================================
CREATE UNIQUE INDEX IF NOT EXISTS uix_mold_active ON mold (tenant_id, LOWER(mold_no))
WHERE
    is_deleted = FALSE;

-- =========================================
-- INDEXES (PERFORMANCE)
-- =========================================
CREATE INDEX IF NOT EXISTS idx_mold_tenant ON mold (tenant_id);

CREATE INDEX IF NOT EXISTS idx_mold_active ON mold (tenant_id)
WHERE
    is_deleted = FALSE;

-- Optional: faster lookup by mold_no
CREATE INDEX IF NOT EXISTS idx_mold_no_search ON mold (tenant_id, LOWER(mold_no));

-- =========================================
-- UPDATED_AT TRIGGER FUNCTION
-- =========================================
CREATE OR REPLACE FUNCTION update_mold_updated_at()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =========================================
-- TRIGGER
-- =========================================
DROP TRIGGER IF EXISTS trigger_mold_updated_at ON mold;

CREATE TRIGGER trigger_mold_updated_at
BEFORE UPDATE ON mold
FOR EACH ROW
EXECUTE FUNCTION update_mold_updated_at();
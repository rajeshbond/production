-- =========================================
-- RESOURCE TABLE (MASTER ONLY)
-- =========================================

CREATE TABLE IF NOT EXISTS resource (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    resource_sub_id BIGINT NOT NULL,
    resource_code TEXT NOT NULL,
    resource_name TEXT,
    resource_type TEXT NOT NULL, -- MOLD / FIXTURE / TOOL
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by BIGINT,
    updated_by BIGINT,
    deleted_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

-- =========================================
-- UNIQUE INDEX (ACTIVE ONLY)
-- =========================================

CREATE UNIQUE INDEX IF NOT EXISTS uix_resource_active ON resource (
    tenant_id,
    LOWER(resource_code),
    resource_sub_id
)
WHERE
    is_deleted = FALSE;

-- =========================================
-- INDEXES
-- =========================================

CREATE INDEX IF NOT EXISTS idx_resource_tenant ON resource (tenant_id);

CREATE INDEX IF NOT EXISTS idx_resource_type ON resource (resource_type);

CREATE INDEX IF NOT EXISTS idx_resource_active ON resource (tenant_id)
WHERE
    is_deleted = FALSE;

-- =========================================
-- UPDATED_AT FUNCTION
-- =========================================

CREATE OR REPLACE FUNCTION update_resource_updated_at()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =========================================
-- TRIGGER
-- =========================================

DROP TRIGGER IF EXISTS trigger_resource_updated_at ON resource;

CREATE TRIGGER trigger_resource_updated_at
BEFORE UPDATE ON resource
FOR EACH ROW
EXECUTE FUNCTION update_resource_updated_at();
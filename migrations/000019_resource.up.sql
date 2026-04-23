-- =========================================
-- RESOURCE TABLE
-- =========================================


CREATE TABLE IF NOT EXISTS resource (
    id BIGSERIAL PRIMARY KEY,

    tenant_id BIGINT NOT NULL,

    resource_code TEXT NOT NULL,
    resource_name TEXT,
    resource_type TEXT NOT NULL,  -- MOLD / FIXTURE / TOOL

    description TEXT,

-- optional mapping

mold_id BIGINT,
    fixture_id BIGINT,
    tool_id BIGINT,

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
    LOWER(resource_code)
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
-- FOREIGN KEYS
-- =========================================

ALTER TABLE resource
ADD CONSTRAINT fk_resource_mold FOREIGN KEY (mold_id) REFERENCES mold (id) ON DELETE CASCADE;

ALTER TABLE resource
ADD CONSTRAINT fk_resource_fixture FOREIGN KEY (fixture_id) REFERENCES fixture (id) ON DELETE CASCADE;

-- Optional (if tool table exists)
-- ALTER TABLE resource
-- ADD CONSTRAINT fk_resource_tool
-- FOREIGN KEY (tool_id) REFERENCES tool(id) ON DELETE CASCADE;

-- =========================================
-- UPDATED_AT TRIGGER
-- =========================================

CREATE OR REPLACE FUNCTION update_resource_updated_at()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_resource_updated_at ON resource;

CREATE TRIGGER trigger_resource_updated_at
BEFORE UPDATE ON resource
FOR EACH ROW
EXECUTE FUNCTION update_resource_updated_at();
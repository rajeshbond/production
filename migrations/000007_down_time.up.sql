-- =========================================
-- DOWNTIME
-- =========================================


CREATE TABLE IF NOT EXISTS downtime (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    downtime_name TEXT NOT NULL,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_at TIMESTAMPTZ,
    deleted_by BIGINT,

    CONSTRAINT fk_downtime_tenant 
        FOREIGN KEY (tenant_id) 
        REFERENCES tenant (id) 
        ON DELETE CASCADE,

-- ✅ REQUIRED for composite FK (VERY IMPORTANT)
CONSTRAINT uix_downtime_tenant_id UNIQUE (tenant_id, id) );

-- =========================================
-- Index
-- =========================================

CREATE INDEX IF NOT EXISTS idx_downtime_tenant ON downtime (tenant_id);

-- =========================================
-- Unique (active only, case-insensitive)
-- =========================================

CREATE UNIQUE INDEX IF NOT EXISTS uix_downtime_active ON downtime (
    tenant_id,
    LOWER(downtime_name)
)
WHERE
    is_deleted = FALSE;

-- =========================================
-- Trigger
-- =========================================

DROP TRIGGER IF EXISTS trg_update_downtime_updated_at ON downtime;

CREATE TRIGGER trg_update_downtime_updated_at
BEFORE UPDATE ON downtime
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
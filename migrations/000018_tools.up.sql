-- =========================================
-- TOOL MASTER TABLE (UP MIGRATION)
-- =========================================

-- 1. TABLE
CREATE TABLE IF NOT EXISTS tool_master (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    type TEXT NOT NULL,
    tool_code TEXT NOT NULL,
    tool_name TEXT NOT NULL,
    description TEXT,
    tool_type TEXT,
    unit TEXT,
    cost NUMERIC(12, 2) DEFAULT 0,
    life_cycles BIGINT DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =========================================
-- 2. SOFT DELETE SAFE UNIQUE INDEX
-- =========================================

-- Ensures: only one ACTIVE tool_code per tenant
CREATE UNIQUE INDEX IF NOT EXISTS uix_tool_tenant_active ON tool_master (tenant_id, tool_code)
WHERE
    is_deleted = false;

-- =========================================
-- 3. PERFORMANCE INDEXES
-- =========================================

CREATE INDEX IF NOT EXISTS idx_tool_tenant ON tool_master (tenant_id);

CREATE INDEX IF NOT EXISTS idx_tool_active ON tool_master (tenant_id, is_active);

CREATE INDEX IF NOT EXISTS idx_tool_deleted ON tool_master (tenant_id, is_deleted);

-- =========================================
-- 4. UPDATED_AT TRIGGER FUNCTION
-- =========================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =========================================
-- 5. TRIGGER
-- =========================================

DROP TRIGGER IF EXISTS trigger_tool_updated_at ON tool_master;

CREATE TRIGGER trigger_tool_updated_at
BEFORE UPDATE ON tool_master
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
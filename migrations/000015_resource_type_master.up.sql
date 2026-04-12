-- =========================================
-- RESOURCE TYPE MASTER
-- =========================================

CREATE TABLE IF NOT EXISTS resource_type_master (
    id SERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    type_name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by BIGINT,
    updated_by BIGINT,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_by BIGINT
);

-- =========================================
-- PARTIAL UNIQUE INDEX (FIX FOR SOFT DELETE)
-- =========================================

CREATE UNIQUE INDEX IF NOT EXISTS unique_active_type_per_tenant ON resource_type_master (tenant_id, LOWER(type_name))
WHERE
    is_deleted = FALSE;

-- =========================================
-- VALIDATION (NO EMPTY TYPE NAME)
-- =========================================

ALTER TABLE resource_type_master
ADD CONSTRAINT chk_type_name_not_empty CHECK (TRIM(type_name) <> '');

-- =========================================
-- COMMON FUNCTION FOR updated_at
-- =========================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =========================================
-- TRIGGER: resource_type_master
-- =========================================

DROP TRIGGER IF EXISTS trg_resource_type_updated_at ON resource_type_master;

CREATE TRIGGER trg_resource_type_updated_at
BEFORE UPDATE ON resource_type_master
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- =========================================
-- TRIGGER: resource_master (SAFE VERSION)
-- =========================================

-- Only create if table exists
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 
        FROM information_schema.tables 
        WHERE table_name = 'resource_master'
    ) THEN
        DROP TRIGGER IF EXISTS trg_resource_updated_at ON resource_master;

        CREATE TRIGGER trg_resource_updated_at
        BEFORE UPDATE ON resource_master
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;
END;
$$;
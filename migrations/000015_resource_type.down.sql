-- =========================================
-- DROP TRIGGER: resource_master (SAFE)
-- =========================================
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 
        FROM information_schema.tables 
        WHERE table_name = 'resource_master'
    ) THEN
        DROP TRIGGER IF EXISTS trg_resource_updated_at ON resource_master;
    END IF;
END;
$$;

-- =========================================
-- DROP TRIGGER: resource_type_master
-- =========================================
DROP TRIGGER IF EXISTS trg_resource_type_updated_at ON resource_type_master;

-- =========================================
-- DROP FUNCTION
-- =========================================
DROP FUNCTION IF EXISTS update_updated_at_column;

-- =========================================
-- DROP CONSTRAINT
-- =========================================
ALTER TABLE IF EXISTS resource_type_master
DROP CONSTRAINT IF EXISTS chk_type_name_not_empty;

-- =========================================
-- DROP INDEX
-- =========================================
DROP INDEX IF EXISTS unique_active_type_per_tenant;

-- =========================================
-- DROP TABLE
-- =========================================
DROP TABLE IF EXISTS resource_type_master;
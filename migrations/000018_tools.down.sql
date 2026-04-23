-- =========================================
-- TOOL MASTER DOWN MIGRATION
-- =========================================

-- 1. Drop triggers first
DROP TRIGGER IF EXISTS trigger_tool_updated_at ON tool_master;

-- 2. Drop trigger function (ONLY if not used elsewhere)
DROP FUNCTION IF EXISTS update_updated_at_column;

-- 3. Drop indexes (safe cleanup)
DROP INDEX IF EXISTS uix_tool_tenant_active;

DROP INDEX IF EXISTS idx_tool_tenant;

DROP INDEX IF EXISTS idx_tool_active;

DROP INDEX IF EXISTS idx_tool_deleted;

-- 4. Drop table
DROP TABLE IF EXISTS tool_master;
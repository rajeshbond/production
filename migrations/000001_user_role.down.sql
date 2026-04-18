-- Drop trigger first (dependency)
DROP TRIGGER IF EXISTS trg_update_user_role_updated_at ON user_role;

-- Drop table
DROP TABLE IF EXISTS user_role;

-- Optional: drop function (only if not used elsewhere)
-- ⚠️ Be careful if reused in other tables
DROP FUNCTION IF EXISTS update_updated_at_column ();
-- =========================
-- DROP INDEXES
-- =========================

DROP INDEX IF EXISTS uix_user_employee_active;

DROP INDEX IF EXISTS uix_user_email_active;

DROP INDEX IF EXISTS idx_user_active;

-- =========================
-- DROP TRIGGER
-- =========================

DROP TRIGGER IF EXISTS trg_update_users_updated_at ON user;

-- =========================
-- DROP TABLE
-- =========================

DROP TABLE IF EXISTS users;

-- =========================
-- DROP FUNCTION (OPTIONAL)
-- =========================

-- ⚠️ Only drop if this function is NOT used by other tables
-- DROP FUNCTION IF EXISTS update_updated_at_column();
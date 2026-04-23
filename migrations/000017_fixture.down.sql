DROP TRIGGER IF EXISTS trigger_fixture_updated_at ON fixture;

DROP FUNCTION IF EXISTS update_fixture_updated_at;

DROP INDEX IF EXISTS uix_fixture_active;

DROP INDEX IF EXISTS idx_fixture_tenant;

DROP INDEX IF EXISTS idx_fixture_active;

DROP TABLE IF EXISTS fixture;
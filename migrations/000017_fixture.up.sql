CREATE TABLE IF NOT EXISTS fixture (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    type TEXT NOT NULL,
    fixture_no TEXT NOT NULL,
    fixture_name TEXT,
    description TEXT,
    cavities INT CHECK (cavities > 0),
    life_shots BIGINT CHECK (life_shots > 0),
    fixture_type TEXT,
    material TEXT,
    special_notes JSONB,
    created_by BIGINT,
    updated_by BIGINT,
    deleted_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE UNIQUE INDEX IF NOT EXISTS uix_fixture_active ON fixture (tenant_id, LOWER(fixture_no))
WHERE
    is_deleted = FALSE;

CREATE INDEX IF NOT EXISTS idx_fixture_tenant ON fixture (tenant_id);

CREATE INDEX IF NOT EXISTS idx_fixture_active ON fixture (tenant_id)
WHERE
    is_deleted = FALSE;

-- trigger
CREATE OR REPLACE FUNCTION update_fixture_updated_at()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_fixture_updated_at ON fixture;

CREATE TRIGGER trigger_fixture_updated_at
BEFORE UPDATE ON fixture
FOR EACH ROW
EXECUTE FUNCTION update_fixture_updated_at();
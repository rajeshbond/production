-- Create Tenant Table

CREATE TABLE IF NOT EXISTS tenant (
    id SERIAL PRIMARY KEY,
    tenant_name VARCHAR NOT NULL,
    tenant_code VARCHAR NOT NULL UNIQUE,
    address VARCHAR NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by INTEGER,
    updated_by INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_set_updated_at
BEFORE UPDATE ON tenant
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
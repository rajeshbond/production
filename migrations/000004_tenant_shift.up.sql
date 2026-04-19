-- Tenant shifts
CREATE TABLE IF NOT EXISTS tenant_shift (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER NOT NULL,
    shift_name VARCHAR NOT NULL,
    created_by INTEGER,
    updated_by INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_tenant_shift_tenant FOREIGN KEY (tenant_id) REFERENCES tenant (id) ON DELETE CASCADE,
    CONSTRAINT uix_tenant_shift UNIQUE (tenant_id, shift_name)
);

-- Trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;  
END;
$$ LANGUAGE plpgsql;

-- Trigger
CREATE TRIGGER trg_update_tenant_shift_updated_at
BEFORE UPDATE ON tenant_shift
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
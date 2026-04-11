-- create Product table
CREATE TABLE IF NOT EXISTS product (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER NOT NULL,
    product_name VARCHAR NOT NULL,
    product_no VARCHAR NOT NULL,
    created_by INTEGER,
    updated_by INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_product_tenant FOREIGN KEY (tenant_id) REFERENCES tenant (id) ON DELETE CASCADE,
    CONSTRAINT uix_tenant_product_no UNIQUE (tenant_id, product_no),
    CONSTRAINT uix_tenant_product_name UNIQUE (tenant_id, product_name),
    CONSTRAINT uix_product_id_tenant UNIQUE (id, tenant_id),
);

-- indexes
CREATE INDEX IF NOT EXISTS ix_product_tenant_id ON product (tenant_id);

CREATE INDEX IF NOT EXISTS ix_product_product_no ON product (product_no);

-- function (safe to re-run)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- trigger (make it idempotent)
DROP TRIGGER IF EXISTS trigger_update_product_updated_at ON product;

CREATE TRIGGER trigger_update_product_updated_at
BEFORE UPDATE ON product
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
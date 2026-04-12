
CREATE TABLE IF NOT EXISTS product (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER NOT NULL,
    product_name VARCHAR NOT NULL,
    product_no VARCHAR NOT NULL,

    created_by INTEGER,
    updated_by INTEGER,

-- ✅ Soft delete fields
is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_by INTEGER,
    deleted_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_product_tenant 
        FOREIGN KEY (tenant_id) REFERENCES tenant (id) ON DELETE CASCADE,

    CONSTRAINT uix_tenant_product_no 
        UNIQUE (tenant_id, product_no),

    CONSTRAINT uix_tenant_product_name 
        UNIQUE (tenant_id, product_name)
);
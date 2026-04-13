
CREATE TABLE IF NOT EXISTS resource_master (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,

    resource_name TEXT NOT NULL,
    resource_type_id BIGINT NOT NULL,

    created_by BIGINT,
    updated_by BIGINT,

    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_by BIGINT,
    deleted_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_resource_tenant 
        FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,

-- ✅ SIMPLE FK (BEST)
CONSTRAINT fk_resource_type 
        FOREIGN KEY (resource_type_id)
        REFERENCES resource_type_master(id)
        ON DELETE CASCADE,

    CONSTRAINT uix_resource_name 
        UNIQUE (tenant_id, resource_name)
);
-- =========================================
-- OPERATION MODE
-- =========================================

CREATE TABLE IF NOT EXISTS operation_mode (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    operation_mode_name TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_by BIGINT,
    updated_by BIGINT,
    deleted_by BIGINT,
    deleted_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_operation_mode_tenant FOREIGN KEY (tenant_id) REFERENCES tenant (id) ON DELETE CASCADE
);

-- =========================================
-- UNIQUE (SOFT DELETE SAFE + CASE INSENSITIVE)
-- =========================================

CREATE UNIQUE INDEX IF NOT EXISTS uix_operation_mode ON operation_mode (
    tenant_id,
    LOWER(operation_mode_name)
)
WHERE
    is_deleted = FALSE;

-- =========================================
-- INDEXES
-- =========================================

CREATE INDEX IF NOT EXISTS idx_operation_mode_tenant ON operation_mode (tenant_id);
-- =========================================
-- TABLE: operation_setup_resource_map
-- =========================================

CREATE TABLE IF NOT EXISTS operation_setup_resource_map (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    setup_id BIGINT NOT NULL,
    resource_id BIGINT NOT NULL,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_osrm_setup FOREIGN KEY (setup_id) REFERENCES operation_setup (id) ON DELETE CASCADE,
    CONSTRAINT fk_osrm_resource FOREIGN KEY (resource_id) REFERENCES resource (id) ON DELETE CASCADE
);

-- =========================================
-- UNIQUE INDEX (PREVENT DUPLICATES)
-- =========================================

CREATE UNIQUE INDEX IF NOT EXISTS uix_osrm_tenant_setup_resource ON operation_setup_resource_map (
    tenant_id,
    setup_id,
    resource_id
);

-- =========================================
-- PERFORMANCE INDEXES
-- =========================================

CREATE INDEX IF NOT EXISTS idx_osrm_setup ON operation_setup_resource_map (setup_id);

CREATE INDEX IF NOT EXISTS idx_osrm_resource ON operation_setup_resource_map (resource_id);
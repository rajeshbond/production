-- =========================================
-- PRODUCTION LOG
-- =========================================

CREATE TABLE IF NOT EXISTS production_log (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    machine_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    operation_id BIGINT NOT NULL,
    shift_hour_slot_id BIGINT NOT NULL,
    quantity INT NOT NULL,
    remarks TEXT,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_pl_tenant FOREIGN KEY (tenant_id) REFERENCES tenant (id) ON DELETE CASCADE,
    CONSTRAINT fk_pl_machine FOREIGN KEY (machine_id) REFERENCES machine (id) ON DELETE CASCADE,
    CONSTRAINT fk_pl_product FOREIGN KEY (product_id) REFERENCES product (id) ON DELETE CASCADE,
    CONSTRAINT fk_pl_operation FOREIGN KEY (operation_id) REFERENCES operation_master (id) ON DELETE CASCADE,
    CONSTRAINT fk_pl_slot FOREIGN KEY (shift_hour_slot_id) REFERENCES shift_hour_slot (id) ON DELETE CASCADE
);

-- ONE ENTRY PER MACHINE PER SLOT
CREATE UNIQUE INDEX uix_pl_unique_slot ON production_log (
    tenant_id,
    machine_id,
    shift_hour_slot_id
);

-- =========================================
-- RESOURCE MAP
-- =========================================

CREATE TABLE IF NOT EXISTS production_log_resource (
    id BIGSERIAL PRIMARY KEY,
    production_log_id BIGINT NOT NULL,
    resource_id BIGINT NOT NULL,
    CONSTRAINT fk_plr_log FOREIGN KEY (production_log_id) REFERENCES production_log (id) ON DELETE CASCADE,
    CONSTRAINT fk_plr_resource FOREIGN KEY (resource_id) REFERENCES resource_master (id) ON DELETE CASCADE,
    CONSTRAINT uix_plr UNIQUE (
        production_log_id,
        resource_id
    )
);

-- =========================================
-- DEFECT MASTER
-- =========================================

CREATE TABLE IF NOT EXISTS defect_master (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    defect_name TEXT NOT NULL
);

-- =========================================
-- PRODUCTION DEFECT
-- =========================================

CREATE TABLE IF NOT EXISTS production_log_defect (
    id BIGSERIAL PRIMARY KEY,
    production_log_id BIGINT NOT NULL,
    defect_id BIGINT NOT NULL,
    quantity INT NOT NULL,
    CONSTRAINT fk_pld_log FOREIGN KEY (production_log_id) REFERENCES production_log (id) ON DELETE CASCADE,
    CONSTRAINT fk_pld_defect FOREIGN KEY (defect_id) REFERENCES defect_master (id) ON DELETE CASCADE
);

-- =========================================
-- DOWNTIME MASTER
-- =========================================

CREATE TABLE IF NOT EXISTS downtime_master (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    downtime_name TEXT NOT NULL
);

-- =========================================
-- PRODUCTION DOWNTIME
-- =========================================

CREATE TABLE IF NOT EXISTS production_log_downtime (
    id BIGSERIAL PRIMARY KEY,
    production_log_id BIGINT NOT NULL,
    downtime_id BIGINT NOT NULL,
    minutes INT NOT NULL,
    CONSTRAINT fk_pldt_log FOREIGN KEY (production_log_id) REFERENCES production_log (id) ON DELETE CASCADE,
    CONSTRAINT fk_pldt_downtime FOREIGN KEY (downtime_id) REFERENCES downtime_master (id) ON DELETE CASCADE
);
-- =========================================
-- CREATE TABLE
-- =========================================
CREATE TABLE IF NOT EXISTS product_operation_sequence (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    operation_id BIGINT NOT NULL,
    sequence_no INT NOT NULL CHECK (sequence_no > 0),
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_pos_product FOREIGN KEY (product_id, tenant_id) REFERENCES product (id, tenant_id) ON DELETE CASCADE,
    CONSTRAINT fk_pos_operation FOREIGN KEY (operation_id, tenant_id) REFERENCES operation_master (id, tenant_id) ON DELETE CASCADE
);

-- =========================================
-- UNIQUE CONSTRAINT (DEFERRABLE for reorder)
-- =========================================
ALTER TABLE product_operation_sequence
DROP CONSTRAINT IF EXISTS uix_product_sequence;

ALTER TABLE product_operation_sequence
ADD CONSTRAINT uix_product_sequence UNIQUE (
    tenant_id,
    product_id,
    sequence_no
) DEFERRABLE INITIALLY DEFERRED;

-- =========================================
-- INDEX (for performance)
-- =========================================
CREATE INDEX IF NOT EXISTS idx_pos_tenant_product ON product_operation_sequence (tenant_id, product_id);

-- =========================================
-- UPDATED_AT FUNCTION
-- =========================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =========================================
-- TRIGGER
-- =========================================
DROP TRIGGER IF EXISTS trigger_update_pos_updated_at ON product_operation_sequence;

CREATE TRIGGER trigger_update_pos_updated_at
BEFORE UPDATE ON product_operation_sequence
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
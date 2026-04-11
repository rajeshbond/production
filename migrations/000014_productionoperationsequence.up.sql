-- Create Production Operation Sequence

CREATE TABLE IF NOT EXISTS product_operation_sequence (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    operation_id BIGINT NOT NULL,
    sequence_no INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_pos_product FOREIGN KEY (product_id, tenant_id) REFERENCES product (id, tenant_id) ON DELETE CASCADE,
    CONSTRAINT fk_pos_operation FOREIGN KEY (operation_id, tenant_id) REFERENCES operation_master (id, tenant_id) ON DELETE CASCADE,
    CONSTRAINT uix_product_sequence UNIQUE (
        tenant_id,
        product_id,
        sequence_no
    ),
    CONSTRAINT uix_pos_id_tenant UNIQUE (id, tenant_id)
);
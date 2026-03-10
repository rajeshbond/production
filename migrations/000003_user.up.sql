-- Create User

CREATE TABLE IF NOT EXISTS "user" (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    employee_id VARCHAR(100) NOT NULL,
    user_name VARCHAR(150) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(150),
    password TEXT NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user_tenant FOREIGN KEY (tenant_id) REFERENCES tenant (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_role FOREIGN KEY (role_id) REFERENCES user_role (id) ON DELETE CASCADE,
    CONSTRAINT uix_tenant_employee UNIQUE (tenant_id, employee_id),
    CONSTRAINT uix_tenant_email UNIQUE (tenant_id, email)
);

-- Composite index for performance
CREATE INDEX idx_user_tenant_user_employee ON "user" (
    tenant_id,
    user_name,
    employee_id
);
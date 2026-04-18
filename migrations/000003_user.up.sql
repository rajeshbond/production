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
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_at TIMESTAMPTZ,
    deleted_by BIGINT,
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user_tenant FOREIGN KEY (tenant_id) REFERENCES tenant (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_role FOREIGN KEY (role_id) REFERENCES user_role (id) ON DELETE CASCADE
);

DROP TRIGGER IF EXISTS trg_update_user_updated_at ON "user";

CREATE TRIGGER trg_update_user_updated_at
BEFORE UPDATE ON "user"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE UNIQUE INDEX IF NOT EXISTS uix_user_employee_active ON "user" (tenant_id, LOWER(employee_id))
WHERE
    is_deleted = FALSE;

CREATE UNIQUE INDEX IF NOT EXISTS uix_user_email_active ON "user" (tenant_id, LOWER(email))
WHERE
    is_deleted = FALSE
    AND email IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_user_active ON "user" (
    tenant_id,
    is_deleted,
    is_active
);
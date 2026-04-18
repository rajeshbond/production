-- =========================
-- TENANT TABLE
-- =========================

CREATE TABLE IF NOT EXISTS tenant (
    id SERIAL PRIMARY KEY,
    tenant_name VARCHAR NOT NULL,
    tenant_code VARCHAR NOT NULL,
    address VARCHAR NOT NULL,

-- Contact
contact_person_name VARCHAR(150),
contact_phone VARCHAR(20),
contact_email VARCHAR(150),
is_verified BOOLEAN NOT NULL DEFAULT FALSE,
is_active BOOLEAN NOT NULL DEFAULT TRUE,

-- Soft delete
is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
deleted_at TIMESTAMPTZ,
deleted_by INTEGER,

-- Audit

created_by INTEGER,
    updated_by INTEGER,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =========================
-- FUNCTION (CREATE ONCE)
-- =========================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =========================
-- TRIGGER
-- =========================

DROP TRIGGER IF EXISTS trg_update_tenant_updated_at ON tenant;

CREATE TRIGGER trg_update_tenant_updated_at
BEFORE UPDATE ON tenant
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- =========================
-- UNIQUE INDEXES (ACTIVE ONLY)
-- =========================

CREATE UNIQUE INDEX IF NOT EXISTS uix_tenant_code_active ON tenant (LOWER(tenant_code))
WHERE
    is_deleted = FALSE;

CREATE UNIQUE INDEX IF NOT EXISTS uix_tenant_contact_email_active ON tenant (LOWER(contact_email))
WHERE
    is_deleted = FALSE
    AND contact_email IS NOT NULL;

-- =========================
-- PERFORMANCE INDEX
-- =========================

CREATE INDEX IF NOT EXISTS idx_tenant_active ON tenant (is_deleted, is_active);
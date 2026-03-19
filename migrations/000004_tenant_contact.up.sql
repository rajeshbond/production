-- =========================
-- TENANT CONTACTS TABLE
-- =========================


CREATE TABLE IF NOT EXISTS tenant_contacts (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER NOT NULL,

    contact_person_name VARCHAR(150) NOT NULL,
    contact_phone VARCHAR(20),
    contact_email VARCHAR(150),

    is_primary BOOLEAN DEFAULT FALSE,

-- ✅ Soft delete


is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_tenant_contact
    FOREIGN KEY (tenant_id)
    REFERENCES tenant(id)
    ON DELETE CASCADE
);

-- =========================
-- ✅ PARTIAL UNIQUE INDEX
-- =========================

CREATE UNIQUE INDEX IF NOT EXISTS unique_primary_contact_per_tenant ON tenant_contacts (tenant_id)
WHERE
    is_primary = TRUE
    AND is_deleted = FALSE;

-- =========================
-- TRIGGER
-- =========================

DROP TRIGGER IF EXISTS trg_set_updated_at_contacts ON tenant_contacts;

CREATE TRIGGER trg_set_updated_at_contacts
BEFORE UPDATE ON tenant_contacts
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
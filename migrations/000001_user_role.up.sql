-- Create table
CREATE TABLE IF NOT EXISTS user_role (
    id SERIAL PRIMARY KEY,
    user_role VARCHAR NOT NULL UNIQUE,
    created_by INTEGER,
    updated_by INTEGER,
    deleted_by INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_at TIMESTAMPTZ
);

-- Create function (safe: only if not exists / replace)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Drop trigger if already exists (important for reruns)
DROP TRIGGER IF EXISTS trg_update_user_role_updated_at ON user_role;

-- Create trigger
CREATE TRIGGER trg_update_user_role_updated_at
BEFORE UPDATE ON user_role
FOR EACH ROW
WHEN (OLD IS DISTINCT FROM NEW)
EXECUTE FUNCTION update_updated_at_column();
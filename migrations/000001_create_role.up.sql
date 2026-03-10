-- Create Table User Role

CREATE TABLE IF NOT EXISTS user_role (
    id SERIAL PRIMARY KEY,
    user_role VARCHAR NOT NULL UNIQUE,
    created_by INTEGER,
    updated_by INTEGER,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL
);
-- Create shift_timing table


CREATE TABLE IF NOT EXISTS shift_timing (
    id SERIAL PRIMARY KEY,

    tenant_shift_id INTEGER NOT NULL,
    shift_start TIME NOT NULL,
    shift_end TIME NOT NULL,
    weekday INTEGER NOT NULL, -- 0 = Monday ... 6 = Sunday

    created_by INTEGER,
    updated_by INTEGER,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_shift_tenant
        FOREIGN KEY (tenant_shift_id)
        REFERENCES tenant_shift(id)
        ON DELETE CASCADE,

-- Prevent duplicate shift timing per day
CONSTRAINT uix_shift_timing UNIQUE (
    tenant_shift_id,
    weekday,
    shift_start,
    shift_end
),

-- Allow overnight shifts, but prevent same start & end
CONSTRAINT check_shift_time_valid CHECK (shift_start <> shift_end),

-- Ensure valid weekday (1–7)
CONSTRAINT check_weekday_valid
        CHECK (weekday BETWEEN 1 AND 7)
);

-- Index for faster queries
CREATE INDEX IF NOT EXISTS idx_shift_tenant_weekday ON shift_timing (tenant_shift_id, weekday);

-- Trigger (reuse same function)
DROP TRIGGER IF EXISTS trg_update_shift_timing_updated_at ON shift_timing;

CREATE TRIGGER trg_update_shift_timing_updated_at
BEFORE UPDATE ON shift_timing
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
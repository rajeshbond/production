-- =========================================
-- SHIFT TIMING TABLE
-- =========================================

CREATE TABLE IF NOT EXISTS shift_timing (
    id SERIAL PRIMARY KEY,
    tenant_shift_id INTEGER NOT NULL,
    shift_start TIME NOT NULL,
    shift_end TIME NOT NULL,
    weekday INTEGER NOT NULL, -- 1 = Monday ... 7 = Sunday
    created_by INTEGER,
    updated_by INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_shift_tenant FOREIGN KEY (tenant_shift_id) REFERENCES tenant_shift (id) ON DELETE CASCADE,
    CONSTRAINT uix_shift_timing UNIQUE (
        tenant_shift_id,
        weekday,
        shift_start,
        shift_end
    ),
    CONSTRAINT check_shift_time_valid CHECK (shift_start <> shift_end),
    CONSTRAINT check_weekday_valid CHECK (weekday BETWEEN 1 AND 7)
);

-- Index
CREATE INDEX IF NOT EXISTS idx_shift_tenant_weekday ON shift_timing (tenant_shift_id, weekday);

-- =========================================
-- FUNCTION: update_updated_at
-- =========================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =========================================
-- FUNCTION: Prevent overlap
-- =========================================

CREATE OR REPLACE FUNCTION check_shift_overlap()
RETURNS TRIGGER AS $$
DECLARE
    new_start INT;
    new_end INT;
    existing_start INT;
    existing_end INT;
BEGIN
    new_start := EXTRACT(HOUR FROM NEW.shift_start) * 60 
               + EXTRACT(MINUTE FROM NEW.shift_start);

    new_end := EXTRACT(HOUR FROM NEW.shift_end) * 60 
             + EXTRACT(MINUTE FROM NEW.shift_end);

    -- Handle overnight
    IF new_end <= new_start THEN
        new_end := new_end + 1440;
    END IF;

    FOR existing_start, existing_end IN
        SELECT 
            EXTRACT(HOUR FROM st.shift_start) * 60 
          + EXTRACT(MINUTE FROM st.shift_start),

            CASE 
                WHEN st.shift_end <= st.shift_start THEN
                    EXTRACT(HOUR FROM st.shift_end) * 60 
                  + EXTRACT(MINUTE FROM st.shift_end) + 1440
                ELSE
                    EXTRACT(HOUR FROM st.shift_end) * 60 
                  + EXTRACT(MINUTE FROM st.shift_end)
            END
        FROM shift_timing st
        JOIN tenant_shift ts ON ts.id = st.tenant_shift_id
        WHERE ts.tenant_id = (
            SELECT tenant_id 
            FROM tenant_shift 
            WHERE id = NEW.tenant_shift_id
        )
        AND st.weekday = NEW.weekday
        AND st.id <> COALESCE(NEW.id, 0)
    LOOP
        IF new_start < existing_end AND new_end > existing_start THEN
            RAISE EXCEPTION 'Shift overlap detected for weekday %', NEW.weekday;
        END IF;
    END LOOP;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =========================================
-- TRIGGERS (IMPORTANT FIX)
-- =========================================

-- ✅ tenant_shift trigger (THIS WAS MISSING FIX)
DROP TRIGGER IF EXISTS trg_update_tenant_shift_updated_at ON tenant_shift;

CREATE TRIGGER trg_update_tenant_shift_updated_at
BEFORE UPDATE ON tenant_shift
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- ✅ overlap trigger
DROP TRIGGER IF EXISTS trg_check_shift_overlap ON shift_timing;

CREATE TRIGGER trg_check_shift_overlap
BEFORE INSERT OR UPDATE ON shift_timing
FOR EACH ROW
EXECUTE FUNCTION check_shift_overlap();

-- ✅ updated_at trigger
DROP TRIGGER IF EXISTS trg_update_shift_timing_updated_at ON shift_timing;

CREATE TRIGGER trg_update_shift_timing_updated_at
BEFORE UPDATE ON shift_timing
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
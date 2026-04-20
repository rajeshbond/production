CREATE TABLE IF NOT EXISTS shift_hour_slot (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    shift_timing_id BIGINT NOT NULL,
    slot_start TIME NOT NULL,
    slot_end TIME NOT NULL,
    slot_index INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_slot_shift FOREIGN KEY (shift_timing_id) REFERENCES shift_timing (id) ON DELETE CASCADE,
    CONSTRAINT uix_slot_unique UNIQUE (
        tenant_id,
        shift_timing_id,
        slot_start
    )
);
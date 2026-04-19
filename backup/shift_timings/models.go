package shifttiming

import "time"

type ShiftTiming struct {
	ID            int64  `db:"id"`
	TenantShiftID int64  `db:"tenant_shift_id"`
	ShiftStart    string `db:"shift_start"` // TIME → use string "HH:MM"
	ShiftEnd      string `db:"shift_end"`
	Weekday       int    `db:"weekday"`

	CreatedBy *int64 `db:"created_by"`
	UpdatedBy *int64 `db:"updated_by"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

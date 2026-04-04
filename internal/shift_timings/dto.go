package shifttiming

import "time"

type ShiftTimingDTO struct {
	ShiftStart string `json:"shift_start" validate:"required"` // "09:00"
	ShiftEnd   string `json:"shift_end" validate:"required"`
	Weekday    int    `json:"weekday" validate:"required"` // 0–6
}

type CreateShiftTimingRequest struct {
	TenantShiftID int64            `json:"tenant_shift_id" validate:"required"`
	Timings       []ShiftTimingDTO `json:"timings" validate:"required,dive"`
}

type ShiftTimingResponse struct {
	ID            int64      `json:"id"`
	TenantShiftID int64      `json:"tenant_shift_id"`
	ShiftStart    string     `json:"shift_start"`
	ShiftEnd      string     `json:"shift_end"`
	Weekday       int        `json:"weekday"`
	CreatedBy     *int64     `json:"created_by"`
	UpdatedBy     *int64     `json:"updated_by"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

package shift

import "time"

//
// ===============================
// COMMON TIMING DTO
// ===============================
// (Single source of truth for timings)
//
type TimingDTO struct {
	ShiftStart string `json:"shift_start" validate:"required"`         // HH:MM
	ShiftEnd   string `json:"shift_end" validate:"required"`           // HH:MM
	Weekday    int    `json:"weekday" validate:"required,min=1,max=7"` // 1–7
}

//
// ===============================
// CREATE (Single Shift)
// ===============================
//
type CreateShiftTimingRequest struct {
	TenantShiftID int64       `json:"tenant_shift_id" validate:"required"`
	Timings       []TimingDTO `json:"timings" validate:"required,dive"`
}

//
// ===============================
// REPLACE (Delete + Insert)
// ===============================
//
type ReplaceShiftTimingRequest struct {
	TenantShiftID int64       `json:"tenant_shift_id" validate:"required"`
	Timings       []TimingDTO `json:"timings" validate:"required,dive"`
}

//
// ===============================
// BULK (Multi Shift + Tenant Code)
// ===============================
//
type BulkShiftRequest struct {
	TenantCode string      `json:"tenant_code" validate:"required"`
	ShiftName  string      `json:"shift_name" validate:"required"`
	Timings    []TimingDTO `json:"timings" validate:"required,dive"`
}

// 🔥 Root request (IMPORTANT)
type BulkCreateShiftRequest []BulkShiftRequest

//
// ===============================
// RESPONSE
// ===============================
//
type ShiftTimingResponse struct {
	ID            int64      `json:"id"`
	TenantShiftID int64      `json:"tenant_shift_id"`
	ShiftStart    string     `json:"shift_start"`
	ShiftEnd      string     `json:"shift_end"`
	Weekday       int        `json:"weekday"`
	CreatedBy     *int64     `json:"created_by,omitempty"`
	UpdatedBy     *int64     `json:"updated_by,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

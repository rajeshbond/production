package tenantshifts

import "time"

type ShiftDTO struct {
	ShiftName string `json:"shift_name" validate:"required"`
}

type CreateShiftRequest struct {
	TenantID int64      `json:"tenant_id" validate:"required"`
	Shifts   []ShiftDTO `json:"shifts" validate:"required,dive"`
}

type TenantShiftResponse struct {
	ID        int64     `json:"id"`
	TenantID  int64     `json:"tenant_id"`
	ShiftName string    `json:"shift_name"`
	CreatedBy *int64    `json:"created_by,omitempty"`
	UpdatedBy *int64    `json:"updated_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// type CreateTenantShiftRequest struct {
// 	TenantID  int64  `json:"tenant_id" validation:"required"`
// 	ShiftName string `json:"shift_name" validation:"required"`
// 	CreatedBy *int64 `json:"created_by"`
// }

// // CreataTenantShiftResponse struct

// type CreateTenantShiftResponse struct {
// 	ID        int64     `json:"id"`
// 	TenantID  int64     `json:"tenant_id"`
// 	ShiftName string    `json:"shift_name"`
// 	CreatedBy *int64    `json:"created_by"`
// 	UpdatedBy *int64    `json:"upated_by"`
// 	CreatedAt time.Time `json:"created_at"`
// 	UpdatedAt time.Time `json:"upated_at"`
// }

// type UpdateTEnantShiftRequest struct {
// 	ID        int64  `json:"id"`
// 	ShiftName string `json:"shift_name"`
// 	UpdatedBy *int64 `json:"updated_by"`
// }

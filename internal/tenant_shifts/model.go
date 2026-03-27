package tenantshifts

import "time"

type TenantShift struct {
	ID        int64     `db:"id"`
	TenantID  int64     `db:"tenant_id"`
	ShiftName string    `db:"shift_name"`
	CreatedBy *int64    `db:"created_at"`
	UpdatedBy *int64    `db:"updated_at"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

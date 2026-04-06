package defect

import "time"

type Defect struct {
	ID         int64  `db:"id"`
	TenantID   int64  `db:"tenant_id"`
	DefectName string `db:"defect_name"`

	CreatedBy int64 `db:"created_by"`
	UpdatedBy int64 `db:"updated_by"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	IsDeleted bool       `db:"is_deleted"`
	DeletedAt *time.Time `db:"deleted_at"`
	DeletedBy *int64     `db:"deleted_by"`
}

package resourcemaster

import "time"

type Resource struct {
	ID             int64  `db:"id"`
	TenantID       int64  `db:"tenant_id"`
	ResourceName   string `db:"resource_name"`
	ResourceTypeID int64  `db:"resource_type_id"`

	CreatedBy *int64 `db:"created_by"`
	UpdatedBy *int64 `db:"updated_by"`

	IsDeleted bool       `db:"is_deleted"`
	DeletedBy *int64     `db:"deleted_by"`
	DeletedAt *time.Time `db:"deleted_at"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

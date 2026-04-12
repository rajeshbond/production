package resourcetype

import "time"

type ResourceType struct {
	ID          int64  `db:"id"`
	TenantID    int64  `db:"tenant_id"`
	TypeName    string `db:"type_name"`
	Description string `db:"description"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	CreatedBy int64 `db:"created_by"`
	UpdatedBy int64 `db:"updated_by"`

	IsDeleted bool  `db:"is_deleted"`
	DeletedBy int64 `db:"deleted_by"`
}

package operationmode

import "time"

type OperationMode struct {
	ID                int64  `db:"id"`
	TenantID          int64  `db:"tenant_id"`
	OperationModeName string `db:"operation_mode_name"`

	IsActive  bool `db:"is_active"`
	IsDeleted bool `db:"is_deleted"`

	CreatedBy *int64 `db:"created_by"`
	UpdatedBy *int64 `db:"updated_by"`
	DeletedBy *int64 `db:"deleted_by"`

	DeletedAt *time.Time `db:"deleted_at"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

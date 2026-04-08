package operations

import "time"

type Operations struct {
	ID            int64     `db:"id"`
	TenantID      int64     `db:"tenant_id"`
	OperationName string    `db:"operation_name"`
	IsDeleted     bool      `db:"is_deleted"`
	DeletedBy     *int64    `db:"deleted_by"`
	CreatedBy     *int64    `db:"created_by"`
	UpdatedBy     *int64    `db:"updated_by"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

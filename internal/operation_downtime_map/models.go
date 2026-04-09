package operationdowntimemap

import "time"

type OperationDowntimeMap struct {
	ID          int64     `db:"id"`
	TenantID    int64     `db:"tenant_id"`
	OperationID int64     `db:"operation_id"`
	DowntimeID  int64     `db:"downtime_id"`
	CreatedBy   *int64    `db:"created_by"`
	UpdatedBy   *int64    `db:"updated_by"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

package operationdefectmap

import "time"

type Operationdefectmap struct {
	ID          int64     `db:"id"`
	TenantID    int64     `db:"tenant_id"`
	OperationID int64     `db:"operation_id"`
	DefectID    int64     `db:"defect_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

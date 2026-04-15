package productiontarget

import "time"

type ProductionTarget struct {
	ID                 int64  `db:"id"`
	TenantID           int64  `db:"tenant_id"`
	ProductID          int64  `db:"product_id"`
	OperationID        int64  `db:"operation_id"`
	MachineID          int64  `db:"machine_id"`
	ProcessType        string `db:"process_type"`
	TargetPerHour      int    `db:"target_per_hour"`
	ExpectedEfficiency *int   `db:"expected_efficiency"`
	ResourceSignature  string `db:"resource_signature"`

	CreatedBy *int64 `db:"created_by"`
	UpdatedBy *int64 `db:"updated_by"`

	IsDeleted *bool      `db:"is_deleted"`
	DeletedBy *int64     `db:"deleted_by"`
	DeletedAt *time.Time `db:"deleted_at"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type ProductionTargetResource struct {
	ID                 int64     `db:"id"`
	TenantID           int64     `db:"tenant_id"`
	ProductionTargetID int64     `db:"production_target_id"`
	ResourceID         int64     `db:"resource_id"`
	CreatedAt          time.Time `db:"created_at"`
}

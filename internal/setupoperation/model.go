package setupoperation

import "time"

type OperationSetup struct {
	ID           int64      `db:"id"`
	TenantID     int64      `db:"tenant_id"`
	PosID        int64      `db:"pos_id"`
	MachineID    int64      `db:"machine_id"`
	SetupName    *string    `db:"setup_name"`
	TargetQty    int        `db:"target_qty"`
	CycleTimeSec *int       `db:"cycle_time_sec"`
	SetupTimeMin *int       `db:"setup_time_min"`
	IsDeleted    bool       `db:"is_deleted"`
	DeletedAt    *time.Time `db:"deleted_at"`

	CreatedBy *int64 `db:"created_by"`
	UpdatedBy *int64 `db:"updated_by"`
	DeletedBy *int64 `db:"deleted_by"`
}

type OperationSetupResource struct {
	ID         int64     `db:"id"`
	TenantID   int64     `db:"tenant_id"`
	SetupID    int64     `db:"setup_id"`
	ResourceID int64     `db:"resource_id"`
	Quantity   int       `db:"quantity"`
	CreatedAt  time.Time `db:"created_at"`
}

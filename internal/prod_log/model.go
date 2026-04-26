package prodlog

import "time"

type ProductionLog struct {
	ID          int64 `db:"id"`
	TenantID    int64 `db:"tenant_id"`
	SetupID     int64 `db:"setup_id"`
	MachineID   int64 `db:"machine_id"`
	ProductID   int64 `db:"product_id"`
	OperationID int64 `db:"operation_id"`

	ProductionDate time.Time `db:"production_date"`

	ShiftID         int64 `db:"shift_name_id"`
	ShiftTimingID   int64 `db:"shift_timing_id"`
	ShiftHourSlotID int64 `db:"shift_hour_slot_id"`
	SlotIndex       int   `db:"slot_index"`

	TargetQty   int `db:"target_qty"`
	ActualQty   int `db:"actual_qty"`
	OkQty       int `db:"ok_qty"`
	RejectedQty int `db:"rejected_qty"`
	Scrap       int `db:"scrap"`

	Remarks *string `db:"remarks"`

	CreatedBy *int64 `db:"created_by"`
	UpdatedBy *int64 `db:"updated_by"`
	DeletedBy *int64 `db:"deleted_by"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`

	IsDeleted bool `db:"is_deleted"`
}

type ProductionDefect struct {
	ID              int64 `db:"id"`
	TenantID        int64 `db:"tenant_id"`
	ProductionLogID int64 `db:"production_log_id"`
	DefectID        int64 `db:"defect_id"`
	Qty             int   `db:"qty"`
}

// Production Downtime
type ProductionDowntime struct {
	ID              int64 `db:"id"`
	TenantID        int64 `db:"tenant_id"`
	ProductionLogID int64 `db:"production_log_id"`
	DowntimeID      int64 `db:"downtime_id"`
	Minutes         int   `db:"downtime_minutes"`
}

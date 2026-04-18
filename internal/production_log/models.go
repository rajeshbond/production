package productionlog

import "time"

type ProductionLog struct {
	ID              int64      `db:"id"`
	TenantID        int64      `db:"tenant_id"`
	MachineID       int64      `db:"machine_id"`
	ProductID       int64      `db:"product_id"`
	OperationID     int64      `db:"operation_id"`
	ProductionDate  time.Time  `db:"production_date"`
	ShiftHourSlotID int64      `db:"shift_hour_slot_id"`
	SlotIndex       int        `db:"slot_index"`
	Quantity        int        `db:"quantity"`
	Remarks         *string    `db:"remarks"`
	IsDeleted       bool       `db:"is_deleted"`
	CreatedBy       *int64     `db:"created_by"`
	UpdatedBy       *int64     `db:"updated_by"`
	DeletedBy       *int64     `db:"deleted_by"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
}

type ProductionLogResource struct {
	ResourceID int64 `db:"resource_id"`
}

type ProductionLogDefect struct {
	DefectID int64 `db:"defect_id"`
	Quantity int   `db:"quantity"`
}

type ProductionLogDowntime struct {
	DowntimeID int64 `db:"downtime_id"`
	Minutes    int   `db:"minutes"`
}

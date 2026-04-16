package productionlog

type ProductionLogDTO struct {
	ID          int64 `db:"id"`
	TenantID    int64 `db:"tenant_id"`
	MachineID   int64 `db:"machine_id"`
	ProductID   int64 `db:"product_id"`
	OperationID int64 `db:"operation_id"`

	SlotID    int64 `db:"shift_hour_slot_id"`
	SlotIndex int   `db:"slot_index"` // ✅ NEW

	Quantity int    `db:"quantity"`
	Remarks  string `db:"remarks"`
}

type DefectInput struct {
	DefectID int64 `json:"defect_id" validate:"required"`
	Quantity int   `json:"quantity" validate:"required"`
}

type DowntimeInput struct {
	DowntimeID int64 `json:"downtime_id" validate:"required"`
	Minutes    int   `json:"minutes" validate:"required"`
}

type CreateProductionLogRequest struct {
	MachineID   int64 `json:"machine_id" validate:"required"`
	ProductID   int64 `json:"product_id" validate:"required"`
	OperationID int64 `json:"operation_id" validate:"required"`
	ShiftSlotID int64 `json:"shift_slot_id" validate:"required"`
	// SlotID    int64	`json:"shift_slot_id" validate:"required"`

	Quantity int `json:"quantity" validate:"required"`

	ResourceIDs []int64         `json:"resource_ids"`
	Defects     []DefectInput   `json:"defects"`
	Downtimes   []DowntimeInput `json:"downtimes"`

	Remarks string `json:"remarks"`
}

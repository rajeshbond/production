package productionlog

// ================= CREATE =================

type CreateProductionLogRequest struct {
	MachineID      int64   `json:"machine_id" validate:"required"`
	ProductID      int64   `json:"product_id" validate:"required"`
	OperationID    int64   `json:"operation_id" validate:"required"`
	SlotID         int64   `json:"shift_hour_slot_id" validate:"required"`
	ProductionDate string  `json:"production_date" validate:"required,datetime=2006-01-02"`
	Quantity       int     `json:"quantity" validate:"required,min=0"`
	Remarks        *string `json:"remarks"`

	ResourceIDs []int64         `json:"resource_ids"`
	Defects     []DefectInput   `json:"defects"`
	Downtimes   []DowntimeInput `json:"downtimes"`
}

// ================= UPDATE =================

type UpdateProductionLogRequest struct {
	ID int64 `json:"id" validate:"required"`

	MachineID   int64 `json:"machine_id" validate:"required"`
	ProductID   int64 `json:"product_id" validate:"required"`
	OperationID int64 `json:"operation_id" validate:"required"`
	SlotID      int64 `json:"shift_hour_slot_id" validate:"required"`

	Quantity int     `json:"quantity" validate:"required,min=0"`
	Remarks  *string `json:"remarks"`

	ResourceIDs []int64         `json:"resource_ids"`
	Defects     []DefectInput   `json:"defects"`
	Downtimes   []DowntimeInput `json:"downtimes"`
}

// ================= CHILD DTO =================

type DefectInput struct {
	DefectID int64 `json:"defect_id" validate:"required"`
	Quantity int   `json:"quantity" validate:"required,min=0"`
}

type DowntimeInput struct {
	DowntimeID int64 `json:"downtime_id" validate:"required"`
	Minutes    int   `json:"minutes" validate:"required,min=0"`
}

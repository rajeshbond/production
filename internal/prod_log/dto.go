package prodlog

type DefectInput struct {
	DefectID int64 `json:"defect_id"`
	Qty      int   `json:"qty"`
}

type DowntimeInput struct {
	DowntimeID int64 `json:"downtime_id"`
	Minutes    int   `json:"minutes"`
}

type CreateRequest struct {
	SetupID     int64 `json:"setup_id"`
	MachineID   int64 `json:"machine_id"`
	ProductID   int64 `json:"product_id"`
	OperationID int64 `json:"operation_id"`

	ProductionDate string `json:"production_date"`

	ShiftID         int64 `json:"shift_id"`
	ShiftTimingID   int64 `json:"shift_timing_id"`
	ShiftHourSlotID int64 `json:"shift_hour_slot_id"`
	SlotIndex       int   `json:"slot_index"`

	TargetQty   int `json:"target_qty"`
	ActualQty   int `json:"actual_qty"`
	OkQty       int `json:"ok_qty"`
	RejectedQty int `json:"rejected_qty"`
	Scrap       int `json:"scrap"`

	Remarks *string `json:"remarks"`

	Defects  []DefectInput   `json:"defects"`
	Downtime []DowntimeInput `json:"downtime"`
}

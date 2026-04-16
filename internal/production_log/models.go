package productionlog

type ProductionLog struct {
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

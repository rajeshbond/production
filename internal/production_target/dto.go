package productiontarget

type CreateProductionTargetRequest struct {
	ProductID          int64   `json:"product_id" validate:"required"`
	OperationID        int64   `json:"operation_id" validate:"required"`
	MachineID          int64   `json:"machine_id" validate:"required"`
	ProcessType        string  `json:"process_type" validate:"required"`
	TargetPerHour      int     `json:"target_per_hour" validate:"required"`
	ExpectedEfficiency *int    `json:"expected_efficiency"`
	ResourceIDs        []int64 `json:"resource_ids"` // 🔥 important
}

type UpdateProductionTargetRequest struct {
	ID                 int64   `json:"id" validate:"required"`
	ProductID          int64   `json:"product_id"`
	OperationID        int64   `json:"operation_id"`
	MachineID          int64   `json:"machine_id"`
	ProcessType        string  `json:"process_type"`
	TargetPerHour      int     `json:"target_per_hour"`
	ExpectedEfficiency *int    `json:"expected_efficiency"`
	ResourceIDs        []int64 `json:"resource_ids"`
}

type ProductionTargetResponseDTO struct {
	ID                 int64   `json:"id"`
	ProductID          int64   `json:"product_id"`
	OperationID        int64   `json:"operation_id"`
	MachineID          int64   `json:"machine_id"`
	ProcessType        string  `json:"process_type"`
	TargetPerHour      int     `json:"target_per_hour"`
	ExpectedEfficiency *int    `json:"expected_efficiency"`
	ResourceIDs        []int64 `json:"resource_ids"`
}

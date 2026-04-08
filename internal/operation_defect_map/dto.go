package operationdefectmap

type OperationDefectCreateRequest struct {
	OperationName string   `json:"operation_name" validate:"required,min = 2"`
	DefectNames   []string `json:"defect_names" validate:"required, min =1,dive,required" `
}

type OperationDefectCreateResponse struct {
	OperationID int64    `json:"operation_id"`
	Inserted    []string `json:"inserted"`
	Skipped     []string `json:"skipped"`
}

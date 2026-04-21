package operationmode

type CreateOperationModeRequest struct {
	OperationModeName string `json:"operation_mode_name" validate:"required"`
}

type UpdateOperationModeRequest struct {
	ID                int64  `json:"id" validate:"required"`
	OperationModeName string `json:"operation_mode_name" validate:"required"`
	IsActive          *bool  `json:"is_active"`
}

type OperationModeResponse struct {
	ID                int64  `json:"id"`
	OperationModeName string `json:"operation_mode_name"`
	IsActive          bool   `json:"is_active"`
}

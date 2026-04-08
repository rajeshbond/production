package operations

type CreateOperationsRequest struct {
	OperationName string `json:"operation_name" validate:"required,min=2,max=100"`
}

type BulkCreateOperationRequest struct {
	Operations []CreateOperationsRequest `json:"operations" validate:"required,dive"`
}

type OperationResponse struct {
	ID            int64  `json:"id"`
	OperationName string `json:"operation_name"`
}

type BulkOperationResult struct {
	Inserted []OperationResponse
	Skipped  []string
}

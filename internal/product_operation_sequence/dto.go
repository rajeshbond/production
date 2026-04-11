package productoperationsequence

type CreateProductOperationRequest struct {
	ProductID   int64 `json:"product_id" validate:"required"`
	OperationID int64 `json:"operation_id" validate:"required"`
}

type ProductOperationResponse struct {
	ID          int64 `json:"id"`
	ProductID   int64 `json:"product_id"`
	OperationID int64 `json:"operation_id"`
	SequenceNo  int   `json:"sequence_no"`
}

type DeleteSequenceFromRequest struct {
	ProductID int64 `json:"product_id"`
	FromSeqNo int   `json:"from_sequence_no"`
}

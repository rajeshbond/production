package operationdowntimemap

type OperationDowntimeCreateRequest struct {
	OperationName string   `json:"operation_name" validate:"required,min=2"`
	DownTimeName  []string `json:"downtime_name" validate:"required,min=1,dive,required,min=1"`
}

type OperationDowntimeCreateResponse struct {
	OperationID int64    `json:"operation_id"`
	Inserted    []string `json:"inserted"`
	Skipped     []string `json:"skipped"`
}

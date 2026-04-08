package downtime

type CreateDowntimeRequest struct {
	DowntimeName string `json:"downtime_name" validate:"required,min=2,max=100"`
}

type BulkCreateDownTimeRequest struct {
	DownTime []CreateDowntimeRequest `json:"downtime" validate:"required,dive"`
}

type DownTimeResponse struct {
	ID           int64  `json:"id"`
	DowntimeName string `json:"downtime_name"`
}

type BulkDownTimeResult struct {
	Inserted []DownTimeResponse
	Skipped  []string
}

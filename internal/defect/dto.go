package defect

type CreateDefectRequest struct {
	DefectName string `json:"defect_name" validate:"required,min=2,max=100"`
}

type BulkCreateDefectRequest struct {
	Defects []CreateDefectRequest `json:"defects" validate:"required,dive"`
}

type DefectResponse struct {
	ID         int64  `json:"id"`
	DefectName string `json:"defect_name"`
}

type BulkDefectResult struct {
	Inserted []DefectResponse
	Skipped  []string
}

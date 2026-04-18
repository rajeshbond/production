package mold

import "encoding/json"

type CreateMoldRequest struct {
	MoldNo       string          `json:"mold_no" validate:"required"`
	Description  *string         `json:"description"`
	Cavities     int             `json:"cavities" validate:"required,gt=0"`
	TargetShots  int             `json:"target_shots" validate:"required,gte=0"`
	SpecialNotes json.RawMessage `json:"special_notes"`
}

type UpdateMoldRequest struct {
	ID           int64           `json:"id" validate:"required"`
	MoldNo       string          `json:"mold_no" validate:"required"`
	Description  *string         `json:"description"`
	Cavities     int             `json:"cavities" validate:"required,gt=0"`
	TargetShots  int             `json:"target_shots" validate:"required,gte=0"`
	SpecialNotes json.RawMessage `json:"special_notes"`
}

type DeleteMoldRequest struct {
	ID int64 `json:"id" validate:"required"`
}

type BulkCreateMoldRequest struct {
	Molds []CreateMoldRequest `json:"molds" validate:"required,dive"`
}

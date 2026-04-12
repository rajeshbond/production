package resourcetype

type CreateResourceTypeRequest struct {
	TypeName    string `json:"type_name" validate:"required"`
	Description string `json:"description"`
}

type UpdateResourceTypeRequest struct {
	ID          int64  `json:"id" validate:"required"`
	TypeName    string `json:"type_name" validate:"required"`
	Description string `json:"description"`
}

type DeleteResourceTypeRequest struct {
	ID int64 `json:"id" validate:"required"`
}

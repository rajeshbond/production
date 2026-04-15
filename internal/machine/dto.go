package machine

type CreateMachineRequest struct {
	MachineCode string  `json:"machine_code" validate:"required"`
	MachineName string  `json:"machine_name" validate:"required"`
	Description *string `json:"description"`
	Capacity    *string `json:"capacity"`
}

type UpdateMachineRequest struct {
	ID          int64   `json:"id" validate:"required"`
	MachineCode string  `json:"machine_code" validate:"required"`
	MachineName string  `json:"machine_name" validate:"required"`
	Description *string `json:"description"`
	Capacity    *string `json:"capacity"`
}

type MachineResponse struct {
	ID          int64   `json:"id"`
	MachineCode string  `json:"machine_code"`
	MachineName string  `json:"machine_name"`
	Description *string `json:"description,omitempty"`
	Capacity    *string `json:"capacity,omitempty"`
}

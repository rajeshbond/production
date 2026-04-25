package setupoperation

type ResourceItem struct {
	ResourceID int64 `json:"resource_id"`
	Quantity   int   `json:"quantity"`
}

type CreateSetupRequest struct {
	PosID        int64          `json:"pos_id"`
	MachineID    int64          `json:"machine_id"`
	TargetQty    int            `json:"target_qty"`
	SetupName    *string        `json:"setup_name"`
	CycleTimeSec *int           `json:"cycle_time_sec"`
	SetupTimeMin *int           `json:"setup_time_min"`
	Resources    []ResourceItem `json:"resources"`
	// Resources []ResourceItem `json:"resources"`
}

type UpdateSetupRequest struct {
	ID        int64          `json:"id"`
	MachineID int64          `json:"machine_id"`
	TargetQty int            `json:"target_qty"`
	SetupName *string        `json:"setup_name"`
	Resources []ResourceItem `json:"resources"`
}

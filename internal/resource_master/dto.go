package resourcemaster

type CreateResourceRequest struct {
	ResourceName   string `json:"resource_name"`
	ResourceTypeID int64  `json:"resource_type_id"`
}

type UpdateResourceRequest struct {
	ID             int64  `json:"id"`
	ResourceName   string `json:"resource_name"`
	ResourceTypeID int64  `json:"resource_type_id"`
}

type ResourceResponse struct {
	ID             int64  `json:"id"`
	ResourceName   string `json:"resource_name"`
	ResourceTypeID int64  `json:"resource_type_id"`
}

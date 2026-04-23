package resource

type CreateResourceRequest struct {
	ResourceSubID int64   `json:"resource_sub_id"`
	ResourceCode  string  `json:"resource_code"`
	ResourceName  *string `json:"resource_name"`
	ResourceType  string  `json:"resource_type"`
	Description   *string `json:"description"`
}

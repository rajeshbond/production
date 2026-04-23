package resource

type CreateResourceRequest struct {
	ResourceCode string  `json:"resource_code"`
	ResourceName *string `json:"resource_name"`
	ResourceType string  `json:"resource_type"`
	Description  *string `json:"description"`

	MoldID    *int64 `json:"mold_id"`
	FixtureID *int64 `json:"fixture_id"`
	ToolID    *int64 `json:"tool_id"`
}

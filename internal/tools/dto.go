package tools

type CreateToolRequest struct {
	Type     string `json:"type"`
	ToolCode string `json:"tool_code"`
	ToolName string `json:"tool_name"`

	Description *string `json:"description"`
	ToolType    *string `json:"tool_type"`
	Unit        *string `json:"unit"`

	Cost       float64 `json:"cost"`
	LifeCycles int64   `json:"life_cycles"`

	CreatedBy int64 `json:"created_by"`
}

type UpdateToolRequest struct {
	ToolName    *string `json:"tool_name"`
	Description *string `json:"description"`
	ToolType    *string `json:"tool_type"`
	Unit        *string `json:"unit"`

	Cost       *float64 `json:"cost"`
	LifeCycles *int64   `json:"life_cycles"`
	IsActive   *bool    `json:"is_active"`

	UpdatedBy int64 `json:"updated_by"`
}

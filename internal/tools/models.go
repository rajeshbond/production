package tools

type Tool struct {
	ID       int64  `json:"id"`
	TenantID int64  `json:"tenant_id"`
	Type     string `json:"type"`

	ToolCode    string  `json:"tool_code"`
	ToolName    string  `json:"tool_name"`
	Description *string `json:"description"`

	ToolType *string `json:"tool_type"`
	Unit     *string `json:"unit"`

	Cost       float64 `json:"cost"`
	LifeCycles int64   `json:"life_cycles"`

	IsActive  bool `json:"is_active"`
	IsDeleted bool `json:"is_deleted"`

	CreatedBy *int64 `json:"created_by"`
	UpdatedBy *int64 `json:"updated_by"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

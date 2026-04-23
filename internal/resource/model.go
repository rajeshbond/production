package resource

import "time"

type Resource struct {
	ID           int64   `db:"id"`
	TenantID     int64   `db:"tenant_id"`
	ResourceCode string  `db:"resource_code"`
	ResourceName *string `db:"resource_name"`
	ResourceType string  `db:"resource_type"`
	Description  *string `db:"description"`

	MoldID    *int64 `db:"mold_id"`
	FixtureID *int64 `db:"fixture_id"`
	ToolID    *int64 `db:"tool_id"`

	IsActive bool `db:"is_active"`

	CreatedBy *int64 `db:"created_by"`
	UpdatedBy *int64 `db:"updated_by"`
	DeletedBy *int64 `db:"deleted_by"`

	IsDeleted bool       `db:"is_deleted"`
	DeletedAt *time.Time `db:"deleted_at"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

package tenant

import "time"

// Create Tenant DTO

type CreateTenantRequied struct {
	TenantName string `json:"tenant_name" validate:"required"`
	TenantCode string `json:"tenant_code" validate:"required"`
	Address    string `json:"address" validate:"required"`
}

// type CreateTenantDTO
type CreateTenantDTO struct {
	TenantName string `json:"tenant_name"`
	TenantCode string `json:"tenant_code"`
	Address    string `json:"address"`
	CreatedBy  int64  `json:"created_by"`
}

// Update Tenant DTO

type UpdateTenantRequest struct {
	TenantName *string `json:"tenant_name,omitempty"`
	Address    *string `json:"address,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
}

type TenantResponse struct {
	ID         int64  `json:"id"`
	TenantName string `json:"tenant_name"`
	TenantCode string `json:"tenant_code"`
	Address    string `json:"address"`
	IsVerified bool   `json:"is_verified"`
	IsActive   bool   `json:"is_active"`
	CreatedBy  *int64 `json:"created_by"`
	UpdatedBy  *int64 `json:"updated_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

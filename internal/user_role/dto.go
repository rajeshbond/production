package userrole

import "time"

type CreateUserRoleDTO struct {
	UserRole  string `json:"user_role"`
	CreatedBy *int64 `json:"created_by,omitempty"`
	UpdatedBy *int64 `json:"updated_at,omitempty"`
}

type UserRoleResponseDTO struct {
	ID        int64     `json:"id"`
	UserRole  string    `json:"user_role"`
	CreatedBy *int64    `json:"created_by,omitempty"`
	UpdatedBy *int64    `json:"updated_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRoleDBDTO struct {
	UserRole  string
	CreatedBy *int64
}

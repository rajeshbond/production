package user

import "time"

// User Response DTO

type CreateUserRequest struct {
	EmployeeID string `json:"employee_id" validate:"required"`
	UserName   string `json:"user_name" validate:"required"`
	TenantName string `json:"tenant_id" validate:"required"`
	Phone      string `json:"phone" validate:"required"`
	Role       string `json:"role" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=6"`
}

// User DTO

type UserDTO struct {
	TenantID   *int64  `json:"tenant_id,omitempty"`
	RoleID     *int64  `json:"role_id,omitempty"`
	EmployeeID *string `json:"employee_id,omitempty"`
	UserName   *string `json:"user_name,omitempty"`
	Phone      *string `json:"phone,omitempty"`
	Email      *string `json:"email,omitempty"`
	Password   *string `json:"password,omitempty"`
	IsVerified *bool   `json:"is_verified,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
	CreatedBy  *int64  `json:"created_by,omitempty"`
	UpdatedBy  *int64  `json:"updated_by,omitempty"`
}

type CreateUserRequestDTO struct {
	RoleID     int64  `json:"role_id"`
	EmployeeID string `json:"employee_id"`
	UserName   string `json:"user_name"`
	Phone      string `json:"phone,omitempty"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	CreatedBy  string `json:"CreatedBy"`
}

// Update User Request DTO

type UpdateUserDTO struct {
	RoleID     *int64  `json:"role_id,omitempty"`
	UserName   *string `json:"user_name,omitempty"`
	Phone      *string `json:"phone,omitempty"`
	Email      *string `json:"email,omitempty"`
	IsVerified *bool   `json:"is_verified,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
	UpdatedBy  *int64  `json:"updated_by,omitempty"`
}

// User DTO

// User ResponseDTO
type UserResponseDTO struct {
	ID         int64   `db:"id" json:"id"`
	TenantID   int64   `db:"tenant_id" json:"tenant_id"`
	RoleID     int64   `db:"role_id" json:"role_id"`
	EmployeeID string  `db:"employee_id" json:"employee_id"`
	UserName   string  `db:"user_name" json:"user_name"`
	Phone      *string `db:"phone" json:"phone,omitempty"`
	Email      *string `db:"email" json:"email,omitempty"`
	// Password   string    `db:"password" json:"-"`
	IsVerified bool      `db:"is_verified" json:"is_verified"`
	IsActive   bool      `db:"is_active" json:"is_active"`
	CreatedBy  *int64    `db:"created_by" json:"created_by,omitempty"`
	UpdatedBy  *int64    `db:"updated_by" json:"updated_by,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

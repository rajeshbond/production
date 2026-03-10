package users

import "time"

type UserCreateRequest struct {
	TenantID   int64   `json:"tenant_id" validate:"required"`
	RoleId     int64   `json:"role_id" validate:"required"`
	EmployeeID string  `json:"employee_id" validate:"required"`
	UserName   string  `json:"user_name" validate:"required"`
	Phone      *string `json:"phone" validate:"required"`
	Role       *string `json:"role" validate:"required"`
	Email      *string `json:"email"`
	Password   string  `json:"password" validate:"required,min=6"`
	// IsVerified *bool  `json:"is_verified,omitempty"`
	// IsActive   *bool  `json:"is_active,omitempty"`
	CreatedBy *int64 `json:"created_by,omitempty"`
	UpdatedBy *int64 `json:"updated_by,omitempty"`
}

// Response DTO

type UserResponse struct {
	ID         int64     `json:"id"`
	TenantID   int64     `json:"tenant_id"`
	RoleID     int64     `json:"role_id"`
	EmployeeID string    `json:"employee_id"`
	UserName   string    `json:"user_name"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"is_verified"`
	IsActive   bool      `json:"is_active"`
	CreatedBy  int64     `json:"created_by"`
	UpdatedBy  int64     `json:"updated_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

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

// Token Payload

type TokenPayload struct {
	TenantID int64  `json:"tenant_id"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	RoleID   int64  `json:"role_id"`
}

type LoginRequest struct {
	EmployeeID string `json:"employee_id" validation:"required,employeeid"`
	Password   string `json:"password" validate:"required,min=6,max=72"`
}

type LoginResponse struct {
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

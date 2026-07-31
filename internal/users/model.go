package users

import "time"

type User struct {
	ID       int64 `db:"id" json:"id"`
	TenantID int64 `db:"tenant_id" json:"tenant_id"`
	RoleID   int64 `db:"role_id" json:"role_id"`

	EmployeeID string  `db:"employee_id" json:"employee_id" `
	UserName   string  `db:"user_name" json:"user_name"`
	Phone      *string `db:"phone" json:"phone"`
	Email      *string `db:"email" json:"email"`

	Password string `db:"password" json:"password"`

	// Status
	IsVerified bool `db:"is_verified" json:"is_verified"`
	IsActive   bool `db:"is_active" json:"is_active"`

	// Soft Delete
	IsDeleted bool       `db:"is_deleted" json:"is_deleted"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
	DeletedBy *int64     `db:"deleted_by" json:"deleted_by"`

	// Audit
	CreatedBy *int64    `db:"created_by" json:"created_by"`
	UpdatedBy *int64    `db:"updated_by" json:"updated_by"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

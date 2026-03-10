package user

import "time"

type User struct {
	ID         int64   `db:"id"`
	TenantID   int64   `db:"tenant_id"`
	RoleID     int64   `db:"role_id"`
	EmployeeID string  `db:"employee_id"`
	UserName   string  `db:"user_name"`
	Phone      *string `db:"phone"`
	Email      *string `db:"email"`
	Password   string  `db:"password"`

	IsVerified bool `db:"is_verified"`
	IsActive   bool `db:"is_active"`

	CreatedBy *int64 `db:"created_by"`
	UpdatedBy *int64 `db:"updated_by"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

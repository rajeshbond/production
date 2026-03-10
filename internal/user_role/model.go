package userrole

import "time"

type UserRole struct {
	ID        int64     `db:"id"`
	UserRole  string    `db:"user_role"`
	CreatedBy *int64    `db:"created_by"`
	UpdatedBy *int64    `db:"updated_by"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

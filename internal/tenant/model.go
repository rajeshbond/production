package tenant

import "time"

type Tenant struct {
	ID         int64  `db:"id"`
	TenantName string `db:"tenant_name"`
	TenantCode string `db:"tenant_code"`
	Address    string `db:"address"`

	IsVerified bool `db:"is_verified"`
	IsActive   bool `db:"is_active"`

	CreatedBy *int64 `db:"created_by"`
	UpdatedBy *int64 `db:"updated_by"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

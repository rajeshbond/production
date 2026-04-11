package products

import "time"

type Product struct {
	ID          int64     `db:"id"`
	TenantID    int64     `db:"tenant_id"`
	ProductName string    `db:"product_name"`
	ProductNo   string    `db:"product_no"`
	CreatedBy   *int64    `db:"created_by"`
	UpdatedBy   *int64    `db:"updated_by"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

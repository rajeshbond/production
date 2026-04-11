package productoperationsequence

import "time"

type ProductOperationSequence struct {
	ID          int64     `db:"id"`
	TenantID    int64     `db:"tenant_id"`
	ProductID   int64     `db:"product_id"`
	OperationID int64     `db:"operation_id"`
	SequenceNo  int       `db:"sequence_no"`
	CreatedBy   *int64    `db:"created_by"`
	UpdatedBy   *int64    `db:"updared_by"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

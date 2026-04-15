package machine

import "time"

type Machine struct {
	ID           int64   `db:"id" json:"id"`
	TenantID     int64   `db:"tenant_id" json:"tenant_id"`
	MachineCode  string  `db:"machine_code" json:"machine_code"`
	MachineName  string  `db:"machine_name" json:"machine_name"`
	Description  *string `db:"description" json:"description,omitempty"`
	Capacity     *string `db:"capacity" json:"capacity,omitempty"`
	SpecialNotes []byte  `db:"special_notes" json:"special_notes,omitempty"`

	IsDeleted bool       `db:"is_deleted" json:"-"`
	DeletedBy *int64     `db:"deleted_by" json:"-"`
	DeletedAt *time.Time `db:"deleted_at" json:"-"`

	CreatedBy *int64    `db:"created_by" json:"created_by,omitempty"`
	UpdatedBy *int64    `db:"updated_by" json:"updated_by,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

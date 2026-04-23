package mould

import "time"

type Mold struct {
	ID           int64   `db:"id"`
	TenantID     int64   `db:"tenant_id"`
	Type         string  `db:"type"`
	MoldName     *string `db:"mold_name"`
	MoldNo       string  `db:"mold_no"`
	Description  *string `db:"description"`
	Cavities     int     `db:"cavities"`
	TargetShots  int64   `db:"target_shots"`
	SpecialNotes []byte  `db:"special_notes"`

	CreatedBy *int64 `db:"created_by"`
	UpdatedBy *int64 `db:"updated_by"`
	DeletedBy *int64 `db:"deleted_by"`

	IsDeleted bool       `db:"is_deleted"`
	DeletedAt *time.Time `db:"deleted_at"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

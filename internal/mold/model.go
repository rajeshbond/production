package mold

import "time"

type Mold struct {
	ID           int64      `db:"id" json:"id"`
	TenantID     int64      `db:"tenant_id" json:"tenant_id"`
	MoldNo       string     `db:"mold_no" json:"mold_no"`
	Description  *string    `db:"description" json:"description"`
	Cavities     int        `db:"cavities" json:"cavities"`
	TargetShots  int        `db:"target_shots" json:"target_shots"`
	SpecialNotes []byte     `db:"special_notes" json:"special_notes"` // JSONB
	CreatedBy    *int64     `db:"created_by" json:"created_by"`
	UpdatedBy    *int64     `db:"updated_by" json:"updated_by"`
	IsDeleted    bool       `db:"is_deleted" json:"is_deleted"`
	DeletedAt    *time.Time `db:"deleted_at" json:"deleted_at"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}

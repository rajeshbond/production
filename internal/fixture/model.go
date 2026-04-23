package fixture

import "time"

type Fixture struct {
	ID       int64  `db:"id"`
	TenantID int64  `db:"tenant_id"`
	Type     string `db:"type"`

	FixtureNo   string  `db:"fixture_no"`
	FixtureName *string `db:"fixture_name"`
	Description *string `db:"description"`

	Cavities  *int   `db:"cavities"`
	LifeShots *int64 `db:"life_shots"`

	FixtureType *string `db:"fixture_type"`
	Material    *string `db:"material"`

	SpecialNotes []byte `db:"special_notes"`

	CreatedBy *int64 `db:"created_by"`
	UpdatedBy *int64 `db:"updated_by"`
	DeletedBy *int64 `db:"deleted_by"`

	IsDeleted bool       `db:"is_deleted"`
	DeletedAt *time.Time `db:"deleted_at"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

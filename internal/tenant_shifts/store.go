package tenantshifts

import (
	"context"
	"database/sql"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateTenantShift(ctx context.Context, req *CreateTenantShiftRequest) (*TenantShift, error) {

	// Create tenant shift
	query := `
	INSERT INTO tenant_shift (
	tenant_id,
	shift_name,
	created_by
	)VALUES($1,$2,$3)
	RETURNING id,
	tenant_id,
	shift_name,
	created_by,
	updated_by,
	created_at,
	updated_at
	`
	var ts TenantShift

	err := s.db.QueryRowContext(ctx, query,
		req.TenantID,
		req.ShiftName,
		req.CreatedBy,
	).Scan(
		&ts.ID,
		&ts.TenantID,
		&ts.ShiftName,
		&ts.CreatedBy,
		&ts.UpdatedBy,
		&ts.CreatedAt,
		&ts.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &ts, nil

}

package tenantshifts

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateTenantShifts(ctx context.Context, tx *sql.Tx, tenantID int64, shiftName string, userID int64) (*TenantShift, error) {

	query := `
		INSERT INTO tenant_shift (tenant_id,shift_name,created_by,updated_by)
		VALUES($1,$2,$3,$4)
		ON CONFLICT (tenant_id, shift_name) DO NOTHING
		RETURNING id,
		tenant_id,
		shift_name,
		created_by,
		updated_by,
		created_at,
		updated_at
	`

	var ts TenantShift

	err := tx.QueryRowContext(
		ctx,
		query,
		tenantID,
		shiftName,
		userID,
		userID,
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
		if err == sql.ErrNoRows {
			return nil, nil // already exists
		}
		return nil, err
	}

	return &ts, nil

	// return ts, err

}

// 🔍 Get existing shifts from DB
func (s *Store) GetExistingShifts(ctx context.Context, tx *sql.Tx, tenantID int64, names []string) (map[string]bool, error) {

	query := `
		SELECT shift_name
		FROM tenant_shift
		WHERE tenant_id = $1
		AND shift_name = ANY($2)
	`

	rows, err := tx.QueryContext(ctx, query, tenantID, pq.Array(names))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	existing := make(map[string]bool)

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		existing[name] = true
	}

	return existing, nil
}

package shift

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Create(
	ctx context.Context,
	tx *sql.Tx,
	userID int64,
	req CreateShiftTimingRequest,
) error {

	for _, t := range req.Timings {
		if err := s.InsertShiftTiming(ctx, tx, req.TenantShiftID, t, userID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) InsertShiftTiming(
	ctx context.Context,
	tx *sql.Tx,
	shiftID int64,
	t TimingDTO,
	userID int64,
) error {

	query := `
	INSERT INTO shift_timing
	(
		tenant_shift_id,
		shift_start,
		shift_end,
		weekday,
		created_by,
		updated_by
	)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (tenant_shift_id, weekday, shift_start, shift_end)
	DO NOTHING;
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		shiftID,
		t.ShiftStart,
		t.ShiftEnd,
		t.Weekday,
		userID,
		userID,
	)

	if err != nil {
		// 🔥 Handle DB trigger error nicely
		if pqErr, ok := err.(*pq.Error); ok {
			return fmt.Errorf("shift timing error: %s", pqErr.Message)
		}
		return err
	}

	return nil
}

func (s *Store) GetTenantIDByCode(
	ctx context.Context,
	tx *sql.Tx,
	code string,
) (int64, error) {

	var id int64

	err := tx.QueryRowContext(ctx,
		`SELECT id FROM tenant WHERE tenant_code = $1`,
		code,
	).Scan(&id)

	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("tenant not found: %s", code)
	}

	return id, err
}

func (s *Store) UpsertTenantShift(
	ctx context.Context,
	tx *sql.Tx,
	tenantID int64,
	shiftName string,
	userID int64,
) (int64, error) {

	query := `
	INSERT INTO tenant_shift (tenant_id, shift_name, created_by, updated_by)
	VALUES ($1,$2,$3,$4)
	ON CONFLICT (tenant_id, shift_name)
	DO UPDATE SET updated_at = NOW()
	RETURNING id;
	`

	var id int64
	err := tx.QueryRowContext(ctx, query,
		tenantID,
		shiftName,
		userID,
		userID,
	).Scan(&id)

	return id, err
}

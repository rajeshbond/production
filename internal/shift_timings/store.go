package shifttiming

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rajesh_bond/production/database"
)

type Store struct {
	sqlDB *sql.DB
	pgx   *pgxpool.Pool
}

func NewStore(db *database.DB) *Store {
	return &Store{
		sqlDB: db.SQLDB,
		pgx:   db.PGX,
	}
}

// type Store struct {
// 	db *sql.DB
// }

// func NewStore(db *sql.DB) *Store {
// 	return &Store{db: db}
// }

func (s *Store) GetExisttingTimings(ctx context.Context, tx *sql.Tx, tenantShiftID int64) (map[string]bool, error) {

	query := `
		SELECT shift_start,shift_end,weekday
		FROM shift_timing
		WHERE tenant_shift_id = $1
	`

	rows, err := tx.QueryContext(ctx, query, tenantShiftID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	existing := make(map[string]bool)

	for rows.Next() {
		var start, end string
		var weekday int
		if err := rows.Scan(&start, &end, &weekday); err != nil {
			return nil, err
		}
		key := start + "-" + end + "-" + strconv.Itoa(weekday)
		existing[key] = true
	}

	return existing, nil

}

func (s *Store) InsertShiftTimingTx(ctx context.Context, tx *sql.Tx, req ShiftTimingDTO, tenantShiftID int64, userID int64) (*ShiftTiming, error) {

	query := `
		INSERT INTO shift_timing 
		(tenant_shift_id, shift_start, shift_end, weekday, created_by, updated_by)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT DO NOTHING
		RETURNING id, tenant_shift_id, shift_start, shift_end, weekday
	`

	var st ShiftTiming

	err := tx.QueryRowContext(
		ctx,
		query,
		tenantShiftID,
		req.ShiftStart,
		req.ShiftEnd,
		req.Weekday,
		userID,
		userID,
	).Scan(
		&st.ID,
		&st.TenantShiftID,
		&st.ShiftStart,
		&st.ShiftEnd,
		&st.Weekday,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &st, nil
}

func (s *Store) ValidateTenantShiftAccess(
	ctx context.Context,
	tx *sql.Tx,
	tenantShiftID int64,
	tenantCode string,
) error {

	query := `
	SELECT EXISTS (
		SELECT 1
		FROM tenant_shift ts
		JOIN tenant t ON t.id = ts.tenant_id
		WHERE ts.id = $1 AND t.tenant_code = $2
	)
	`

	var exists bool
	err := tx.QueryRowContext(ctx, query, tenantShiftID, tenantCode).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return ErrUnAuthorized
	}

	return nil
}

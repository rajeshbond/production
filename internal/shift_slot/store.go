package shiftslot

import (
	"context"
	"database/sql"
	"time"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetShiftTimingsByTenant(
	ctx context.Context,
	tenantID int64,
) ([]struct {
	ID    int64
	Start time.Time
	End   time.Time
}, error) {

	query := `
	SELECT 
		st.id,
		st.shift_start,
		st.shift_end
	FROM shift_timing st
	JOIN tenant_shift ts 
		ON st.tenant_shift_id = ts.id
	WHERE ts.tenant_id = $1
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []struct {
		ID    int64
		Start time.Time
		End   time.Time
	}

	for rows.Next() {
		var r struct {
			ID    int64
			Start time.Time
			End   time.Time
		}

		if err := rows.Scan(&r.ID, &r.Start, &r.End); err != nil {
			return nil, err
		}

		result = append(result, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
func (s *Store) SlotsExist(ctx context.Context, tx *sql.Tx, shiftTimingID int64) (bool, error) {

	var exists bool

	query := `
	SELECT EXISTS (
		SELECT 1 FROM shift_hour_slot
		WHERE shift_timing_id = $1
	)
	`

	err := tx.QueryRowContext(ctx, query, shiftTimingID).Scan(&exists)
	return exists, err
}

func (s *Store) InsertSlots(
	ctx context.Context,
	tx *sql.Tx,
	tenantID int64,
	shiftTimingID int64,
	slots []ShiftSlot,
) error {

	query := `
	INSERT INTO shift_hour_slot
	(tenant_id, shift_timing_id, slot_start, slot_end, slot_index)
	VALUES ($1,$2,$3,$4,$5)
	`

	for _, slot := range slots {
		_, err := tx.ExecContext(ctx, query,
			tenantID,
			shiftTimingID,
			slot.Start,
			slot.End,
			slot.Index,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) DeleteSlots(
	ctx context.Context,
	tx *sql.Tx,
	shiftTimingID int64,
) error {

	query := `DELETE FROM shift_hour_slot WHERE shift_timing_id = $1`
	_, err := tx.ExecContext(ctx, query, shiftTimingID)
	return err
}

// func (s *Store) GetShiftTimingsByTenant(
// 	ctx context.Context,
// 	tenantID int64,
// ) ([]struct {
// 	ID    int64
// 	Start time.Time
// 	End   time.Time
// }, error) {

// 	query := `
// 	SELECT id, shift_start, shift_end
// 	FROM shift_timing st
// 	JOIN tenant_shift ts ON st.tenant_shift_id = ts.id
// 	WHERE ts.tenant_id = $1
// 	`

// 	rows, err := s.db.QueryContext(ctx, query, tenantID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var result []struct {
// 		ID    int64
// 		Start time.Time
// 		End   time.Time
// 	}

// 	for rows.Next() {
// 		var r struct {
// 			ID    int64
// 			Start time.Time
// 			End   time.Time
// 		}

// 		if err := rows.Scan(&r.ID, &r.Start, &r.End); err != nil {
// 			return nil, err
// 		}

// 		result = append(result, r)
// 	}

// 	return result, nil
// }

// func (s *Store) SlotsExist(
// 	ctx context.Context,
// 	tx *sql.Tx,
// 	shiftTimingID int64,
// ) (bool, error) {

// 	var exists bool

// 	query := `
// 	SELECT EXISTS (
// 		SELECT 1 FROM shift_hour_slot
// 		WHERE shift_timing_id = $1
// 	)
// 	`

// 	err := tx.QueryRowContext(ctx, query, shiftTimingID).Scan(&exists)
// 	return exists, err
// }

// func (s *Service) GenerateAllShiftSlots(ctx context.Context, req GenerateSlotRequest) error {

// 	shiftTimings, err := s.Store.GetShiftTimingsByTenant(ctx, req.TenantID)
// 	if err != nil {
// 		return err
// 	}

// 	tx, err := s.Store.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return err
// 	}

// 	for _, st := range shiftTimings {

// 		// check already exists
// 		exists, err := s.Store.SlotsExist(ctx, tx, st.ID)
// 		if err != nil {
// 			tx.Rollback()
// 			return err
// 		}

// 		if exists {
// 			continue // ✅ skip existing
// 		}

// 		slots := GenerateSlots(st.Start, st.End)

// 		err = s.Store.InsertSlots(ctx, tx, req.TenantID, st.ID, slots)
// 		if err != nil {
// 			tx.Rollback()
// 			return err
// 		}
// 	}

// 	return tx.Commit()
// }

// func (s *Store) GetShiftTiming(ctx context.Context, shiftTimingID int64) (time.Time, time.Time, error) {

// 	var start, end time.Time

// 	query := `
// 	SELECT shift_start, shift_end
// 	FROM shift_timing
// 	WHERE id = $1
// 	`

// 	err := s.db.QueryRowContext(ctx, query, shiftTimingID).
// 		Scan(&start, &end)

// 	return start, end, err
// }

// func (s *Store) DeleteSlots(
// 	ctx context.Context,
// 	tx *sql.Tx,
// 	shiftTimingID int64,
// ) error {

// 	query := `DELETE FROM shift_hour_slot WHERE shift_timing_id = $1`
// 	_, err := tx.ExecContext(ctx, query, shiftTimingID)
// 	return err
// }

// func (s *Store) InsertSlots(
// 	ctx context.Context,
// 	tx *sql.Tx,
// 	tenantID int64,
// 	shiftTimingID int64,
// 	slots []ShiftSlot,
// ) error {

// 	query := `
// 	INSERT INTO shift_hour_slot
// 	(tenant_id, shift_timing_id, slot_start, slot_end, slot_index)
// 	VALUES ($1,$2,$3,$4,$5)
// 	`

// 	for _, slot := range slots {
// 		_, err := tx.ExecContext(ctx, query,
// 			tenantID,
// 			shiftTimingID,
// 			slot.Start,
// 			slot.End,
// 			slot.Index,
// 		)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

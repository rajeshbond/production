package productionlog

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateProductionLog(ctx context.Context, tx *sql.Tx, tenantID, userID int64, req CreateProductionLogRequest) (int64, error) {
	// Validate date format

	// Validate date format
	prodDate, err := time.Parse("2006-01-02", req.ProductionDate)
	if err != nil {
		return 0, fmt.Errorf("invalid production_date format (YYYY-MM-DD)")
	}

	// =========================================
	// STEP 1: GET SLOT INDEX (DON'T TRUST UI)
	// =========================================

	var slotIndex int

	slotIndexQuery := `
		SELECT slot_index 
		FROM shift_hourly_slot
		WHERE id = $1
	`
	err = tx.QueryRowContext(ctx, slotIndexQuery, req.SlotID).Scan(&slotIndex)
	if err != nil {
		return 0, fmt.Errorf("Invalid shift slot")
	}

	// =========================================
	// STEP 2: SEQUENCE VALIDATION (WITH LOCK 🔥)
	// =========================================

	var maxSlot sql.NullInt64

	seqValidationQuery := `
		SELECT MAX(slot_index)
		FROM production_log 
		WHERE tenant_id =$1
			AND machine_id = $2
			AND shift_hour_slot_id = $3
			AND is_deleted = FALSE
		FOR UPDATE
	`
	err = tx.QueryRowContext(ctx, seqValidationQuery,
		tenantID,
		req.MachineID,
		req.SlotID,
	).Scan(&maxSlot)

	if err != nil {
		return 0, err
	}

	if maxSlot.Valid {
		if int(maxSlot.Int64)+1 != slotIndex {
			return 0, fmt.Errorf("invalid slot sequence,expected slot_index %d", maxSlot.Int64+1)
		}
	} else {
		if slotIndex != 1 {
			return 0, fmt.Errorf("first entry must be slot_index 1")
		}
	}

	// =========================================
	// STEP 3: INSERT LOG
	// =========================================

	var id int64

	insertQuery := `
	INSERT INTO production_log (
			tenant_id, machine_id, product_id, operation_id,
			production_date, shift_hour_slot_id, slot_index,
			quantity, remarks, created_by, updated_by
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$10)
		RETURNING id
	`
	err = tx.QueryRowContext(ctx, insertQuery,
		tenantID,
		req.MachineID,
		req.ProductID,
		req.OperationID,
		prodDate,
		req.SlotID,
		slotIndex,
		req.Quantity,
		req.Remarks,
		userID,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	// =========================================
	// STEP 4: RESOURCE
	// =========================================

	if len(req.ResourceIDs) > 0 {
		resourceQuery := `
			INSERT INTO production_log_resource(production_log_id,rescorce_id)
			VALUES($1,$2)
		`
		for _, r := range req.ResourceIDs {
			if _, err := tx.ExecContext(ctx, resourceQuery, id, r); err != nil {
				return 0, err
			}
		}
	}

	// =========================================
	// STEP 5: DEFECT
	// =========================================

	if len(req.Defects) > 0 {
		defectQuery := `
			INSERT INTO production_log_defect (production_log_id, defect_id, quantity)
			VALUES ($1,$2,$3)
		`

		for _, d := range req.Defects {
			if _, err := tx.ExecContext(ctx, defectQuery, id, d.DefectID, d.Quantity); err != nil {
				return 0, err
			}
		}
	}

	// =========================================
	// STEP 6: DOWNTIME
	// =========================================

	if len(req.Downtimes) > 0 {
		dtQuery := `
			INSERT INTO production_log_downtime (production_log_id, downtime_id, minutes)
			VALUES ($1,$2,$3)
		`

		for _, dt := range req.Downtimes {
			if _, err := tx.ExecContext(ctx, dtQuery, id, dt.DowntimeID, dt.Minutes); err != nil {
				return 0, err
			}
		}
	}

	return id, nil

}

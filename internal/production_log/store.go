package productionlog

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateProductionLog(ctx context.Context, tx *sql.Tx, tenantID, userID int64, req CreateProductionLogRequest) (int64, error) {

	// =========================================
	// STEP 1: GET SLOT INDEX (DON'T TRUST UI)
	// =========================================

	var slotIndex int

	slotIndexQuery := `
		SELECT slot_index
		FROM shift_hour_slot
		WHERE id =$1
	`
	err := tx.QueryRowContext(ctx, slotIndexQuery, req.SlotID).Scan(&slotIndex)

	if err != nil {
		return 0, fmt.Errorf("invalid shift slot")
	}

	// =========================================
	// STEP 2: SEQUENCE VALIDATION (🔥 IMPORTANT)
	// =========================================

	var maxSlot sql.NullInt64

	seqVlaildationQuery := `
		SELECT MAX(slot_index)
		FROM production_log
		WHERE tenant_id = $1
			AND machine_id = $2
			AND shift_hour_slot_id = $3
	`
	err = tx.QueryRowContext(ctx, seqVlaildationQuery, tenantID, req.MachineID, req.SlotID).Scan(&maxSlot)
	if err != nil {
		return 0, nil
	}

	if maxSlot.Valid {
		if int(maxSlot.Int64)+1 != slotIndex {
			return 0, fmt.Errorf("invalid slot sequence, expected slot_index %d", maxSlot.Int64+1)
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
		INSERT INTO production_log (tenant_id,machine_id,product_id,operation_id,shift_hour_slot_id,slot_index,quantity,remarks,created_by,updated_by)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$9)
		RETURNING id
	`

	err = tx.QueryRowContext(ctx, insertQuery,
		tenantID,
		req.MachineID,
		req.ProductID,
		req.OperationID,
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

	resourceIDQuery := `
		INSERT INTO production_log_resource(prodution_log_id,resource_id)
		VALUES($1,$2)
	`

	for _, r := range req.ResourceIDs {
		_, err := tx.ExecContext(ctx, resourceIDQuery, id, r)
		if err != nil {
			return 0, err
		}
	}

	// =========================================
	// STEP 5: DEFECT
	// =========================================

	defectInsertQuery := `
		INSERT INTO production_log_defect (production_log_id,defect_id,quantity)
		VALUES($1,$2,$3)
	`

	for _, d := range req.Defects {
		_, err := tx.ExecContext(ctx, defectInsertQuery, id, d.DefectID, d.Quantity)
		if err != nil {
			return 0, nil
		}
	}

	// =========================================
	// STEP 6: DOWNTIME
	// =========================================

	dtinsertQuery := `
		INSERT INTO production_log_downtime(production_log_id,downtime_id,minutes)
		VALUES($1,$2,$3)
	`

	for _, dt := range req.Downtimes {
		_, err := tx.ExecContext(ctx, dtinsertQuery, id, dt.DowntimeID, dt.Minutes)
		if err != nil {
			return 0, err
		}
	}

	return id, nil

}

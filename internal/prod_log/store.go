package prodlog

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

func (s *Store) CreateProdcutionLog(ctx context.Context, tx *sql.Tx, p *ProductionLog) (int64, error) {

	query := `
		INSERT INTO production_log (
		tenant_id,setup_id,machine_id,product_id,operation_id,production_date,shift_name_id,shift_timing_id,shift_hour_slot_id,slot_index,target_qty,actual_qty,ok_qty,rejected_qty,scrap,remarks,created_by,updated_by
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)
		RETURNING id
	`
	var id int64

	err := tx.QueryRowContext(ctx, query,
		p.TenantID,
		p.SetupID,
		p.MachineID,
		p.ProductID,
		p.OperationID,
		p.ProductionDate,
		p.ShiftID,
		p.ShiftTimingID,
		p.ShiftHourSlotID,
		p.SlotIndex,
		p.TargetQty,
		p.ActualQty,
		p.OkQty,
		p.RejectedQty,
		p.Scrap,
		p.Remarks,
		p.CreatedBy,
		p.UpdatedBy,
	).Scan(&id)

	fmt.Println(id)

	return id, err

}

// Defect Prodyuction Log

func (s *Store) UpsertDefectLog(ctx context.Context, tx *sql.Tx, tenantID, productionLogID int64, d DefectInput) error {

	// Upsert defect

	if d.Qty <= 0 {
		return nil
	}

	query := `
		INSERT INTO production_defect (tenant_id,production_log_id,defect_id,qty)
		VALUES($1,$2,$3,$4)
		ON CONFLICT (tenant_id,production_log_id,defect_id)
		DO UPDATE SET qty = EXCLUDED.qty
	`

	_, err := tx.ExecContext(ctx, query,
		tenantID,
		productionLogID,
		d.DefectID,
		d.Qty,
	)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Defect done ")
	return err
	// return err

}

// Upsert Downtime

func (s *Store) UpsertDownTimeLog(ctx context.Context, tx *sql.Tx, tenantID, productionLogID int64, d DowntimeInput) error {
	if d.Minutes <= 0 {
		return nil
	}

	// downtime query

	query := `
	INSERT INTO production_downtime (tenant_id,production_log_id,downtime_id,downtime_minutes)
	VALUES($1,$2,$3,$4)
	ON CONFLICT (tenant_id, production_log_id, downtime_id)
	DO UPDATE SET downtime_minutes = EXCLUDED.downtime_minutes
`

	_, err := tx.ExecContext(ctx, query,
		tenantID,
		productionLogID,
		d.DowntimeID,
		d.Minutes,
	)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Downtime done ")
	return err

}

// soft delete

func (s *Store) SoftDelete(ctx context.Context, id, tenantID, userID int64) error {
	query := `
		UPDATE production_log
		SET 
			is_deleted = TRUE, 
			deleted_at = NOW(),
			deleted_by = $1,
		WHERE 
			id = $2,
			AND tenant_id = $3
			AND is_deleted = FALSE
	`

	result, err := s.db.ExecContext(ctx, query, userID, id, tenantID)

	if err != nil {
		return err
	}

	RowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if RowsAffected == 0 {
		return fmt.Errorf("No record found or already deleted")
	}

	return nil

}

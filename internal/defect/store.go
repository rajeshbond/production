package defect

import (
	"context"
	"database/sql"
	"strings"

	operationdefectmap "github.com/rajesh_bond/production/internal/operation_defect_map"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) BulkCreateDefect(ctx context.Context, tx *sql.Tx, tenantID int64, userID int64, defects []CreateDefectRequest) (BulkDefectResult, error) {

	checkQuery := `
		SELECT id 
		FROM defect
		WHERE tenant_id = $1
		AND LOWER(defect_name) = LOWER($2)
		AND is_deleted = FALSE
	`
	insertQuery := `
		INSERT INTO defect (tenant_id,defect_name,created_by,updated_by)
		VALUES ($1,$2,$3,$3)
		RETURNING id,defect_name
	`
	var results BulkDefectResult

	for _, d := range defects {
		// Step 1: check if exists
		var existingID int64
		err := tx.QueryRowContext(
			ctx,
			checkQuery,
			tenantID,
			strings.ToLower(strings.TrimSpace(d.DefectName)),
		).Scan(&existingID)
		if err != nil && err != sql.ErrNoRows {
			return results, err
		}

		// Already exists -> skip
		if err == nil {
			results.Skipped = append(results.Skipped, d.DefectName)
			continue
		}

		// Step 2: Insert if not exists

		var res DefectResponse

		err = tx.QueryRowContext(ctx, insertQuery, tenantID, d.DefectName, userID).Scan(&res.ID, &res.DefectName)
		if err != nil {
			return results, err
		}
		results.Inserted = append(results.Inserted, res)
	}

	return results, nil
}

func (s *Store) GetDefectIDByName(ctx context.Context, tx *sql.Tx, tenantID int64, defectName string) (int64, error) {
	var id int64

	query := `
		SELECT id
		FROM defect
		WHERE tenant_id = $1
			AND LOWER(defect_name) = LOWER($2)
			AND is_deleted = FALSE
	`
	if err := tx.QueryRowContext(ctx, query, tenantID, defectName).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return id, nil

}

func (s *Store) CreateDefects(ctx context.Context, tx *sql.Tx, tenantID int64, userID int64, defectName []string) ([]int64, error) {

	var reqs []CreateDefectRequest

	for _, name := range defectName {
		reqs = append(reqs, CreateDefectRequest{
			DefectName: name,
		})
	}

	// call your existing functions

	res, err := s.BulkCreateDefect(ctx, tx, tenantID, userID, reqs)
	if err != nil {
		return nil, err
	}

	var ids []int64

	for _, d := range res.Inserted {
		ids = append(ids, d.ID)
	}

	return ids, nil

}

func (s *Store) CreateDefect(ctx context.Context, tx *sql.Tx, tenantID int64, userID int64, defectName string) (int64, error) {
	query := `
		INSERT INTO defect(tenant_id,defect_name,created_by,updated_by)
		VALUES($1,$2,$3,$3)
		ON CONFLICT (tenant_id,defect_name) DO NOTHING
		RETURNING id
	`
	var id int64
	err := tx.QueryRowContext(ctx, query, tenantID, defectName, userID).Scan(&id)
	if err != nil {
		return 0, nil
	}

	return id, nil

}

var _ operationdefectmap.DefectProvider = (*Store)(nil)

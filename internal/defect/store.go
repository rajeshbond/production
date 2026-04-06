package defect

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
		err := tx.QueryRowContext(ctx, checkQuery, tenantID, d.DefectName).Scan(&existingID)
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

package downtime

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

func (s *Store) CheckExistingDownTime(ctx context.Context, tx *sql.Tx, tenantID int64, defectName string) (bool, int64, error) {
	checkQuery := `
		SELECT id 
		FROM downtime
		WHERE tenant_id = $1
		AND LOWER(downtime_name) = LOWER($2)
		AND is_deleted = FALSE
	`

	var id int64

	err := tx.QueryRowContext(ctx, checkQuery, tenantID, defectName).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0, nil // Not found
		}
		return false, 0, err //actual error
	}

	return true, id, nil

}

func (s *Store) BulkCreateDownTime(ctx context.Context, tx *sql.Tx, tenantID int64, userID int64, downtimes []CreateDowntimeRequest) (BulkDownTimeResult, error) {

	insertQuery := `
		INSERT INTO downtime (tenant_id,downtime_name,created_by,updated_by)
		VALUES($1,$2,$3,$3)
		RETURNING id,downtime_name
	`
	var result BulkDownTimeResult

	for _, d := range downtimes {

		// Step 1: Check if exists (using seperate function)
		exists, _, err := s.CheckExistingDownTime(ctx, tx, tenantID, d.DowntimeName)
		if err != nil {
			return result, nil
		}

		// Already Exisits
		if exists {
			result.Skipped = append(result.Skipped, d.DowntimeName)
			continue
		}

		// step 2: Insert if not exists

		var res DownTimeResponse

		err = tx.QueryRowContext(ctx, insertQuery,
			tenantID,
			d.DowntimeName,
			userID,
		).Scan(&res.DowntimeName, &res.DowntimeName)

		if err != nil {
			return result, err
		}

		result.Inserted = append(result.Inserted, res)

	}

	return result, nil

}

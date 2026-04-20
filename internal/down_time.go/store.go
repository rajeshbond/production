package downtime

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	operationdowntimemap "github.com/rajesh_bond/production/internal/operation_downtime_map"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CheckExistingDownTime(ctx context.Context, tx *sql.Tx, tenantID int64, downtimeName string) (bool, int64, error) {
	checkQuery := `
		SELECT id 
		FROM downtime
		WHERE tenant_id = $1
		AND LOWER(downtime_name) = LOWER($2)
		AND is_deleted = FALSE
	`
	fmt.Println("CheckExisiting Down time")
	var id int64

	err := tx.QueryRowContext(ctx, checkQuery, tenantID, downtimeName).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0, nil // Not found
		}
		return false, 0, err //actual error
	}

	return true, id, nil

}

func (s *Store) BulkCreateDownTime(
	ctx context.Context,
	tx *sql.Tx,
	tenantID int64,
	userID int64,
	downtimes []CreateDowntimeRequest,
) (BulkDownTimeResult, error) {

	insertQuery := `
        INSERT INTO downtime (tenant_id, downtime_name, created_by, updated_by)
        VALUES ($1, $2, $3, $3)
        RETURNING id, downtime_name
    `

	var result BulkDownTimeResult

	for _, d := range downtimes {

		// ✅ Normalize input
		normalized := strings.ToLower(strings.TrimSpace(d.DowntimeName))

		// ✅ Skip empty after trim
		if normalized == "" {
			continue
		}

		// Step 1: Check if exists
		exists, _, err := s.CheckExistingDownTime(ctx, tx, tenantID, normalized)
		if err != nil {
			return result, err
		}

		if exists {
			result.Skipped = append(result.Skipped, normalized)
			continue
		}

		// Step 2: Insert
		var res DownTimeResponse

		err = tx.QueryRowContext(ctx, insertQuery,
			tenantID,
			normalized,
			userID,
		).Scan(&res.ID, &res.DowntimeName)

		if err != nil {
			return result, err
		}

		result.Inserted = append(result.Inserted, res)
	}

	return result, nil
}

// func (s *Store) BulkCreateDownTime(ctx context.Context, tx *sql.Tx, tenantID int64, userID int64, downtimes []CreateDowntimeRequest) (BulkDownTimeResult, error) {

// 	insertQuery := `
// 		INSERT INTO downtime (tenant_id,downtime_name,created_by,updated_by)
// 		VALUES($1,$2,$3,$3)
// 		RETURNING id,downtime_name
// 	`
// 	var result BulkDownTimeResult

// 	for _, d := range downtimes {

// 		// Step 1: Check if exists (using seperate function)
// 		exists, _, err := s.CheckExistingDownTime(ctx, tx, tenantID, d.DowntimeName)
// 		if err != nil {
// 			return result, nil
// 		}

// 		// Already Exisits
// 		if exists {
// 			result.Skipped = append(result.Skipped, d.DowntimeName)
// 			continue
// 		}

// 		// step 2: Insert if not exists

// 		var res DownTimeResponse

// 		err = tx.QueryRowContext(ctx, insertQuery,
// 			tenantID,
// 			strings.ToLower(d.DowntimeName),
// 			userID,
// 		).Scan(&res.DowntimeName, &res.DowntimeName)

// 		if err != nil {
// 			return result, err
// 		}

// 		result.Inserted = append(result.Inserted, res)

// 	}

// 	return result, nil

// }

func (s *Store) CreateDowntime(ctx context.Context, tx *sql.Tx, tenantID int64, userID int64, downtimeName string) (int64, error) {

	query := `
		INSERT INTO downtime (tenant_id, downtime_name, created_by, updated_by)
		VALUES ($1, $2, $3, $3)
		ON CONFLICT (tenant_id, LOWER(downtime_name))
		WHERE is_deleted = FALSE
		DO NOTHING
		RETURNING id
	`
	var id int64
	err := tx.QueryRowContext(ctx, query, tenantID, downtimeName, userID).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return id, nil

}

func (s *Store) GetDowntimeIDByName(ctx context.Context, tx *sql.Tx, tenantID int64, downtimeName string) (int64, error) {
	var id int64
	// fmt.Print("Inside get Downtime ------------>", downtimeName)
	query := `
		SELECT id 
		FROM downtime
		WHERE tenant_id = $1
		AND LOWER(downtime_name) = LOWER($2)
		AND is_deleted = FALSE
	`
	err := tx.QueryRowContext(ctx, query, tenantID, downtimeName).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // Not found
		}
		return 0, err //actual error
	}

	return id, nil

}

var _ operationdowntimemap.DowntimeProvider = (*Store)(nil)

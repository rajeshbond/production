package operations

import (
	"context"
	"database/sql"
	"fmt"

	operationdefectmap "github.com/rajesh_bond/production/internal/operation_defect_map"
	operationdowntimemap "github.com/rajesh_bond/production/internal/operation_downtime_map"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CheckExistingOperations(ctx context.Context, tx *sql.Tx, tenantID int64, operationName string) (bool, int64, error) {
	checkQuery := `
	SELECT id 
	FROM operation_master
	WHERE tenant_id = $1
	AND LOWER(operation_name) = LOWER($2)
	AND is_deleted = FALSE
`

	var id int64

	if err := tx.QueryRowContext(ctx, checkQuery, tenantID, operationName).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, 0, nil // Not found
		}

		return false, 0, err // actual error
	}

	return true, id, nil

}

func (s *Store) BulkCreateOperation(ctx context.Context, tx *sql.Tx, tenantID int64, userID int64, operations []CreateOperationsRequest) (BulkOperationResult, error) {

	insertQuery := `
		INSERT INTO operation_master(tenant_id,operation_name,created_by,updated_by)
		VALUES ($1,$2,$3,$3)
		RETURNING id,operation_name
	`

	var results BulkOperationResult

	for _, o := range operations {

		// Step1 : Check if exists (using seperate functions)

		exists, _, err := s.CheckExistingOperations(ctx, tx, tenantID, o.OperationName)
		if err != nil {
			return results, nil
		}

		// Already Exists
		if exists {
			results.Skipped = append(results.Skipped, o.OperationName)
			continue
		}

		// step 2: Insert if not exists

		var res OperationResponse

		err = tx.QueryRowContext(ctx, insertQuery, tenantID, o.OperationName, userID).Scan(&res.ID, &res.OperationName)
		if err != nil {
			return results, err
		}

		results.Inserted = append(results.Inserted, res)

	}

	return results, nil

}

func (s *Store) GetOperationIDByName(ctx context.Context, tx *sql.Tx, tenantID int64, operationName string) (int64, error) {
	var id int64

	query := `
		SELECT id 
		from operation_master
		WHERE tenant_id = $1
		 AND LOWER(operation_name) = LOWER($2)
		 AND is_deleted = FALSE
	`

	if err := tx.QueryRowContext(ctx, query, tenantID, operationName).Scan(&id); err != nil {

		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return id, nil
}

func (s *Store) CreateOperation(ctx context.Context, tx *sql.Tx, tenantID int64, userID int64, operationName string) (int64, error) {
	var id int64
	fmt.Println("Inside create operation of store")
	query := `
		INSERT INTO operation_master ( tenant_id, operation_name,created_by,updated_by)
		VALUES($1,$2,$3,$3)
		ON CONFLICT (tenant_id,operation_name) DO NOTHING
		RETURNING id
	`

	if err := tx.QueryRowContext(ctx, query, tenantID, operationName, userID).Scan(&id); err != nil {
		return 0, err
	}

	fmt.Println("store ID", id)

	return id, nil

}

var _ operationdefectmap.OperationProvider = (*Store)(nil)

var _ operationdowntimemap.OperationProvider = (*Store)(nil)

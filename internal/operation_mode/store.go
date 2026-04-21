package operationmode

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

func (s *Store) Create(ctx context.Context, tx *sql.Tx, tenantID, userID int64, name string) (int64, error) {

	query := `
		INSERT INTO operation_mode (tenant_id,operation_mode_name,created_by, updated_by)
		VALUES($1,$2,$3,$3)
		RETURNING id;
	`

	var id int64

	err := tx.QueryRowContext(ctx, query, tenantID, name, userID).Scan(&id)
	return id, err

}

func (s *Store) Update(ctx context.Context, tx *sql.Tx, tenantID, userID int64, req UpdateOperationModeRequest) error {
	query := `
	UPDATE operation_mode
	SET operation_mode_name = $1,
	    is_active = COALESCE($2, is_active),
	    updated_by = $3,
	    updated_at = NOW()
	WHERE id = $4
	AND tenant_id = $5
	AND is_deleted = FALSE;
	`

	_, err := tx.ExecContext(ctx, query,
		req.OperationModeName,
		req.IsActive,
		userID,
		req.ID,
		tenantID,
	)

	return err
}

func (s *Store) Delete(ctx context.Context, tx *sql.Tx, tenantID, userID, id int64) error {

	query := `
	UPDATE operation_mode
	SET is_deleted = TRUE,
	    deleted_by = $1,
	    deleted_at = NOW()
	WHERE id = $2
	AND tenant_id = $3
	AND is_deleted = FALSE;
	`

	_, err := tx.ExecContext(ctx, query, userID, id, tenantID)
	return err
}

func (s *Store) List(ctx context.Context, tenantID int64) ([]OperationMode, error) {

	query := `
	SELECT id, operation_mode_name, is_active
	FROM operation_mode
	WHERE tenant_id = $1
	AND is_deleted = FALSE;
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []OperationMode

	for rows.Next() {
		var m OperationMode
		if err := rows.Scan(&m.ID, &m.OperationModeName, &m.IsActive); err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, nil
}

func (s *Store) GetIDByName(ctx context.Context, tenantID int64, name string) (int64, error) {

	query := `
	SELECT id
	FROM operation_mode
	WHERE tenant_id = $1
	  AND LOWER(operation_mode_name) = LOWER($2)
	  AND is_deleted = FALSE
	LIMIT 1;
	`

	var id int64
	err := s.db.QueryRowContext(ctx, query, tenantID, name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("operation mode not found")
		}
		return 0, err
	}

	return id, nil
}

package resourcemaster

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

func (s *Store) CreateResource(ctx context.Context, tx *sql.Tx, tenantID int64, resourceName string, typeID, userID int64) (int64, error) {
	var id int64

	query := `
		INSERT INTO resource_master (tenant_id,resource_name,resource_type_id,created_by,updated_by)
		VALUES($1,$2,$3,$4,$4)
		ON CONFLICT (tenant_id,resresource_name)
		DO UPDATE SET 
			resource_type_id = EXCLUDED.resource_type_id,
			updated_by = EXCLUDED.updated_by,
			updated_at = NOW()
		RETURNING id;
	`

	err := tx.QueryRowContext(ctx, query, tenantID, resourceName, typeID, userID).Scan(&id)

	return id, err
}

// GET BY ID

func (s *Store) GetResourceByID(ctx context.Context, tx *sql.Tx, tenantID, resourceID int64) (ResourceResponse, error) {

	qyery := `
		SELECT id , resource_name, resource_type_id
		FROM resource_master
		WHERE id = $1
		AND tenant_id = $2
		AND is_deleted = FALSE;
	`
	var r ResourceResponse

	if err := tx.QueryRowContext(ctx, qyery, resourceID, tenantID).Scan(&r.ID, &r.ResourceName, &r.ResourceTypeID); err != nil {
		if err == sql.ErrNoRows {
			return ResourceResponse{}, fmt.Errorf("resouce not found")
		}
		return ResourceResponse{}, nil
	}

	return r, nil

}

func (s *Store) GetAllResources(ctx context.Context, tenantID int64) ([]ResourceResponse, error) {
	query := `
		SELECT id,resource_name,resource_type_id
		FROM resource_master
		WHERE tenant_id =$1
		AND is_deleted = FALSE
		ORDER BY id DESC;
	`
	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []ResourceResponse

	for rows.Next() {
		var r ResourceResponse

		err := rows.Scan(&r.ID, &r.ResourceName, &r.ResourceTypeID)
		if err != nil {
			return nil, err
		}
		list = append(list, r)

	}

	return list, rows.Err()

}

func (s *Store) UpdateResorce(ctx context.Context, tx *sql.Tx, tenantID, resourceID int64, resourceName string, typeID, userID int64) error {
	query := `
	UPDATE resource_master
	SET resource_name = $1,
	    resource_type_id = $2,
	    updated_by = $3,
	    updated_at = NOW()
	WHERE id = $4
	AND tenant_id = $5
	AND is_deleted = FALSE;
	`

	result, err := tx.ExecContext(
		ctx,
		query,
		resourceName,
		typeID,
		userID,
		resourceID,
		tenantID,
	)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("resource not found")
	}

	return nil
}

func (s *Store) DeleteResource(ctx context.Context, tx *sql.Tx, tenantID, resourceID, userID int64) error {

	query := `
	UPDATE resource_master
	SET is_deleted = TRUE,
	    deleted_by = $1,
	    deleted_at = NOW()
	WHERE id = $2
	AND tenant_id = $3
	AND is_deleted = FALSE;
	`

	result, err := tx.ExecContext(ctx, query, userID, resourceID, tenantID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("resource not found")
	}

	return nil
}

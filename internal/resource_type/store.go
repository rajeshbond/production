package resourcetype

import (
	"context"
	"database/sql"
	"errors"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateResourceType(ctx context.Context, tenantID int64, userID int64, typeName, description string) (int64, error) {
	query := `
	INSERT INTO resource_type_master
	    (tenant_id, type_name, description, created_by, updated_by)
	VALUES ($1, $2, $3, $4, $4)
	RETURNING id;
	`
	var id int64

	err := s.db.QueryRowContext(ctx, query, tenantID, typeName, description, userID).Scan(&id)

	return id, err
}

func (s *Store) UpdateResourceType(ctx context.Context, tenantID, userID, id int64, typeName, description string) (int64, error) {
	query := `
		UPDATE resource_type_master
		SET 
			type_name = $1,
			description = $2,
			updated_by = $3,
		WHERE id = $4,
			AND tenant_id = $5,
			AND is_deleted = FALSE
		RETURNING id;
	`
	var updatedID int64

	err := s.db.QueryRowContext(ctx, query, typeName, description, userID, id, tenantID).Scan(&updatedID)

	if err == sql.ErrNoRows {
		return 0, errors.New("resource type not found or duplicate")
	}

	return updatedID, err
}

func (s *Store) DeleteResourceType(ctx context.Context, tenantID, userID, id int64) (int64, error) {

	query := `
	UPDATE resource_type_master
	SET 
	    is_deleted = TRUE,
	    deleted_by = $1
	WHERE id = $2
	  AND tenant_id = $3
	  AND is_deleted = FALSE
	RETURNING id;
	`

	var deletedID int64
	err := s.db.QueryRowContext(ctx, query,
		userID, id, tenantID,
	).Scan(&deletedID)

	if err == sql.ErrNoRows {
		return 0, errors.New("resource type not found")
	}

	return deletedID, err
}

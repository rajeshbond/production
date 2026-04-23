package resource

import (
	"context"
	"database/sql"

	"github.com/rajesh_bond/production/internal/fixture"
	"github.com/rajesh_bond/production/internal/mould"
	"github.com/rajesh_bond/production/internal/tools"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Create(ctx context.Context, tx *sql.Tx, r *Resource) (int64, error) {
	query := `
	INSERT INTO resource (resource_sub_id,tenant_id,resource_code,resource_name,resource_type,description,created_by,updated_by)
	VALUES($1,$2,$3,$4,$5,$6,$7,$7)
	RETURNING id
`

	var id int64
	err := tx.QueryRowContext(ctx, query,
		r.ResourceSubID,
		r.TenantID,
		r.ResourceCode,
		r.ResourceName,
		r.ResourceType,
		r.Description,
		r.CreatedBy,
		r.CreatedBy,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil

}

func (s *Store) CreateResource(ctx context.Context, tx *sql.Tx, rid, tenantID, userID int64, resourceCode, resourceName, resourceType, description string) (int64, error) {
	query := `
	INSERT INTO resource (resource_sub_id,tenant_id,resource_code,resource_name,resource_type,description,created_by,updated_by)
	VALUES($1,$2,$3,$4,$5,$6,$7,$7)
	RETURNING id
`
	var id int64
	err := tx.QueryRowContext(ctx, query,
		rid,
		tenantID,
		resourceCode,
		resourceName,
		resourceType,
		description,
		userID,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil

}

//  for Tool Resourece Provide

var _ tools.ResourceProvider = (*Store)(nil)
var _ mould.ResourceProvider = (*Store)(nil)
var _ fixture.ResourceProvider = (*Store)(nil)

package resource

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

func (s *Store) Create(ctx context.Context, tx *sql.Tx, r *Resource) (int64, error) {
	query := `
	INSERT INTO resource (tenant_id,rescouce_code,resource_name,resource_type,description,mold_id,fixture_id,tool_id,created_by,updated_by)
	VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$9)
	RETURNING id
`

	var id int64
	err := tx.QueryRowContext(ctx, query,
		r.TenantID,
		r.ResourceCode,
		r.ResourceName,
		r.ResourceType,
		r.Description,
		r.MoldID,
		r.FixtureID,
		r.ToolID,
		r.CreatedBy,
		r.CreatedBy,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil

}

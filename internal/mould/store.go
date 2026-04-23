package mould

import (
	"context"
	"database/sql"
	"strings"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Create(ctx context.Context, tx *sql.Tx, m *Mold) (int64, error) {

	query := `
		INSERT INTO mold (tenant_id,type,mold_name,mold_no,description,cavities,target_shots,created_by,updated_by)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$8)
		RETURNING id
	`
	var id int64

	typeName := strings.ToLower(m.Type)

	err := tx.QueryRowContext(
		ctx,
		query,
		m.TenantID,
		typeName,
		m.MoldName,
		m.MoldNo,
		m.Description,
		m.Cavities,
		m.TargetShots,
		m.CreatedBy,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil

}

package tools

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

func (s *Store) Create(ctx context.Context, tx *sql.Tx, t *Tool) (int64, error) {
	query := `
		INSERT INTO tool_master(tenant_id,type,tool_code,tool_name,description,tool_type,unit,cost,life_cycles,created_by,updated_by)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$10)
		RETURNING id
	`

	typeName := strings.ToLower(t.Type)

	var id int64
	err := tx.QueryRowContext(ctx, query,
		t.TenantID,
		typeName,
		t.ToolCode,
		t.ToolName,
		t.Description,
		t.ToolType,
		t.Unit,
		t.Cost,
		t.LifeCycles,
		t.CreatedBy,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil

}

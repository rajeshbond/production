package operationdefectmap

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

func (s *Store) InsertOperationDefectMap(ctx context.Context, tx *sql.Tx, tenantID, operationID, defectID int64) (int64, error) {

	query := `
		INSERT INTO operation_defect_map(tenant_id,operation_id,defect_id)
		VALUES($1,$2,$3)
		ON CONFLICT (tenant_id,operation_id,defect_id)
		DO NOTHING
		RETURNING id
	`
	var id int64

	if err := tx.QueryRowContext(ctx, query, tenantID, operationID, defectID).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return id, nil

}

package operationdowntimemap

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

func (s *Store) InsertOperationDowntimeMap(ctx context.Context, tx *sql.Tx, tenantID int64, operationID int64, downTimeID int64, userID int64) (int64, error) {

	var id int64

	query := `
		INSERT INTO operation_downtime_map(tenant_id,operation_id,downtime_id,created_by,updated_by)
		VALUES($1,$2,$3,$4,$4)
		ON CONFLICT (tenant_id,operation_id,downtime_id)
		DO NOTHING
		RETURNING id
	`
	err := tx.QueryRowContext(ctx, query, tenantID, operationID, downTimeID, userID).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		return 0, err
	}
	return id, nil
}

func (s *Store) GetOperationDowntimeMap(ctx context.Context, tx *sql.Tx, tenantID int64, operationID int64, downTimeID int64) (int64, error) {
	var id int64

	query := `
		SELECT id 
		FROM operation_downtime_map 
		WHERE tenant_id = $1
		AND operation_id = $2
		AND downtime_id = $3
	`
	if err := tx.QueryRowContext(ctx, query, tenantID, operationID, downTimeID).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return id, nil

}

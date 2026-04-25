package setupoperation

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateSetup(ctx context.Context, tx *sql.Tx, o *OperationSetup) (int64, error) {
	query := `
	INSERT INTO operation_setup (tenant_id,pos_id,machine_id,setup_name,target_qty,cycle_time_sec,setup_time_min,created_by,updated_by)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	RETURNING id
`
	var id int64

	err := tx.QueryRowContext(ctx, query,
		o.TenantID,
		o.PosID,
		o.MachineID,
		o.SetupName,
		o.TargetQty,
		o.CycleTimeSec,
		o.SetupTimeMin,
		o.CreatedBy,
		o.UpdatedBy,
	).Scan(&id)

	if err != nil {
		fmt.Println("Resource========>", err)
		return 0, err
	}

	return id, nil

}

func (s *Store) UpsertResource(ctx context.Context, tx *sql.Tx, tenantID, setupID int64, items []ResourceItem) error {
	var ids []int64

	var qtys []int

	for _, r := range items {
		ids = append(ids, r.ResourceID)
		if r.Quantity <= 0 {
			qtys = append(qtys, 1)
		} else {
			qtys = append(qtys, r.Quantity)
		}
	}

	_, err := tx.ExecContext(ctx, `
	INSERT INTO operation_setup_resource_map (tenant_id,setup_id,resource_id,quantity)
	SELECT $1,$2,x.resource_id,x.quantity
	FROM UNNEST($3::BIGINT[],$4::INT[]) AS x(resource_id,quantity)
	ON CONFLICT (tenant_id,setup_id,resource_id)
	DO UPDATE SET quantity=EXCLUDED.quantity`,
		tenantID, setupID,
		pq.Array(ids), pq.Array(qtys),
	)
	fmt.Println("Resource========>", err)
	return err
}

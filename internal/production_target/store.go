package productiontarget

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// Generarte Signature

func generateSignature(resourceIDs []int64) string {
	if len(resourceIDs) == 0 {
		return ""
	}

	sort.Slice(resourceIDs, func(i, j int) bool {
		return resourceIDs[i] < resourceIDs[j]
	})

	var parts []string
	for _, id := range resourceIDs {
		parts = append(parts, fmt.Sprintf("%d", id))
	}

	return strings.Join(parts, ",")

}

// Create

func (s *Store) CreateProductionTarget(ctx context.Context, tx *sql.Tx, tenantID int64, req *CreateProductionTargetRequest, userID int64) (int64, error) {

	signature := generateSignature(req.ResourceIDs)

	query := `
		INSERT INTO production_target (tenant_id,product_id,operation_id,machine_id,process_type,target_per_hour,expected_efficiency,resource_signature,created_by,updated_by)
		VALUES($1	,$2,$3,$4,$5,$6,$7,$8,$9,$9)
		RETURNING id;
	`
	var id int64
	fmt.Println("Store Create Production Target")
	err := tx.QueryRowContext(ctx, query,
		tenantID,
		req.ProductID,
		req.OperationID,
		req.MachineID,
		req.ProcessType,
		req.TargetPerHour,
		req.ExpectedEfficiency,
		signature,
		userID,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	fmt.Println("ID:-", id)
	// insert resource mapping

	queryInsertResource := `
		INSERT INTO production_target_resource_map(tenant_id,production_target_id,resource_id)VALUES(	$1,$2,$3);
	`

	for _, rID := range req.ResourceIDs {
		_, err := tx.ExecContext(ctx, queryInsertResource, tenantID, id, rID)
		if err != nil {
			fmt.Print(err)
			return 0, err
		}

	}

	return id, nil

}

func (s *Store) UpdateProductionTarget(ctx context.Context, tx *sql.Tx, tenantID, userID int64, req *UpdateProductionTargetRequest) error {

	signature := generateSignature(req.ResourceIDs)

	query := `
	UPDATE production_target
	SET product_id = $1,
	    operation_id = $2,
	    machine_id = $3,
	    process_type = $4,
	    target_per_hour = $5,
	    expected_efficiency = $6,
	    resource_signature = $7,
	    updated_by = $8,
	    updated_at = NOW()
	WHERE id = $9
	AND tenant_id = $10
	AND is_deleted = FALSE;
	`

	res, err := tx.ExecContext(ctx, query,
		req.ProductID,
		req.OperationID,
		req.ProcessType,
		req.TargetPerHour,
		req.ExpectedEfficiency,
		signature,
		userID,
		req.ID,
		tenantID,
	)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	// delete old resource mapping
	deleteQuery := `
		DELETE FROM production_target_resource_map
		WHERE production_target_id = $1;`

	_, err = tx.ExecContext(ctx, deleteQuery, req.ID)

	if err != nil {
		return err
	}

	// insert new mapping

	insertQuery := `
		INSERT INTO production_target_resource_map
		(tenant_id,production_target_id,resource_id)
		VALUES($1,$2,$3);
	`
	for _, rID := range req.ResourceIDs {
		_, err = tx.ExecContext(ctx, insertQuery, tenantID, req.ID, rID)
		if err != nil {
			return err
		}
	}

	return nil

}

// Get Production Target By ID

func (s *Store) GetProductionTargetByID(
	ctx context.Context,
	tenantID int64,
	id int64,
) (*ProductionTargetResponseDTO, error) {

	query := `
	SELECT 
		pt.id,
		pt.product_id,
		pt.operation_id,
		pt.machine_id,
		pt.process_type,
		pt.target_per_hour,
		pt.expected_efficiency,
		COALESCE(array_agg(ptr.resource_id), '{}') as resource_ids
	FROM production_target pt
	LEFT JOIN production_target_resource_map ptr
		ON pt.id = ptr.production_target_id
	WHERE pt.id = $1
	AND pt.tenant_id = $2
	AND pt.is_deleted = FALSE
	GROUP BY pt.id;
	`

	var resp ProductionTargetResponseDTO
	var resourceIDs []int64

	err := s.db.QueryRowContext(ctx, query, id, tenantID).
		Scan(
			&resp.ID,
			&resp.ProductID,
			&resp.OperationID,
			&resp.MachineID,
			&resp.ProcessType,
			&resp.TargetPerHour,
			&resp.ExpectedEfficiency,
			pq.Array(&resourceIDs),
		)

	if err != nil {
		return nil, err
	}

	resp.ResourceIDs = resourceIDs

	return &resp, nil
}

func (s *Store) GetAllProductionTargets(ctx context.Context, tenantID int64) ([]ProductionTargetResponseDTO, error) {

	query := `
	SELECT 
		pt.id,
		pt.product_id,
		pt.operation_id,
		pt.machine_id,
		pt.process_type,
		pt.target_per_hour,
		pt.expected_efficiency,
		COALESCE(array_agg(ptr.resource_id), '{}') as resource_ids
	FROM production_target pt
	LEFT JOIN production_target_resource_map ptr
		ON pt.id = ptr.production_target_id
	WHERE pt.tenant_id = $1
	AND pt.is_deleted = FALSE
	GROUP BY pt.id
	ORDER BY pt.id DESC;
	`
	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []ProductionTargetResponseDTO

	for rows.Next() {
		var r ProductionTargetResponseDTO
		var resourceIDs []int64

		err := rows.Scan(
			&r.ID,
			&r.ProductID,
			&r.OperationID,
			&r.MachineID,
			&r.ProcessType,
			&r.TargetPerHour,
			&r.ExpectedEfficiency,
			pq.Array(&resourceIDs),
		)
		if err != nil {
			return nil, err
		}
		r.ResourceIDs = resourceIDs
		result = append(result, r)
	}

	return result, nil

}

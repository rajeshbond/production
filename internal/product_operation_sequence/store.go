package productoperationsequence

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

// 🔒 Lock rows for concurrency (IMPORTANT FIX HERE)
func (s *Store) LockRows(ctx context.Context, tx *sql.Tx, tenantID int64, productID int64) error {
	query := `
		SELECT id 
		FROM product_operation_sequence
		WHERE tenant_id = $1 AND product_id = $2
		FOR UPDATE
	`

	_, err := tx.ExecContext(ctx, query, tenantID, productID)
	return err
}

// 📌 Get next sequence number
func (s *Store) GetNextSequenceNo(ctx context.Context, tx *sql.Tx, tenantID int64, productID int64) (int, error) {
	query := `
		SELECT COALESCE(MAX(sequence_no), 0) + 1
		FROM product_operation_sequence 
		WHERE tenant_id = $1 AND product_id = $2
	`

	var seq int
	err := tx.QueryRowContext(ctx, query, tenantID, productID).Scan(&seq)

	return seq, err
}

// ✅ Create single record
func (s *Store) Create(ctx context.Context, tx *sql.Tx, p *ProductOperationSequence) error {

	// ensure updated_by = created_by
	p.UpdatedBy = p.CreatedBy

	query := `
	INSERT INTO product_operation_sequence
	(tenant_id, product_id, operation_id, sequence_no, created_by, updated_by)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at, updated_at;
	`

	err := tx.QueryRowContext(
		ctx,
		query,
		p.TenantID,
		p.ProductID,
		p.OperationID,
		p.SequenceNo,
		p.CreatedBy,
		p.UpdatedBy,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return fmt.Errorf("insert failed: %w", err)
	}

	return nil
}

// 🚀 Bulk insert (row-by-row safe version)
func (s *Store) BulkCreate(ctx context.Context, tx *sql.Tx, records []ProductOperationSequence) error {

	query := `
	INSERT INTO product_operation_sequence
	(tenant_id, product_id, operation_id, sequence_no, created_by, updated_by)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at, updated_at;
	`

	for i := range records {

		records[i].UpdatedBy = records[i].CreatedBy

		err := tx.QueryRowContext(
			ctx,
			query,
			records[i].TenantID,
			records[i].ProductID,
			records[i].OperationID,
			records[i].SequenceNo,
			records[i].CreatedBy,
			records[i].UpdatedBy,
		).Scan(
			&records[i].ID,
			&records[i].CreatedAt,
			&records[i].UpdatedAt,
		)

		if err != nil {
			return fmt.Errorf("bulk insert failed: %w", err)
		}
	}

	return nil
}

// 📊 Get current max sequence
func (s *Store) GetCurrentMaxSequence(ctx context.Context, tx *sql.Tx, tenantID, productID int64) (int, error) {

	query := `
	SELECT COALESCE(MAX(sequence_no), 0)
	FROM product_operation_sequence
	WHERE tenant_id = $1 AND product_id = $2;
	`

	var maxSeq int
	err := tx.QueryRowContext(ctx, query, tenantID, productID).Scan(&maxSeq)

	return maxSeq, err
}

// package productoperationsequence

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// )

// type Store struct {
// 	db *sql.DB
// }

// func NewStore(db *sql.DB) *Store {
// 	return &Store{db: db}
// }

// func (s *Store) LockRows(ctx context.Context, tx *sql.Tx, tenantID int64, productID int64) error {
// 	query := `
// 		SELECT id FROM product_operation_sequence
// 		WHERE tenant_id = $1 AND product_id = $2
// 		FROM UPDATE
// 	`

// 	_, err := tx.ExecContext(ctx, query, tenantID, productID)

// 	return err
// }

// func (s *Store) GetNextSequenceNo(ctx context.Context, tx *sql.Tx, tenantID int64, productID int64) (int, error) {
// 	query := `
// 		SELECT COALESCE(MAX(sequence_no),0) + 1
// 		FROM product_operation_sequence
// 		WHERE tenant_id =$1 AND product_id = $2
// 	`
// 	var seq int

// 	err := tx.QueryRowContext(ctx, query, tenantID, productID).Scan(&seq)

// 	return seq, err
// }

// func (s *Store) Create(ctx context.Context, tx *sql.Tx, p *ProductOperationSequence) error {

// 	// ensure updated_by = created_by
// 	p.UpdatedBy = p.CreatedBy

// 	query := `
// 	INSERT INTO product_operation_sequence
// 	(tenant_id, product_id, operation_id, sequence_no, created_by, updated_by)
// 	VALUES ($1, $2, $3, $4, $5, $6)
// 	RETURNING id, created_at, updated_at;
// 	`

// 	err := tx.QueryRowContext(
// 		ctx,
// 		query,
// 		p.TenantID,
// 		p.ProductID,
// 		p.OperationID,
// 		p.SequenceNo,
// 		p.CreatedBy,
// 		p.UpdatedBy,
// 	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

// 	if err != nil {
// 		return fmt.Errorf("insert failed: %w", err)
// 	}

// 	return nil
// }

// func (s *Store) BulkCreate(ctx context.Context, tx *sql.Tx, records []ProductOperationSequence) error {

// 	query := `
// 	INSERT INTO product_operation_sequence
// 	(tenant_id, product_id, operation_id, sequence_no, created_by, updated_by)
// 	VALUES ($1, $2, $3, $4, $5, $6)
// 	RETURNING id, created_at, updated_at;
// 	`

// 	for i := range records {

// 		records[i].UpdatedBy = records[i].CreatedBy

// 		err := tx.QueryRowContext(
// 			ctx,
// 			query,
// 			records[i].TenantID,
// 			records[i].ProductID,
// 			records[i].OperationID,
// 			records[i].SequenceNo,
// 			records[i].CreatedBy,
// 			records[i].UpdatedBy,
// 		).Scan(
// 			&records[i].ID,
// 			&records[i].CreatedAt,
// 			&records[i].UpdatedAt,
// 		)

// 		if err != nil {
// 			return fmt.Errorf("bulk insert failed: %w", err)
// 		}
// 	}

// 	return nil
// }

// func (s *Store) GetCurrentMaxSequence(ctx context.Context, tx *sql.Tx, tenantID, productID int64) (int, error) {

// 	query := `
// 	SELECT COALESCE(MAX(sequence_no), 0)
// 	FROM product_operation_sequence
// 	WHERE tenant_id = $1 AND product_id = $2;
// 	`

// 	var maxSeq int
// 	err := tx.QueryRowContext(ctx, query, tenantID, productID).Scan(&maxSeq)
// 	return maxSeq, err
// }

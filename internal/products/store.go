package products

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateProduct(ctx context.Context, tx *sql.Tx, tenantID int64, productName string,
	productNumber string, userID int64,
) (int64, error) {

	var id int64

	query := `
		INSERT INTO product 
		(tenant_id, product_name, product_no, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $4)
		ON CONFLICT (tenant_id, product_no)
		DO UPDATE SET
			product_name = EXCLUDED.product_name,
			updated_by = EXCLUDED.updated_by,
			updated_at = NOW()
		RETURNING id;
	`

	err := tx.QueryRowContext(
		ctx,
		query,
		tenantID,
		productName,
		productNumber,
		userID,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to upsert product: %w", err)
	}

	return id, nil
}

// Get Product by ID

func (s *Store) GetAllProductsByTenant(ctx context.Context, tenantID int64) ([]Product, error) {

	query := `
	SELECT id, tenant_id, product_name, product_no, created_at, updated_at
	FROM product
	WHERE tenant_id = $1
	AND is_deleted = FALSE
	ORDER BY id DESC;
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var p Product

		err := rows.Scan(
			&p.ID,
			&p.TenantID,
			&p.ProductName,
			&p.ProductNo,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (s *Store) DeleteProduct(ctx context.Context, tx *sql.Tx, tenantID int64, productID int64, userID int64) error {

	query := `
		UPDATE product
		SET is_deleted = TRUE,
		    deleted_by = $1,
		    deleted_at = NOW()
		WHERE id = $2
		AND tenant_id = $3
		AND is_deleted = FALSE;
	`

	result, err := tx.ExecContext(ctx, query, userID, productID, tenantID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("product not found or already deleted")
	}

	return nil
}

func (s *Store) GetProductByID(ctx context.Context, tx *sql.Tx, tenantID int64, productID int64) (string, string, error) {

	query := `
		SELECT product_name, product_no
		FROM product
		WHERE id = $1
		AND tenant_id = $2
		AND is_deleted = FALSE;
	`

	var productName, productNumber string

	err := tx.QueryRowContext(ctx, query, productID, tenantID).
		Scan(&productName, &productNumber)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", fmt.Errorf("product not found")
		}
		return "", "", err
	}

	return productName, productNumber, nil
}

// Product Search

func (s *Store) SearchProducts(ctx context.Context, tenantID int64, search string) ([]ProductResponse, error) {
	query := `
		SELECT product_name,product_id
		FROM product 
		WHERE tenant_id = $1 
		AND is_deleted = FALSE
		AND(
			product_name = ILIKE '%' || $2 || '%'
			OR 
			product_no = ILIKE '%' || $2 || '%'
		)
		ORDER BY id DESC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID, search)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []ProductResponse

	for rows.Next() {
		var p ProductResponse
		err := rows.Scan(&p.ID, &p.ProductName, &p.ProductNo)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (s *Store) GetProductByNoOrName(ctx context.Context, tx *sql.Tx, tenantID int64, value string) (ProductResponse, error) {
	query := `
		SELECT id,product_name,product_no
		FROM product 
		WHERE tenant_id = $1
		AND is_deleted = FALSE
		AND(
		product_no = $2 
		OR product_name = $2 
		)
		LIMIT;
	`

	var p ProductResponse

	err := tx.QueryRowContext(ctx, query, tenantID, value).Scan(&p.ID, &p.ProductName, &p.ProductNo)
	if err != nil {
		if err == sql.ErrNoRows {
			return ProductResponse{}, fmt.Errorf("Product not found")
		}
		return ProductResponse{}, err
	}

	return p, nil
}

func (s *Store) UpdateProduct(ctx context.Context, tx *sql.Tx, tenantID int64, userID int64, req *UpdateProductRequest) error {
	query := `
		UPDATE product
		SET product_name =$1
			product_no = $2,
			updated_by = $3,
			updated_at = NOW()
		WHERE id = $4
		AND tenant_id = $5
		AND is_deleted = FALSE;
	`
	result, err := tx.ExecContext(ctx, query, req.ProductName, req.ProductNo, userID, req.ID, tenantID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("Product bot found or already deleted")
	}

	return nil

}

// func (s *Store) DeleteProduct(ctx context.Context, tx *sql.Tx, tenantID int64, productID int64, userID int64) error {
// 	query := `
// 		UPDATE product
// 		SET is_deleted = TRUE,
// 			deleted_by = $1,
// 			deleted_at = NOW()
// 		WHERE id = $2,
// 		AND tenantID = $3,
// 		AND is_deleted = FALSE;
// 	`

// 	result, err := tx.ExecContext(ctx, query, userID, productID, tenantID)

// 	if err != nil {
// 		return err
// 	}

// 	rowsAffected, _ := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		return fmt.Errorf("Product not found or already deleted")
// 	}

// 	return nil

// }

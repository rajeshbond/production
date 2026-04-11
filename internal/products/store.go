package products

import (
	"context"
	"database/sql"
	"errors"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateProduct(ctx context.Context, tx *sql.Tx, tenantID int64, productName string, productNumber string, userID int64) (int64, error) {
	var id int64

	query := `
		INSERT INTO product (tenant_id,product_name,product_no,created_by,updated_by)
		VALUES ($1,$2,$3,$4,$4)
		ON CONFLICT (tenant_id,product_no)
		ON UPDATE SET
			product_name = EXCLUDED.productName,
			update_by = EXCLUDED.updated_by
			updated_at = NOW()
		RETURNING id;
	`

	err := tx.QueryRowContext(ctx, query, tenantID, productName, productNumber, userID).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil

}

// Get Product by ID

func (s *Store) GetProductByID(ctx context.Context, tx *sql.Tx, tenantID int64, productID int64) (string, string, error) {

	query := `

		SELECT product_name,product_no
		FROM product
		WHERE id = $1
		AND 
		tenant_id = $2 
	`

	var productName, productNumber string

	err := tx.QueryRowContext(ctx, query, productID, tenantID).Scan(&productName, productNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", errors.New("product not found")
		}

		return "", "", err
	}

	return productName, productNumber, nil

}

func (s *Store) GetAllProductsByTenant(ctx context.Context, tenantID int64,
) ([]Product, error) {

	query := `
	SELECT id, tenant_id, product_name, product_no, created_at, updated_at
	FROM product
	WHERE tenant_id = $1
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

	// ✅ VERY IMPORTANT: check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// ✅ return result
	return products, nil
}

package tenant

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

func (s *Store) Create(ctx context.Context, dto CreateTenantDTO) (*Tenant, error) {
	query := `
		INSERT INTO tenant (tenant_name, tenant_code, address,created_by,updated_by) 
		VALUES ($1,$2,$3,$4,$4)
		RETURNING id, tenant_name, tenant_code, address,
		          is_verified, is_active, created_by, updated_by,
		          created_at, updated_at
	`

	var t Tenant

	err := s.db.QueryRowContext(
		ctx,
		query,
		dto.TenantName,
		dto.TenantCode,
		dto.Address,
		dto.CreatedBy,
	).Scan(
		&t.ID,
		&t.TenantName,
		&t.TenantCode,
		&t.Address,
		&t.IsVerified,
		&t.IsActive,
		&t.CreatedBy,
		&t.UpdatedBy,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		return nil, err // ⭐ important improvement
	}

	return &t, nil
}

func (s *Store) GetTenantIDByName(ctx context.Context, tenantName string) (int64, error) {
	query := `
		SELECT id
		FROM tenant
		WHERE LOWER(tenant_name) = LOWER($1)
	`

	var tenantID int64

	err := s.db.QueryRowContext(ctx, query, tenantName).Scan(&tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("Role not found")
		}

		return 0, err
	}

	return tenantID, nil
}

func (s *Store) GetTenantNameByID(ctx context.Context, tenantID int64) (string, error) {
	query := `
		SELECT tenant_name
		FROM tenant
		WHERE id = $1
	`

	var tenatName string

	err := s.db.QueryRowContext(
		ctx,
		query,
		tenantID).Scan(&tenantID)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("No Role found")
		}

		return "", err
	}

	return tenatName, nil

}

func (s *Store) GetTenantCodeByID(ctx context.Context, tenantID int64) (string, error) {

	query := `
		SELECT tenant_code
		FROM tenant
		WHERE id = $1
	`

	var tenantCode string

	err := s.db.QueryRowContext(ctx, query, tenantID).Scan(&tenantCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("tenant not found")
		}
		return "", err
	}

	return tenantCode, nil
}

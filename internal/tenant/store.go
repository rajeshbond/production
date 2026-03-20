package tenant

// Index ---->

// 1. Create Tenant by suer Admin or Admin

// 2. Get Tenant ID by Name

// 3. Get Tenant by Name

// 4. Get Tenant Code by ID

// 5. Create super Tenant Tx - Development

// 6. Check the tenant code in DB

// 7. Tenant verification store function

// 8. Delete Tenant by Code

// 10. Update the Tenant by Tenant code

//<--- Code Starts here --> //

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// 1. Create Tenant by suer Admin or Admin

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

// 2. Get Tenant ID by Name

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

// 3. Get Tenant by Name

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

// 4. Get Tenant Code by ID
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

// 5. Create super Tenant Tx - Development

func (s *Store) CreateSuperTenantTx(ctx context.Context, tx *sql.Tx, dto CreateTenantDTO) (int64, error) {

	query := `
		INSERT INTO tenant (tenant_name,tenant_code,address,created_by,updated_by)
		VALUES($1,$2,$3,$4,$5)
		RETURNING id
	`
	var tenantID int64
	err := tx.QueryRowContext(
		ctx,
		query,
		dto.TenantName,
		dto.TenantCode,
		dto.Address,
		dto.CreatedBy,
		dto.UpdatedBy,
	).Scan(&tenantID)

	if err != nil {

		// Handle duplicate key error
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return 0, errors.New("role already exists")
			}
		}

		return 0, err
	}
	return tenantID, nil
}

// 6. Check the tenant code in DB

func (s *Store) TenantCodeInDB(ctx context.Context, tenantCode string) (bool, error) {

	query := `
		SELECT EXISTS(
    SELECT 1
    FROM tenant
    WHERE tenant_code = LOWER($1)
)
`

	var exists bool

	err := s.db.QueryRowContext(ctx, query, tenantCode).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil

}

// 7. Tenant verification store function
func (s *Store) TenantVerification(ctx context.Context, tenantCode string) (bool, error) {

	query := `
	UPDATE tenant
	SET is_verified = TRUE
	WHERE LOWER(tenant_code) = LOWER($1)
  AND is_deleted = FALSE
`

	result, err := s.db.ExecContext(ctx, query, tenantCode)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	// ✅ success check
	if rows == 0 {
		return false, nil // not found / no update
	}

	return true, nil
}

func (s *Store) GetTenantbyCode(ctx context.Context, tenantCode string) (*TenantResponse, error) {

	query := `
		SELECT id,
		tenant_name,
		tenant_code,
		contact_person_name,
		contact_phone,
		contact_email,
		address,
		is_verified,
		is_active,
		is_deleted,
		created_by,
		updated_by,
		created_at,
		updated_at
		FROM tenant
		WHERE tenant_code = LOWER($1)
	`

	var t TenantResponse

	err := s.db.QueryRowContext(ctx, query, tenantCode).Scan(
		&t.ID,
		&t.TenantName,
		&t.TenantCode,
		&t.ContactPersonName,
		&t.ContactPhone,
		&t.ContactEmail,
		&t.Address,
		&t.IsVerified,
		&t.IsActive,
		&t.IsDeleted,
		&t.CreatedBy,
		&t.UpdatedBy,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tenant not found")
		}
		return nil, err
	}

	return &t, nil
}

// 8. Delete Tenant by Code

func (s *Store) DeleteTenant(ctx context.Context, tenantCode string, deletedBy int64) (bool, error) {

	// Normalize input (important)

	tenantCode = strings.ToLower(strings.TrimSpace(tenantCode))

	query := `
	UPDATE tenant
	SET is_deleted = TRUE,
		deleted_at = NOW(),
		deleted_by = $2
	WHERE LOWER(tenant_code) = $1
		AND is_deleted = FALSE
	`

	result, err := s.db.ExecContext(ctx, query, tenantCode, deletedBy)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, ErrTenantNotDeleted
	}

	// ✅ success check
	if rows == 0 {
		return false, ErrTenantNotDeleted
	}

	return true, nil
}

// 9. Restore Deleted tenant

func (s *Store) RecoveryTenant(ctx context.Context, tenantCode string, deletedBy int64) (bool, error) {

	// Normalize input (important)

	tenantCode = strings.ToLower(strings.TrimSpace(tenantCode))

	query0 := `
		SELECT EXISTS (
		SELECT 1
		FROM tenant
		WHERE tenant_code = $1
)
`
	var exist bool
	err := s.db.QueryRowContext(ctx, query0, tenantCode).Scan(&exist)

	if err != nil {
		return false, ErrSameTenantCode
	}
	if exist {
		return false, ErrSameTenantCode
	}

	query := `
	UPDATE tenant
	SET is_deleted = TRUE,
		deleted_at = NOW(),
		deleted_by = $2
	WHERE LOWER(tenant_code) = $1
		AND is_deleted = FALSE
	`

	result, err := s.db.ExecContext(ctx, query, tenantCode, deletedBy)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, ErrTenantNotDeleted
	}

	// ✅ success check
	if rows == 0 {
		return false, ErrTenantNotDeleted
	}

	return true, nil
}

// 10. Update the Tenant by Tenant code

func (s *Store) UpdateTenant(ctx context.Context, t *Tenant) (bool, error) {

	query := `
		UPDATE tenant
		SET tenant_name = $1,
		    contact_person_name = $2,
		    contact_phone = $3,
		    contact_email = $4,
		    address = $5,
		    is_active = $6,
		    updated_by = $7,
		    updated_at = NOW()
		WHERE tenant_code = $8
		  AND is_deleted = false
	`

	result, err := s.db.ExecContext(ctx, query,
		t.TenantName,
		t.ContactPersonName,
		t.ContactPhone,
		t.ContactEmail,
		t.Address,
		t.IsActive,
		t.UpdatedBy,
		t.TenantCode,
	)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsAffected == 0 {
		return false, fmt.Errorf("%w: tenant_code=%s", ErrTenantNotUpdated, t.TenantCode)
	}

	return true, nil
}

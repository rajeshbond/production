package user

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

func (s *Store) CreateUser(ctx context.Context, dto UserDTO) (*User, error) {

	query := `
INSERT INTO "user" (
    tenant_id,
    role_id,
    employee_id,
    user_name,
    phone,
    email,
    password,
    is_verified,
    is_active,
    created_by,
    updated_by
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$10)
RETURNING 
    id,
    tenant_id,
    role_id,
    employee_id,
    user_name,
    phone,
    email,
    is_verified,
    is_active,
    created_by,
    updated_by,
    created_at,
    updated_at
`

	var user User

	err := s.db.QueryRowContext(ctx, query,
		dto.TenantID,
		dto.RoleID,
		dto.EmployeeID,
		dto.UserName,
		dto.Phone,
		dto.Email,
		dto.Password,
		dto.IsVerified,
		dto.IsActive,
		dto.CreatedBy,
	).Scan(
		&user.ID,
		&user.TenantID,
		&user.RoleID,
		&user.EmployeeID,
		&user.UserName,
		&user.Phone,
		&user.Email,
		&user.IsVerified,
		&user.IsActive,
		&user.CreatedBy,
		&user.UpdatedBy,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (s *Store) GetUserDetailByID(ctx context.Context, userID int64) (*User, error) {

	query := `
SELECT 
    id,
    tenant_id,
    role_id,
    employee_id,
    user_name,
    phone,
    email,
    is_verified,
    is_active,
    created_by,
    updated_by,
    created_at,
    updated_at
FROM "user"
WHERE id = $1
`

	var user User

	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.TenantID,
		&user.RoleID,
		&user.EmployeeID,
		&user.UserName,
		&user.Phone,
		&user.Email,
		&user.IsVerified,
		&user.IsActive,
		&user.CreatedBy,
		&user.UpdatedBy,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

// Get All Users by Tenant

func (s *Store) GetUsersByTenantID(ctx context.Context, tenantID int64) ([]User, error) {

	query := `
SELECT
    id,
    tenant_id,
    role_id,
    employee_id,
    user_name,
    phone,
    email,
    is_verified,
    is_active,
    created_by,
    updated_by,
    created_at,
    updated_at
FROM "user"
WHERE tenant_id = $1
ORDER BY id
`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID,
			&user.TenantID,
			&user.RoleID,
			&user.EmployeeID,
			&user.UserName,
			&user.Phone,
			&user.Email,
			&user.IsVerified,
			&user.IsActive,
			&user.CreatedBy,
			&user.UpdatedBy,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

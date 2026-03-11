package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rajesh_bond/production/internal/common/response"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// Create User

func (s *Store) CreateUser(ctx context.Context, dto UserCreateRequest) (*User, error) {

	query := `
	INSERT INTO "user" (
		tenant_id,
		role_id,
		employee_id,
		user_name,
		phone,
		email,
		password,
		created_by,
		updated_by
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
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

	err := s.db.QueryRowContext(
		ctx,
		query,
		dto.TenantID,
		dto.RoleId,
		dto.EmployeeID,
		dto.UserName,
		dto.Phone,
		dto.Email,
		dto.Password,
		dto.CreatedBy,
		dto.UpdatedBy,
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
		return nil, response.HandlePostgresError(err)
	}

	return &user, nil
}

// Get User detaiils by ID

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

// Get HashPassword by Employee ID

func (s *Store) GetPasswordHashbyEmplopeeID(ctx context.Context, employeeID string) (*UserPayload, string, error) {
	var passwordHash string

	query := `
		SELECT id, 
		tenant_id,
	  user_name,
	  role_id,
	  password
		FROM "user"
		WHERE employee_id = $1
	`
	payload := &UserPayload{}

	err := s.db.QueryRowContext(ctx, query, employeeID).Scan(
		&payload.UserID,
		&payload.TenantID,
		&payload.Username,
		&payload.RoleID,
		&passwordHash,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", errors.New("user not found")
		}
		return nil, "", err
	}
	return payload, passwordHash, nil
}

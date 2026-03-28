package users

/*
// Users Index
//////////////////////////////////
// 1.Create User
//2. Get User detaiils by ID
//3. Get All Users by Tenant
//4. Get HashPassword by Employee ID
//5. Create Super Admin
//6.create Tenant Users
//7. Is Employee Exist
//8. Is Tenant Exist
//9. Get Verify Tenant User
//10. Get Verification Status
//11. Get Tenant ID by Code
// //12. Delete Tenant User

//////////////////////////////////
*/

// imports
import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/rajesh_bond/production/internal/common/response"
)

// Store Stuct
type Store struct {
	db *sql.DB
}

// Store constructor
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// 1.Create User

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
		dto.RoleID,
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

//2. Get User detaiils by ID

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
		is_deleted,
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
		&user.IsDeleted,
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

//3. Get All Users by Tenant

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

// 4. Get HashPassword by Employee ID
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

// 5. Create Super Admin
func (s *Store) CreateSuperAdminTx(ctx context.Context, tx *sql.Tx, dto UserCreateRequest) (int64, error) {

	query := `
		INSERT INTO "user" 
		(tenant_id, role_id, employee_id, user_name, phone, email, password, created_by, updated_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id
	`

	var id int64

	err := tx.QueryRowContext(
		ctx,
		query,
		dto.TenantID,
		dto.RoleID,
		dto.EmployeeID,
		dto.UserName,
		dto.Phone,
		dto.Email,
		dto.Password,
		dto.CreatedBy,
		dto.UpdatedBy,
	).Scan(&id)

	if err != nil {

		// Handle duplicate key error
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return 0, errors.New("user already exists")
			}
		}

		return 0, err
	}

	return id, nil
}

// 6.create Tenant Users
func (s *Store) CreateTenantUser(ctx context.Context, dto *UserCreateRequest) (*CreateUserResponse, error) {

	query := `
	INSERT INTO "user"
	(tenant_id, role_id, employee_id, user_name, phone, email, password, created_by, updated_by)
	VALUES
	($1,$2,$3,$4,$5,$6,$7,$8,$9)
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
		is_deleted,
		deleted_by,
		created_by,
		updated_by
	`

	var resp CreateUserResponse

	// nullable fields
	var phone sql.NullString
	var email sql.NullString
	var deletedBy sql.NullInt64

	err := s.db.QueryRowContext(
		ctx,
		query,
		dto.TenantID,
		dto.RoleID,
		dto.EmployeeID,
		dto.UserName,
		dto.Phone,
		dto.Email,
		dto.Password,
		dto.CreatedBy,
		dto.UpdatedBy,
	).Scan(
		&resp.ID,
		&resp.TenantID,
		&resp.RoleID,
		&resp.EmployeeID,
		&resp.UserName,
		&phone,
		&email,
		&resp.IsVerified,
		&resp.IsActive,
		&resp.IsDeleted, // ✅ direct bool
		&deletedBy,
		&resp.CreatedBy,
		&resp.UpdatedBy,
	)

	if err != nil {
		return nil, err
	}

	// ✅ Handle nullable fields
	if phone.Valid {
		resp.Phone = &phone.String
	}
	if email.Valid {
		resp.Email = &email.String
	}
	if deletedBy.Valid {
		resp.DeletedBy = &deletedBy.Int64
	}

	return &resp, nil
}

// 7. Is Employee Exist
func (s *Store) IsEmployeeExist(ctx context.Context, employeeID string, tenantID int64) (bool, error) {

	query := `
		SELECT EXISTS(
		SELECT 1
		FROM "user"
		WHERE employee_id = $1
		AND tenant_id = $2
		AND is_deleted = false
		)`
	var exists bool

	err := s.db.QueryRowContext(ctx, query, employeeID, tenantID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil

}

// 8. Is Tenant Exist
func (s *Store) IsTenantExist(ctx context.Context, tennatCode string) (bool, error) {
	query := `
		SELECT EXISITS(
		SELECT 1,
		FROM tenant,
		WHERE tenant_code = $1
		AND is_deleted = false
	)`

	var exists bool

	err := s.db.QueryRowContext(ctx, query, tennatCode).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil

}

// 9. Get Verify Tenant User
func (s *Store) GetVerifyTenantUser(ctx context.Context, employeeID string, tenantID int64) (bool, error) {

	query := `
	UPDATE "user"
	SET is_verified = true 
	WHERE employee_id = $1
	AND tenant_id = $2
	AND is_deleted = false
	`

	result, err := s.db.ExecContext(ctx, query, employeeID, tenantID)
	if err != nil {
		return false, err
	}

	resultRowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, nil
	}

	// If no rows updated -> not found or already deleetd

	if resultRowsAffected == 0 {
		return false, err
	}

	return true, nil

}

// 10. Get Verification Status
func (s *Store) GetVerificationStatus(ctx context.Context, employeeID string, tenantID int64) (bool, bool, error) {

	query := `
		SELECT is_verified
		FROM "user"
		WHERE employee_id = $1
		  AND tenant_id = $2
		  AND is_deleted = false
	`

	var isVerified bool

	err := s.db.QueryRowContext(ctx, query, employeeID, tenantID).Scan(&isVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, false, nil // not found
		}
		return false, false, err
	}

	return true, isVerified, nil
}

// 11. Get Tenant ID by Code
func (s *Store) GetTenantIDByCode(ctx context.Context, tenantName string) (int64, error) {
	query := `
		SELECT id
		FROM tenant
		WHERE LOWER(tenant_code) = LOWER($1)
	`

	var tenantID int64

	err := s.db.QueryRowContext(ctx, query, tenantName).Scan(&tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("Tenant Not Found")
		}

		return 0, err
	}

	return tenantID, nil
}

// 12. Delete Tenant User
func (s *Store) DeleteTenantUser(ctx context.Context, employeeID string, tenantID int64, userID int64) (bool, error) {
	query := `
	UPDATE "user"
	SET is_deleted = TRUE,
		deleted_at = NOW(),
		deleted_by = $3,
		updated_by = $3,
		updated_at = NOW()
	WHERE employee_id = $1
	  AND tenant_id = $2
	  AND is_deleted = false
	`

	result, err := s.db.ExecContext(ctx, query, employeeID, tenantID, userID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rows == 0 {
		return false, errors.New("user not found or already deleted")
	}

	return true, nil
}
func (s *Store) GetUserbyEmploeeID(ctx context.Context, employeeID string, tenantID int64) (*CreateUserResponse, error) {

	query := `
		SELECT id,
			employee_id,
			tenant_id,
			user_name,
			email,
			phone,
			is_verified,
			is_active,
			is_deleted,
			created_by,
			updated_by
		FROM "user"
		WHERE employee_id = $1
		AND tenant_id = $2
		AND is_deleted = false
	`

	var resp CreateUserResponse

	err := s.db.QueryRowContext(ctx, query, employeeID, tenantID).Scan(
		&resp.ID,
		&resp.EmployeeID,
		&resp.TenantID,
		&resp.UserName,
		&resp.Email,
		&resp.Phone,
		&resp.IsVerified,
		&resp.IsActive,
		&resp.IsDeleted,
		&resp.CreatedBy,
		&resp.UpdatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &resp, nil

}

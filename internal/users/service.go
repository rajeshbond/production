package users

/*
///////////////////////////////////////
// Index - Service Layer
///////////////////////////////////////
// 1. Create User
// 2. Login
// 3. Create Super User
// 4. Create Tenant User
// 5. Check Employee
// 6. Check Tenant
// 7. Verify Tenant User
// 8. Delete Tenant User




///////////////////////////////////////
*/
import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/rajesh_bond/production/cmd/service"
	"github.com/rajesh_bond/production/internal/auth"
	"github.com/rajesh_bond/production/internal/common/utils"
)

type Service struct {
	Store          *Store
	RoleProvider   RoleProvider
	TenantProvider TenantProvider
}

func NewService(store *Store, roleProvider RoleProvider, tenantProvider TenantProvider) *Service {
	return &Service{
		Store:          store,
		RoleProvider:   roleProvider,
		TenantProvider: tenantProvider,
	}
}

// 1. Create User - Corrected
func (ser *Service) CreateTenantAdmin(ctx context.Context, claims *auth.UserClaims, req UserCreateRequest) (*UserResponse, error) {

	// Auth check (who can create the Tenant Admin)

	// 1. Validate request
	if err := utils.Validate.Struct(req); err != nil {
		return nil, fmt.Errorf("%w,%v", ErrInvalidRequest, err)
	}
	// 2. Validate employee_id format (extra safety if needed)
	if !strings.Contains(req.EmployeeID, "@") {
		return nil, ErrInvalidRequest
	}
	// 1. Validate request

	if err := utils.Validate.Struct(req); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidRequest, err)
	}

	// 2.Validate empolyee_id format (extra saftey of needed)

	if !strings.Contains(req.EmployeeID, "@") {
		return nil, ErrInvalidRequest
	}

	tenandCode, err := auth.Tcode(req.EmployeeID)
	// 3. check Tenant status

	employeeTenantID, err := ser.TenantProvider.GetTenantIDByCode(ctx, tenandCode)
	if err != nil {
		return nil, err
	}

	if req.TenantID != employeeTenantID {
		return nil, ErrTenantIDMismatched
	}

	isVerified, isActive, isDeleted, err := ser.TenantProvider.GetTenantStatus(ctx, tenandCode)
	if err != nil {
		return nil, err
	}

	if isDeleted {
		return nil, ErrTenantDeleted
	}

	if !isActive {
		return nil, ErrTenantInActive
	}
	if !isVerified {
		return nil, ErrTenantVerified
	}

	reqRole, err := ser.RoleProvider.GetRoleNameByID(ctx, req.RoleID)
	if err != nil {
		return nil, err
	}

	// if reqRole != "tenantadmin" && reqRole != "tenantowner" {
	// 	return nil, ErrOnlyTenantAdminCreate
	// }

	if !auth.IsTenatAdminRole(reqRole) {
		return nil, ErrOnlyTenantAdminCreate
	}

	// continue normal flow

	// Asining the Created and updated by to struct
	req.CreatedBy = &claims.UserID
	req.UpdatedBy = &claims.UserID

	// Check the role is Present in the DB

	// 4. check User already exsits (options but recommended)
	exists, err := ser.Store.IsEmployeeExist(ctx, req.EmployeeID, req.TenantID)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrUserAlreadyExistForThisTenant
	}

	// 5. Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	req.Password = hashedPassword

	// Call store
	user, err := ser.Store.CreateUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create user failed: %w", err)
	}

	// 🔹 7. Map to response
	res := &UserResponse{
		ID:         user.ID,
		TenantID:   user.TenantID,
		RoleID:     user.RoleID,
		EmployeeID: user.EmployeeID,
		UserName:   user.UserName,
		Phone:      utils.SafeString(user.Phone),
		Email:      utils.SafeString(user.Email),
		IsVerified: user.IsVerified,
		IsActive:   user.IsActive,
		CreatedBy:  utils.SafeInt(user.CreatedBy),
		UpdatedBy:  utils.SafeInt(user.UpdatedBy),
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}

	return res, nil

}

// 2. Login - Corrected
func (ser *Service) LoginUser(ctx context.Context, req LoginRequest) (*LoginResponse, error) {

	// validate request
	if err := utils.Validate.Struct(req); err != nil {
		return nil, err
	}

	tcode, err := auth.Tcode(req.EmployeeID)
	if err != nil {
		return nil, err
	}

	tenantID, err := ser.TenantProvider.GetTenantIDByCode(ctx, tcode)
	if err != nil {
		return nil, err
	}

	// ✅ Check status
	found, isVerified, err := ser.Store.GetVerificationStatus(ctx, req.EmployeeID, tenantID)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, errors.New("user not found")
	}

	if !isVerified {
		return nil, errors.New("User Not Verified please contact Admin")
	}

	// fetch user data + password
	tokenPayload, hashedPassword, err := ser.Store.GetPasswordHashbyEmplopeeID(ctx, req.EmployeeID)
	if err != nil {
		return nil, err
	}

	// compare password
	if err := utils.CompareHash(hashedPassword, req.Password); err != nil {
		return nil, err
	}

	role, err := ser.RoleProvider.GetRoleNameByID(ctx, tokenPayload.RoleID)
	if err != nil {
		return nil, err
	}

	// prepare jwt payload
	payload := service.TokenPayload{
		TenantID: tokenPayload.TenantID,
		UserID:   tokenPayload.UserID,
		Username: tokenPayload.Username,
		RoleID:   tokenPayload.RoleID,
		Role:     role,
	}

	tokenString, err := service.GenerateToken(payload, req.EmployeeID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		// UserID: tokenPayload.UserID,
		Token: tokenString}, nil
}

// 3. Create Super User - Corrected
func (s *Service) CreateSuperUserTx(ctx context.Context, tx *sql.Tx, tenantID int64, roleID int64, dto UserSuperRequest) (int64, error) {
	createdBy := int64(1)
	hasshedPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		return 0, err
	}
	req := UserCreateRequest{
		TenantID:   tenantID,
		RoleID:     roleID,
		EmployeeID: dto.EmployeeID,
		UserName:   dto.UserName,
		Phone:      dto.Phone,
		Email:      dto.Email,
		Password:   hasshedPassword,
		CreatedBy:  &createdBy,
		UpdatedBy:  &createdBy,
	}

	return s.Store.CreateSuperAdminTx(ctx, tx, req)

}

// 4. Create Tenant User
func (s *Service) CreateTenantUser(ctx context.Context, claims *auth.UserClaims, req *UserCreateRequest) (*CreateUserResponse, error) {

	// Basic validation
	if strings.TrimSpace(req.EmployeeID) == "" {
		return nil, ErrEmployeeIDReqyured
	}

	// Authorization

	if err := auth.ValidateTenantAccess(
		claims.Role,
		claims.EmployeeID,
		req.EmployeeID,
	); err != nil {
		return nil, err
	}

	reqRole, err := s.RoleProvider.GetRoleNameByID(ctx, req.RoleID)

	if err != nil {
		return nil, err
	}

	fmt.Println("role check", reqRole)

	err = auth.TenantRoleCheck(reqRole)

	if err != nil {
		return nil, err
	}

	// Duplicate check

	exists, err := s.Store.IsEmployeeExist(ctx, req.EmployeeID, req.TenantID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExistForThisTenant
	}

	// Create user

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	req.Password = hashedPassword

	req.CreatedBy = &claims.UserID
	req.UpdatedBy = &claims.UserID

	user, err := s.Store.CreateTenantUser(ctx, req)

	if err != nil {
		return nil, err
	}

	return user, nil

}

// 5. Check Employee
func (s *Service) CheckEmployeeExist(ctx context.Context, employeeID string, tenantID int64) error {
	// Basic Validation
	if strings.TrimSpace(employeeID) == "" {
		return ErrEmployeeIDRequired
	}

	// call Store

	exists, err := s.Store.IsEmployeeExist(ctx, employeeID, tenantID)
	if err != nil {
		return err
	}
	if exists {
		return ErrUserAlreadyExists
	}

	return nil

}

// 6. Check Tenant
func (ser *Service) CheckTenantExist(ctx context.Context, tenantCode string) error {
	exists, err := ser.Store.IsTenantExist(ctx, tenantCode)
	if err != nil {
		return err
	}

	if exists {
		return ErrAlreadyTenantPresent
	}

	return nil

}

// 7. Verify Tenant User
func (s *Service) VerifyTenantUser(ctx context.Context, claims *auth.UserClaims, employeeID string, tenantID int64) error {

	fmt.Println("Claims Role ", claims.Role)
	user, err := s.Store.GetUserbyEmploeeID(ctx, employeeID, tenantID)
	if err != nil {
		return err
	}
	userRole, err := s.RoleProvider.GetRoleNameByID(ctx, user.RoleID)

	if err != nil {
		return err
	}

	if user.IsDeleted {
		return ErrUserDeleted
	}

	if user.IsVerified {
		return errors.New("user already verified")
	}

	// ✅ Auth check
	if err := auth.ValidateTenantAccess(
		claims.Role,
		claims.EmployeeID,
		employeeID,
	); err != nil {
		return err
	}

	// ✅ Basic validation

	if (claims.Role == "superadmin" || claims.Role == "admin") || (userRole == "tenantadmin" || claims.Role == "tenantowner") {
		fmt.Println("tenant admin")
		updated, err := s.Store.VerifyTenantUser(ctx, employeeID, tenantID, claims.UserID)
		if err != nil {
			return err
		}

		if !updated {
			return errors.New("verification failed")
		}

		return nil
	}

	if claims.Role == "tenantadmin" {

		if claims.TenantID != user.TenantID {
			return ErrTenantIDMismatched
		}
		updated, err := s.Store.VerifyTenantUser(ctx, employeeID, tenantID, claims.UserID)
		if err != nil {
			return err
		}

		if !updated {
			return errors.New("verification failed")
		}

		return nil
	}

	return ErrUnauthorized
}

// 8. Delete Tenant User
func (ser *Service) DeleteTenantUser(ctx context.Context, claims *auth.UserClaims, employeeID string, tenantID int64) error {

	// ✅ Basic validation
	if employeeID == "" {
		return errors.New("employee_id is required")
	}

	if tenantID == 0 {
		return errors.New("tenant_id is required")
	}

	user, err := ser.Store.GetUserbyEmploeeID(ctx, employeeID, tenantID)
	if err != nil {
		return err
	}
	userRole, err := ser.RoleProvider.GetRoleNameByID(ctx, user.RoleID)

	if err != nil {
		return err
	}

	if user.IsDeleted {
		return ErrUserDeleted
	}

	// ✅ Auth check
	if err := auth.ValidateTenantAccess(
		claims.Role,
		claims.EmployeeID,
		employeeID,
	); err != nil {
		return err
	}

	// ✅ Basic validation

	if (claims.Role == "superadmin" || claims.Role == "admin") && (userRole == "tenantadmin") {

		deleted, err := ser.Store.DeleteTenantUser(ctx, employeeID, tenantID, claims.UserID)
		if err != nil {
			return err
		}

		if !deleted {
			return errors.New("failed to delete user")
		}

		return nil
	}

	if claims.Role == "tenantadmin" {

		if claims.TenantID != user.TenantID {
			return ErrTenantIDMismatched
		}

		deleted, err := ser.Store.DeleteTenantUser(ctx, employeeID, tenantID, claims.UserID)
		if err != nil {
			return err
		}

		if !deleted {
			return errors.New("failed to delete user")
		}

		return nil
	}

	return ErrUnauthorized

}

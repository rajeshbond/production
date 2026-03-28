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
	RoleProvide    RoleProvide
	TenantProvider TenantProvider
}

func NewService(store *Store, roleProvider RoleProvide, tenantProvider TenantProvider) *Service {
	return &Service{
		Store:          store,
		RoleProvide:    roleProvider,
		TenantProvider: tenantProvider,
	}
}

// 1. Create User
func (ser *Service) CreateUser(ctx context.Context, req UserCreateRequest) (*UserResponse, error) {

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	req.Password = hashedPassword

	// Call store
	user, err := ser.Store.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	// convert DB model -> Response DTO

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

// 2. Login
func (ser *Service) LoginUser(ctx context.Context, req LoginRequest) (*LoginResponse, error) {

	// validate request
	if err := utils.Validate.Struct(req); err != nil {
		return nil, err
	}

	tcode, err := auth.Tcode(req.EmployeeID)
	if err != nil {
		return nil, err
	}

	fmt.Println("Tenant Code:-", tcode)

	tenantID, err := ser.TenantProvider.GetTenantIDByCode(ctx, tcode)
	if err != nil {
		return nil, err
	}

	// tenantID, err := ser.Store.GetTenantIDByCode(ctx, tcode)
	// if err != nil {
	// 	return nil, err
	// }

	// ✅ Check status
	found, isVerified, err := ser.Store.GetVerificationStatus(ctx, req.EmployeeID, tenantID)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, errors.New("user not found")
	}

	if isVerified {
		return nil, errors.New("user already verified")
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

	role, err := ser.RoleProvide.GetRoleNameByID(ctx, tokenPayload.RoleID)
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

// 3. Create Super User
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

	// ✅ Auth check
	if err := auth.ValidateTenantAccess(
		claims.Role,
		claims.EmployeeID,
		employeeID,
	); err != nil {
		return err
	}

	// ✅ Check status
	found, isVerified, err := s.Store.GetVerificationStatus(ctx, employeeID, tenantID)
	if err != nil {
		return err
	}

	if !found {
		return errors.New("user not found")
	}

	if isVerified {
		return errors.New("user already verified")
	}

	// ✅ Update
	updated, err := s.Store.GetVerifyTenantUser(ctx, employeeID, tenantID)
	if err != nil {
		return err
	}

	if !updated {
		return errors.New("verification failed")
	}

	return nil
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

	// ✅ Auth check (who can delete whom)
	if err := auth.ValidateTenantAccess(claims.Role, claims.EmployeeID, employeeID); err != nil {
		return err
	}

	// ✅ Fetch user
	cUser, err := ser.Store.GetUserbyEmploeeID(ctx, employeeID, tenantID)
	if err != nil {
		return err
	}

	if cUser == nil {
		return errors.New("user not found")
	}

	// ✅ Already deleted check
	if cUser.IsDeleted {
		return errors.New("user already deleted")
	}

	// ✅ Soft delete
	deleted, err := ser.Store.DeleteTenantUser(ctx, employeeID, tenantID, claims.UserID)
	if err != nil {
		return err
	}

	if !deleted {
		return errors.New("failed to delete user")
	}

	return nil
}

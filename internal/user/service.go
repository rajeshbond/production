package user

import (
	"context"

	"github.com/rajesh_bond/production/internal/common/utils"
	"github.com/rajesh_bond/production/internal/tenant"
	userrole "github.com/rajesh_bond/production/internal/user_role"
)

type Service struct {
	store       *Store
	roleStore   *userrole.Store
	tenantStore *tenant.Store
}

func NewService(store *Store, roleStore *userrole.Store, tenantStore *tenant.Store) *Service {
	return &Service{
		store:       store,
		roleStore:   roleStore,
		tenantStore: tenantStore,
	}
}

func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponseDTO, error) {

	roleID, err := s.roleStore.GetRoleIDByName(ctx, req.Role)
	if err != nil {
		return nil, err
	}

	tenantID, err := s.tenantStore.GetTenantIDByName(ctx, req.TenantName)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	isVerified := false
	isActive := true

	dto := UserDTO{
		TenantID:   &tenantID,
		RoleID:     &roleID,
		EmployeeID: &req.EmployeeID,
		UserName:   &req.UserName,
		Password:   &hashedPassword,
		Phone:      &req.Phone,
		Email:      &req.Email,
		IsVerified: &isVerified,
		IsActive:   &isActive,
	}

	user, err := s.store.CreateUser(ctx, dto)
	if err != nil {
		return nil, err
	}

	resp := &UserResponseDTO{
		ID:         user.ID,
		TenantID:   user.TenantID,
		RoleID:     user.RoleID,
		UserName:   user.UserName,
		Phone:      user.Phone,
		Email:      user.Email,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
	}

	return resp, nil
}

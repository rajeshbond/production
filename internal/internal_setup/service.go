package internalsetup

import (
	"context"
	"database/sql"

	"github.com/rajesh_bond/production/internal/tenant"
	userrole "github.com/rajesh_bond/production/internal/user_role"
	"github.com/rajesh_bond/production/internal/users"
)

type Service struct {
	db            *sql.DB
	tenantService *tenant.Service
	roleService   *userrole.Service
	userService   *users.Service
}

func NewService(db *sql.DB,
	tenantService *tenant.Service,
	roleService *userrole.Service,
	userService *users.Service) *Service {
	return &Service{
		db:            db,
		tenantService: tenantService,
		roleService:   roleService,
		userService:   userService,
	}
}

func (s *Service) SetupSuperAdmin(ctx context.Context, dto *SetupSuperAdminDTO) (*SetupResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create Role
	roleID, err := s.roleService.CreateRoleTx(ctx, tx, dto.Role.UserRole)

	// Create Tenant
	tenantID, err := s.tenantService.CreateSuperTenantTx(ctx, tx, tenant.CreateTenantRequied(dto.Tenant))

	// Create User
	var phonePtr *string
	var emailPtr *string

	if dto.User.Phone != "" {
		phonePtr = &dto.User.Phone
	}

	if dto.User.Email != "" {
		emailPtr = &dto.User.Email
	}

	userreq := users.UserSuperRequest{
		EmployeeID: dto.User.EmployeeID,
		UserName:   dto.User.UserName,
		Phone:      phonePtr,
		Email:      emailPtr,
		Password:   dto.User.Password,
	}
	userID, err := s.userService.CreateSuperUserTx(ctx, tx, tenantID, roleID, userreq)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &SetupResponse{
		RoleID:   roleID,
		TenantID: tenantID,
		UserID:   userID,
		Message:  "Sucessfully",
	}, nil

}

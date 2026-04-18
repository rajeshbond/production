package userrole

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (s *Service) Create(
	ctx context.Context,
	dto CreateRole, userID int64,
) (*UserRoleResponseDTO, error) {

	fmt.Println("Service", dto)

	sendrole := dto.UserRole

	if sendrole == "" {
		return nil, fmt.Errorf("user_role is required")
	}

	exist, err := s.store.RoleInDB(ctx, sendrole)

	if err != nil {
		return nil, err
	}

	if exist {
		return nil, fmt.Errorf("role already exists")
	}

	// role, err := s.store.Create(ctx, dto)
	role, err := s.store.Create(ctx, dto, userID)
	if err != nil {
		return nil, err
	}

	return &UserRoleResponseDTO{
		ID:        role.ID,
		UserRole:  role.UserRole,
		CreatedBy: role.CreatedBy,
		UpdatedBy: role.UpdatedBy,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}, nil
}

// This is for intial route only allowed in developemt phase
func (s *Service) CreateRoleTx(ctx context.Context, tx *sql.Tx, role string) (int64, error) {

	createdBy := int64(1)

	dto := CreateUserRoleDTO{
		UserRole:  strings.ToLower(role),
		CreatedBy: &createdBy,
		UpdatedBy: &createdBy,
	}

	return s.store.CreateRoleSuperTx(ctx, tx, dto)

}

func (s *Service) GetRoleNameByID(ctx context.Context, roleID int64) (string, error) {

	return s.store.GetRoleNameByID(ctx, roleID)
}

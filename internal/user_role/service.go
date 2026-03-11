package userrole

import (
	"context"
	"database/sql"
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
	dto CreateUserRoleDTO,
) (*UserRoleResponseDTO, error) {

	// fmt.Printf("DTO:%+V\n", dto)
	dbDTO := CreateUserRoleDBDTO{
		UserRole:  strings.ToLower(dto.UserRole),
		CreatedBy: dto.CreatedBy,
	}

	role, err := s.store.Create(ctx, dbDTO)
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

func (s *Service) CreateRoleTx(ctx context.Context, tx *sql.Tx, role string) (int64, error) {

	createdBy := int64(1)

	dto := CreateUserRoleDTO{
		UserRole:  strings.ToLower(role),
		CreatedBy: &createdBy,
		UpdatedBy: &createdBy,
	}

	return s.store.CreateRoleSuperTx(ctx, tx, dto)

}

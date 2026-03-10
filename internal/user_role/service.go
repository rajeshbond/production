package userrole

import (
	"context"
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

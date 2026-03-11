package tenant

import (
	"context"
	"database/sql"

	"github.com/rajesh_bond/production/internal/common/response"
)

type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (s *Service) CreateTenant(ctx context.Context, dto CreateTenantDTO) (*TenantResponse, error) {

	// example Validate

	if dto.TenantName == "" {
		return nil, response.ErrTeantNameRequired
	}
	if dto.TenantCode == "" {

		return nil, response.ErrTenantCondeRequired
	}

	tenant, err := s.store.Create(ctx, dto)
	if err != nil {
		return nil, err
	}

	TenantResponse := &TenantResponse{
		ID:         tenant.ID,
		TenantName: tenant.TenantName,
		TenantCode: tenant.TenantCode,
		Address:    tenant.Address,
		IsVerified: tenant.IsVerified,
		IsActive:   tenant.IsActive,
		CreatedBy:  tenant.CreatedBy,
		UpdatedBy:  tenant.UpdatedBy,
		CreatedAt:  tenant.CreatedAt,
		UpdatedAt:  tenant.UpdatedAt,
	}

	return TenantResponse, nil

}

func (s *Service) CreateSuperTenantTx(ctx context.Context, tx *sql.Tx, dto CreateTenantRequied) (int64, error) {

	req := CreateTenantDTO{
		TenantName: dto.TenantName,
		TenantCode: dto.TenantCode,
		Address:    dto.Address,
		CreatedBy:  int64(1),
		UpdatedBy:  int64(1),
	}

	return s.store.CreateSuperTenantTx(ctx, tx, req)

}

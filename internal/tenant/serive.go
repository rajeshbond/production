package tenant

// Index - Tenant Service Layer
/////////////////////////////////

// 1. Create Tenant

// 2. Create Tenant for Dev section only

// 3. Tenant Verifcation toggle

/////////////////////////////////

// Code starts ---------->>>
/////////////////////////////////

// import

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rajesh_bond/production/internal/common/response"
)

// define the servive layer sturts

type Service struct {
	store *Store
}

// define struct constructor

func NewService(store *Store) *Service {
	return &Service{store: store}
}

// 1. Create Tenant

func (s *Service) CreateTenant(ctx context.Context, dto CreateTenantDTO) (*TenantResponse, error) {

	// example Validate

	if dto.TenantName == "" {
		return nil, ErrTenantNameRequired
	}
	if dto.TenantCode == "" {

		return nil, ErrTenantCodeRequired
	}

	if dto.Address == "" {
		return nil, ErrTenantAddressRequired
	}

	exists, err := s.store.TenantCodeInDB(ctx, dto.TenantCode)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrTenantCodeExists
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

// 2. Create Tenant for Dev section only

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

// 3. Tenant Verifcation toggle

func (ser *Service) TenantVerifcation(ctx context.Context, tenantCode string) (string, error) {

	tenant, err := ser.store.GetTenantbyCode(ctx, tenantCode)
	fmt.Println("tenant:-", tenant.IsDeleted)
	if err != nil {
		return "Intenal server Error", err
	}

	if tenant.IsDeleted {
		return TenantDeleted, nil
	}

	if tenant.IsVerified {
		return TenantAlreadyVerifed, nil
	}

	TenantVerifed, err := ser.store.TenantVerification(ctx, tenantCode)

	if err != nil {
		return "Intenal server Error", err
	}

	if TenantVerifed {
		return TenantVerified, nil

	}

	return TenantVerified, nil

}

// 4. Delete Tenant

func (ser *Service) DeleteTenant(ctx context.Context, tenantCode string, deletedBy int64) (string, error) {

	tenant, err := ser.store.GetTenantbyCode(ctx, tenantCode)
	if err != nil {
		return response.InternalServerError, err
	}

	if tenant.IsDeleted {
		return TenantDeleted, nil
	}

	deleted, err := ser.store.DeleteTenant(ctx, tenantCode, deletedBy)
	if err != nil {
		return response.InternalServerError, err
	}

	if !deleted {
		return TenantNotDeleted, nil
	}

	return TenantNotDeleted, nil

}

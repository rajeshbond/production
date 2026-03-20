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
	"strings"
	"time"

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

// 5.Update Tenant

func (ser *Service) UpdateTenant(ctx context.Context, tenantCode string, dto UpdateTenantDTO) (bool, error) {

	// Normalize tenant code
	tenantCode = strings.ToLower(strings.TrimSpace(tenantCode))

	// fmt.Println("tenantCode:-", tenantCode)

	if tenantCode == "" {
		return false, ErrTenantCodeRequired
	}

	// 1. Get existing tenant
	existing, err := ser.store.GetTenantbyCode(ctx, tenantCode)
	if err != nil {
		return false, fmt.Errorf("fetch tenant failed: %w", err)
	}

	if existing == nil {
		return false, fmt.Errorf("%w: tenant_code=%s", ErrTenantCodeNotFount, tenantCode)
	}
	if existing.IsDeleted {
		return false, ErrTenantDeletedInPast
	}

	// 2. Apply updates
	isChanged := false

	if dto.TenantName != nil {
		val := strings.TrimSpace(*dto.TenantName)
		if val != "" && val != existing.TenantName {
			existing.TenantName = val
			isChanged = true
		}
	}

	if dto.ContactPersonName != nil {
		val := strings.TrimSpace(*dto.ContactPersonName)
		if val != "" && val != *existing.ContactPersonName {
			existing.ContactPersonName = &val
			isChanged = true
		}
	}

	if dto.ContactEmail != nil {
		val := strings.TrimSpace(*dto.ContactEmail)
		if val != "" && val != *existing.ContactEmail {
			existing.ContactEmail = &val
			isChanged = true
		}
	}

	if dto.ContactPhone != nil {
		val := strings.TrimSpace(*dto.ContactPhone)
		if val != "" && val != *existing.ContactPhone {
			existing.ContactPhone = &val
			isChanged = true
		}
	}

	if dto.Address != nil {
		val := strings.TrimSpace(*dto.Address)
		if val != "" && val != *existing.Address {
			existing.Address = &val
			isChanged = true
		}
	}

	if dto.IsActive != nil {
		if *dto.IsActive != existing.IsActive {
			existing.IsActive = *dto.IsActive
			isChanged = true
		}
	}

	if dto.UpdatedBy != nil {
		existing.UpdatedBy = dto.UpdatedBy
		isChanged = true
	}

	// 3. No changes → skip DB
	if !isChanged {
		return false, ErrTenantNotUpdated
	}

	// 4. Update timestamp
	existing.UpdatedAt = time.Now().UTC()

	// 5. Call store

	existingTenant := &Tenant{
		ID:                existing.ID,
		TenantName:        existing.TenantName,
		TenantCode:        existing.TenantCode,
		ContactPersonName: existing.ContactPersonName,
		ContactPhone:      existing.ContactPhone,
		ContactEmail:      existing.ContactEmail,
		Address:           existing.Address,
		IsVerified:        existing.IsVerified,
		IsActive:          existing.IsActive,
		IsDeleted:         existing.IsDeleted,
		CreatedBy:         existing.CreatedBy,
		UpdatedBy:         existing.UpdatedBy,
		CreatedAt:         existing.CreatedAt,
		UpdatedAt:         existing.UpdatedAt,
	}

	updated, err := ser.store.UpdateTenant(ctx, existingTenant)
	if err != nil {
		return false, fmt.Errorf("update tenant failed: %w", err)
	}

	return updated, nil
}

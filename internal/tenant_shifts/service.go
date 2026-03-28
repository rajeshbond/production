package tenantshifts

import (
	"context"

	"github.com/rajesh_bond/production/internal/auth"
	"github.com/rajesh_bond/production/internal/common/utils"
)

type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (ser *Service) CreateTenantShift(ctx context.Context, claims *auth.UserClaims, req *CreateTenantShiftRequest) (*TenantShift, error) {
	// Validate struct
	if err := utils.Validate.Struct(req); err != nil {
		return nil, err
	}
	// Auth check
	// // 🔐 Always trust claims (override)
	// req.TenantID = claims.TenantID
	// req.CreatedBy = &claims.UserID

	if err := utils.Validate.Struct(req); err != nil {
		return nil, err
	}

	if err := auth.ValidateTenantAccesswithTenantCode(claims.Role, claims.TenantID, req.TenantID); err != nil {
		return nil, err
	}

	// if err := auth.ValidateTenantAccess(claims.Role, claims.EmployeeID, claims.EmployeeID); err != nil {
	// 	return nil, err
	// }

	return ser.store.CreateTenantShift(ctx, req)

}

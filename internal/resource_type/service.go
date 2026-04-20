package resourcetype

import (
	"context"

	"github.com/rajesh_bond/production/internal/auth"
)

type Service struct {
	Store *Store
}

func NewService(store *Store) *Service {
	return &Service{Store: store}
}

func (ser *Service) CreateResourceType(ctx context.Context, req CreateResourceTypeRequest, claims *auth.UserClaims) (int64, error) {
	return ser.Store.CreateResourceType(ctx, claims.TenantID, claims.UserID, req.TypeName, req.Description )
}

func (ser *Service) UpdateResourceType(ctx context.Context, req UpdateResourceTypeRequest, claims *auth.UserClaims) (int64, error) {
	return ser.Store.UpdateResourceType(ctx, claims.TenantID, claims.UserID, req.ID, req.TypeName, req.Description)
}

func (ser *Service) DeleteResourceType(ctx context.Context, id int64, claims *auth.UserClaims) (int64, error) {
	return ser.Store.DeleteResourceType(ctx, claims.TenantID, claims.UserID, id)
}



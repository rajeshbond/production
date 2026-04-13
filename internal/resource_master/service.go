package resourcemaster

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

func (ser *Service) CreateResource(ctx context.Context, req CreateResourceRequest, claims *auth.UserClaims) (int64, error) {
	tx, err := ser.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	id, err := ser.Store.CreateResource(ctx, tx, claims.TenantID, req.ResourceName, req.ResourceTypeID, claims.UserID)
	if err != nil {
		return 0, nil
	}

	return id, nil

}

func (ser *Service) GetAllResources(ctx context.Context, tenantID int64) ([]ResourceResponse, error) {
	return ser.Store.GetAllResources(ctx, tenantID)
}

func (ser *Service) GetResourceByID(ctx context.Context, typeID, tenanrID int64) (ResourceResponse, error) {
	tx, err := ser.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return ResourceResponse{}, nil
	}
	defer tx.Rollback()
	return ser.Store.GetResourceByID(ctx, tx, tenanrID, typeID)

}

func (ser *Service) UpdateResorce(ctx context.Context, req UpdateResourceRequest, claims *auth.UserClaims) error {
	tx, _ := ser.Store.db.BeginTx(ctx, nil)

	defer tx.Rollback()
	err := ser.Store.UpdateResorce(
		ctx,
		tx,
		claims.TenantID,
		req.ID,
		req.ResourceName,
		req.ResourceTypeID,
		claims.UserID,
	)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (ser *Service) DeleteResource(ctx context.Context, typeID int64, claims *auth.UserClaims) error {
	tx, _ := ser.Store.db.BeginTx(ctx, nil)

	defer tx.Rollback()

	err := ser.Store.DeleteResource(ctx, tx, claims.TenantID, typeID, claims.UserID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

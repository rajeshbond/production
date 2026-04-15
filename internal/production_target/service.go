package productiontarget

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

func (ser *Service) CreateProductionTarget(ctx context.Context, req *CreateProductionTargetRequest, claims *auth.UserClaims) (int64, error) {
	tx, err := ser.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, nil
	}

	id, err := ser.Store.CreateProductionTarget(ctx, tx, claims.TenantID, req, claims.UserID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (ser *Service) UpdateProductionTarget(ctx context.Context, req *UpdateProductionTargetRequest, claims *auth.UserClaims) error {
	tx, err := ser.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = ser.Store.UpdateProductionTarget(ctx, tx, claims.TenantID, claims.UserID, req)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *Service) GetProductionTargetByID(
	ctx context.Context,
	tenantID int64,
	id int64,
) (*ProductionTargetResponseDTO, error) {

	return s.Store.GetProductionTargetByID(ctx, tenantID, id)
}

func (s *Service) GetAllProductionTargets(
	ctx context.Context,
	tenantID int64,
) ([]ProductionTargetResponseDTO, error) {

	return s.Store.GetAllProductionTargets(ctx, tenantID)
}

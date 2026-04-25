package setupoperation

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

func (s *Service) Create(ctx context.Context, req *CreateSetupRequest, claims *auth.UserClaims) (int64, error) {

	tx, err := s.Store.db.BeginTx(ctx, nil)

	if err != nil {
		return 0, err
	}

	// defer tx.Rollback()

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	obj := OperationSetup{
		TenantID:     claims.TenantID,
		PosID:        req.PosID,
		MachineID:    req.MachineID,
		TargetQty:    req.TargetQty,
		SetupName:    req.SetupName,
		CycleTimeSec: req.CycleTimeSec,
		SetupTimeMin: req.SetupTimeMin,
		CreatedBy:    &claims.UserID,
		UpdatedBy:    &claims.UserID,
	}

	id, err := s.Store.CreateSetup(ctx, tx, &obj)
	if err != nil {
		return 0, nil
	}

	err = s.Store.UpsertResource(ctx, tx, claims.TenantID, id, req.Resources)

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return id, nil

}

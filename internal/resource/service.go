package resource

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

func (ser *Service) Create(ctx context.Context, req *CreateResourceRequest, claims *auth.UserClaims) (int64, error) {

	tx, err := ser.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	r := Resource{
		TenantID:     claims.TenantID,
		ResourceCode: req.ResourceCode,
		ResourceName: req.ResourceName,
		ResourceType: req.ResourceType,
		Description:  req.Description,
		MoldID:       req.MoldID,
		FixtureID:    req.FixtureID,
		ToolID:       req.ToolID,
		CreatedBy:    &claims.UserID,
		UpdatedBy:    &claims.UserID,
	}

	id, err := ser.Store.Create(ctx, tx, &r)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil

}

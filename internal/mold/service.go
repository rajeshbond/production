package mold

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

func (ser *Service) BulkCreate(ctx context.Context, req BulkCreateMoldRequest, claims *auth.UserClaims) (int, int, error) {

	tx, err := ser.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}

	defer tx.Rollback()

	inserted, skipped, err := ser.Store.BulkCreate(ctx, tx, claims.TenantID, claims.UserID, req)
	if err != nil {
		return 0, 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, err
	}

	return inserted, skipped, nil

}

package productionlog

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

func (ser *Service) CreateProdctionLog(ctx context.Context, req CreateProductionLogRequest, claims *auth.UserClaims) (int64, error) {
	tx, err := ser.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	defer tx.Rollback()
	id, err := ser.Store.CreateProductionLog(ctx, tx, claims.TenantID, claims.UserID, req)
	if err != nil {
		return 0, err
	}

	return id, tx.Commit()

}

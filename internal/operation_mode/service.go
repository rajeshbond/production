package operationmode

import (
	"context"
	"strings"

	"github.com/rajesh_bond/production/internal/auth"
)

type Service struct {
	Store *Store
}

func NewService(store *Store) *Service {
	return &Service{Store: store}
}

func (s *Service) Create(ctx context.Context, req CreateOperationModeRequest, claims *auth.UserClaims) error {

	tx, err := s.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	operationModeName := strings.ToLower(req.OperationModeName)

	_, err = s.Store.Create(ctx, tx, claims.TenantID, claims.UserID, operationModeName)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Service) Update(ctx context.Context, tenantID, userID int64, req UpdateOperationModeRequest) error {

	tx, err := s.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.Store.Update(ctx, tx, tenantID, userID, req); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Service) Delete(ctx context.Context, tenantID, userID, id int64) error {

	tx, err := s.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.Store.Delete(ctx, tx, tenantID, userID, id); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Service) GetIDByName(ctx context.Context, tenantID int64, name string) (int64, error) {
	return s.Store.GetIDByName(ctx, tenantID, name)
}

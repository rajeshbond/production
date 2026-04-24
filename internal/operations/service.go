package operations

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

func (ser *Service) CreateOperations(ctx context.Context, req BulkCreateOperationRequest, claims *auth.UserClaims) (BulkOperationResult, error) {

	tx, err := ser.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return BulkOperationResult{}, nil
	}

	// Remove duplicate in request

	seen := make(map[string]bool)
	var filtered []CreateOperationsRequest

	for _, o := range req.Operations {
		if !seen[o.OperationName] {
			seen[o.OperationName] = true
			filtered = append(filtered, o)
		}
	}

	result, err := ser.Store.BulkCreateOperation(ctx, tx, claims.TenantID, claims.UserID, filtered)

	if err != nil {
		return result, err
	}

	if err := tx.Commit(); err != nil {
		return result, nil
	}

	return result, nil

}

func (ser *Service) GetAllOpeationByTenant(ctx context.Context, tenantID int64) ([]OperationResponse, error) {
	return ser.Store.GetAllOpeationByTenant(ctx, tenantID)
}

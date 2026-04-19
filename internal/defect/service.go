package defect

import (
	"context"

	"github.com/rajesh_bond/production/internal/auth"
)

type Service struct {
	Store          *Store
	TenantProvider TenantProvider
}

func NewService(store *Store) *Service {
	return &Service{Store: store}
}

func (ser *Service) CreateDefects(ctx context.Context, req BulkCreateDefectRequest, claims *auth.UserClaims) (BulkDefectResult, error) {

	tx, err := ser.Store.db.BeginTx(ctx, nil)

	if err != nil {
		return BulkDefectResult{}, err
	}

	// Remove duplicate in request

	seen := make(map[string]bool)
	var filtered []CreateDefectRequest

	for _, d := range req.Defects {
		if !seen[d.DefectName] {
			seen[d.DefectName] = true
			filtered = append(filtered, d)
		}
	}

	result, err := ser.Store.BulkCreateDefect(ctx, tx, claims.TenantID, claims.UserID, filtered)

	if err != nil {
		return result, err
	}

	if err := tx.Commit(); err != nil {
		return result, nil
	}

	return result, nil

}

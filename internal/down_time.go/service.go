package downtime

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

func (ser *Service) CreateDownTime(ctx context.Context, req BulkCreateDownTimeRequest, claims *auth.UserClaims) (BulkDownTimeResult, error) {

	tx, err := ser.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return BulkDownTimeResult{}, err
	}

	// Remove duplicate in request

	seen := make(map[string]bool)
	var filtered []CreateDowntimeRequest

	for _, d := range req.DownTime {
		if !seen[d.DowntimeName] {
			seen[d.DowntimeName] = true
			filtered = append(filtered, d)
		}

	}

	result, err := ser.Store.BulkCreateDownTime(ctx, tx, claims.TenantID, claims.UserID, filtered)

	if err != nil {
		return result, err
	}

	if err := tx.Commit(); err != nil {
		return result, nil
	}

	return result, nil
}

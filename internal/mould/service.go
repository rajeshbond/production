package mould

import (
	"github.com/rajesh_bond/production/internal/auth"
	"golang.org/x/net/context"
)

type Service struct {
	Store *Store
}

func NewService(store *Store) *Service {
	return &Service{Store: store}
}

func (ser *Service) Create(ctx context.Context, req *CreateMoldRequest, claims *auth.UserClaims) (int64, error) {

	tx, err := ser.Store.db.BeginTx(ctx, nil)

	if err != nil {
		return 0, err
	}

	// defer tx.Rollback()

	// rollback only if error happens
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	m := Mold{
		TenantID:    claims.TenantID, // ✅ IMPORTANT
		Type:        req.Type,
		MoldNo:      req.MoldNo,
		Description: nil,
		Cavities:    req.Cavities,
		TargetShots: req.TargetShots,
		CreatedBy:   &claims.UserID,
		UpdatedBy:   &claims.UserID,
	}

	// optional description handling
	if req.Description != "" {
		m.Description = &req.Description
	}

	id, err := ser.Store.Create(ctx, tx, &m)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil

}

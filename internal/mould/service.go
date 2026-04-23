package mould

import (
	"context"
	"strings"

	"github.com/rajesh_bond/production/internal/auth"
)

type Service struct {
	Store            *Store
	ResourceProvider ResourceProvider
}

func NewService(store *Store, resourceProvider ResourceProvider) *Service {
	return &Service{Store: store, ResourceProvider: resourceProvider}
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
	typeString := strings.ToLower(req.Type)
	m := Mold{
		TenantID:    claims.TenantID, // ✅ IMPORTANT
		Type:        typeString,
		MoldNo:      req.MoldNo,
		MoldName:    &req.MoldName,
		Description: &req.Description,
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

	_, err = ser.ResourceProvider.CreateResource(ctx, tx, id, claims.TenantID, claims.UserID, req.MoldNo, req.MoldName, typeString, req.Description)

	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil

}

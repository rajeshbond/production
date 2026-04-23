package tools

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

func (ser *Service) Create(ctx context.Context, req *CreateToolRequest, claims *auth.UserClaims) (int64, error) {

	tx, err := ser.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	// rollback only if error happens or panic
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	t := Tool{
		TenantID:    claims.TenantID,
		Type:        req.Type,
		ToolCode:    req.ToolCode,
		ToolName:    req.ToolName,
		Description: req.Description,
		ToolType:    req.ToolType,
		Unit:        req.Unit,
		Cost:        req.Cost,
		LifeCycles:  req.LifeCycles,
		CreatedBy:   &claims.UserID,
		UpdatedBy:   &claims.UserID,
	}

	// 🔥 CALL STORE
	id, err := ser.Store.Create(ctx, tx, &t)
	if err != nil {
		return 0, err
	}

	// ✅ COMMIT TRANSACTION
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return id, nil
}

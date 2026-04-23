package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/rajesh_bond/production/internal/auth"
)

type Service struct {
	Store            *Store
	ResourceProvider ResourceProvider
}

func NewService(store *Store, resourcePointer ResourceProvider) *Service {
	return &Service{Store: store, ResourceProvider: resourcePointer}
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
	typeName := strings.ToLower(req.Type)
	t := Tool{
		TenantID:    claims.TenantID,
		Type:        typeName,
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

	// 🔥 Create Resource (IMPORTANT: check error)
	_, err = ser.ResourceProvider.CreateResource(
		ctx,
		tx,
		id,
		claims.TenantID,
		claims.UserID,
		req.ToolCode,
		req.ToolName,
		typeName,
		*req.Description,
	)
	if err != nil {
		return 0, err
	}

	// ✅ COMMIT TRANSACTION
	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return id, nil
}

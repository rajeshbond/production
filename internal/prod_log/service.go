package prodlog

import (
	"context"
	"errors"
	"time"

	"github.com/rajesh_bond/production/internal/auth"
)

type Service struct {
	Store *Store
}

func NewService(store *Store) *Service {
	return &Service{Store: store}
}

func (ser *Service) CreateProdctionLog(ctx context.Context, req *CreateRequest, claims *auth.UserClaims) (int64, error) {

	tx, err := ser.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, nil
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if req.ActualQty != req.OkQty+req.RejectedQty {
		return 0, errors.New("actual qty mismatch")
	}

	date, _ := time.Parse("2006-01-02", req.ProductionDate)

	productionLog := &ProductionLog{
		TenantID:        claims.TenantID,
		SetupID:         req.SetupID,
		MachineID:       req.MachineID,
		ProductID:       req.ProductID,
		OperationID:     req.OperationID,
		ProductionDate:  date,
		ShiftID:         req.ShiftID,
		ShiftTimingID:   req.ShiftTimingID,
		ShiftHourSlotID: req.ShiftHourSlotID,
		SlotIndex:       req.SlotIndex,
		TargetQty:       req.TargetQty,
		ActualQty:       req.ActualQty,
		OkQty:           req.OkQty,
		RejectedQty:     req.RejectedQty,
		Scrap:           req.Scrap,
		Remarks:         req.Remarks,
		CreatedBy:       &claims.UserID,
		UpdatedBy:       &claims.UserID,
	}

	logID, err := ser.Store.CreateProdcutionLog(ctx, tx, productionLog)

	if err != nil {
		return 0, err
	}

	for _, d := range req.Defects {
		if err = ser.Store.UpsertDefectLog(ctx, tx, claims.TenantID, logID, d); err != nil {
			return 0, err
		}
	}

	for _, def := range req.Downtime {
		if err = ser.Store.UpsertDownTimeLog(ctx, tx, claims.TenantID, logID, def); err != nil {
			return 0, nil
		}
	}

	err = tx.Commit()

	return logID, nil

}

func (ser *Service) SoftDelete(ctx context.Context, logID int64, claims *auth.UserClaims) error {
	return ser.Store.SoftDelete(ctx, logID, claims.TenantID, claims.UserID)

}

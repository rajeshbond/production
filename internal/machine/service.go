package machine

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

// Create Machine

func (ser *Service) CreateMachine(ctx context.Context, req CreateMachineRequest, claims *auth.UserClaims) (*MachineResponse, error) {
	tx, _ := ser.Store.db.BeginTx(ctx, nil)

	m := &Machine{
		TenantID:    claims.TenantID,
		MachineCode: req.MachineCode,
		MachineName: req.MachineName,
		Description: req.Description,
		Capacity:    req.Capacity,
		CreatedBy:   &claims.UserID,
	}

	err := ser.Store.Create(ctx, tx, m)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &MachineResponse{
		ID:          m.ID,
		MachineCode: m.MachineCode,
		MachineName: m.MachineName,
		Description: m.Description,
		Capacity:    m.Capacity,
	}, nil

}

// Update Machine ID

func (ser *Service) UpdateMachine(ctx context.Context, req *UpdateMachineRequest, claims *auth.UserClaims) error {

	tx, _ := ser.Store.db.BeginTx(ctx, nil)

	m := &Machine{
		ID:          req.ID,
		TenantID:    claims.TenantID,
		MachineCode: req.MachineCode,
		MachineName: req.MachineName,
		Description: req.Description,
		Capacity:    req.Capacity,
		UpdatedBy:   &claims.UserID,
	}

	err := ser.Store.Update(ctx, tx, m)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()

}

func (ser *Service) DeleteMachine(ctx context.Context, tenantID, userID, machineID int64) error {
	tx, _ := ser.Store.db.BeginTx(ctx, nil)

	err := ser.Store.DeleteMachine(ctx, tx, tenantID, userID, machineID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()

}

func (ser *Service) GetMachineBYID(ctx context.Context, tenantID, machineID int64) (*MachineResponse, error) {
	m, err := ser.Store.GetMachineByID(ctx, tenantID, machineID)
	if err != nil {
		return nil, err
	}

	return &MachineResponse{
		ID:          m.ID,
		MachineCode: m.MachineCode,
		MachineName: m.MachineName,
		Description: m.Description,
		Capacity:    m.Capacity,
	}, nil
}

func (ser *Service) GetAllMachineByTenant(ctx context.Context, tenantID int64) ([]MachineResponse, error) {
	list, err := ser.Store.GetAllMachineByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	var res []MachineResponse
	for _, m := range list {
		res = append(res,
			MachineResponse{
				ID:          m.ID,
				MachineCode: m.MachineCode,
				MachineName: m.MachineName,
				Description: m.Description,
				Capacity:    m.Capacity,
			})
	}

	return res, nil
}

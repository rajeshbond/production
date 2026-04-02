package tenantshifts

import (
	"context"
	"errors"
)

type Service struct {
	Store *Store
}

func NewService(store *Store) *Service {
	return &Service{Store: store}
}

func (ser *Service) CreateTenantShifts(
	ctx context.Context,
	req CreateShiftRequest,
	userID int64,
) ([]TenantShiftResponse, error) {

	if len(req.Shifts) == 0 {
		return nil, errors.New("no shifts provided")
	}

	tx, err := ser.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	seen := make(map[string]bool)
	var shiftNames []string

	// ✅ Step 1: Validate + remove duplicates in request
	for _, s := range req.Shifts {

		if s.ShiftName == "" {
			return nil, errors.New("shift name cannot be empty")
		}

		if seen[s.ShiftName] {
			continue // skip duplicate in request
		}
		seen[s.ShiftName] = true

		shiftNames = append(shiftNames, s.ShiftName)
	}

	// ✅ Step 2: Get existing shifts from DB
	existingMap, err := ser.Store.GetExistingShifts(ctx, tx, req.TenantID, shiftNames)
	if err != nil {
		return nil, err
	}

	var result []TenantShiftResponse

	// ✅ Step 3: Insert only new shifts
	for _, name := range shiftNames {

		if existingMap[name] {
			continue // skip already in DB
		}

		data, err := ser.Store.CreateTenantShifts(ctx, tx, req.TenantID, name, userID)
		if err != nil {
			return nil, err
		}

		if data == nil {
			continue // safety check
		}

		result = append(result, TenantShiftResponse{
			ID:        data.ID,
			TenantID:  data.TenantID,
			ShiftName: data.ShiftName,
			CreatedBy: data.CreatedBy,
			UpdatedBy: data.UpdatedBy,
			CreatedAt: data.CreatedAt,
			UpdatedAt: data.UpdatedAt,
		})
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

// func (ser *Service) CreateTenantShifts(ctx context.Context, req CreateShiftRequest, userID int64) ([]TenantShiftResponse, error) {

// 	// if claims.Role != "tenantadmin" {
// 	// 	return nil, fmt.Errorf("only tenant admin can create tenant shifts")
// 	// }

// 	if len(req.Shifts) == 0 {
// 		return nil, errors.New("no shifts provided")
// 	}

// 	tx, err := ser.Store.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer tx.Rollback()

// 	var result []TenantShiftResponse

// 	seen := make(map[string]bool)

// 	for _, shift := range req.Shifts {
// 		if shift.ShiftName == "" {
// 			return nil, errors.New("Shift name cannot be empty")
// 		}
// 		seen[shift.ShiftName] = true

// 		data, err := ser.Store.CreateTenantShifts(ctx, tx, req.TenantID, shift.ShiftName, userID)
// 		if err != nil {
// 			return nil, err
// 		}

// 		result = append(result, TenantShiftResponse{
// 			ID:        data.ID,
// 			TenantID:  data.TenantID,
// 			ShiftName: data.ShiftName,
// 			CreatedBy: data.CreatedBy,
// 			UpdatedBy: data.UpdatedBy,
// 			CreatedAt: data.CreatedAt,
// 			UpdatedAt: data.UpdatedAt,
// 		})
// 	}

// 	if err := tx.Commit(); err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }

// func (ser *Service) CreateTenantShift(ctx context.Context, claims *auth.UserClaims, req *CreateTenantShiftRequest) (*TenantShift, error) {
// 	// Validate struct
// 	if err := utils.Validate.Struct(req); err != nil {
// 		return nil, err
// 	}
// 	// Auth check
// 	// // 🔐 Always trust claims (override)
// 	// req.TenantID = claims.TenantID
// 	// req.CreatedBy = &claims.UserID

// 	if err := utils.Validate.Struct(req); err != nil {
// 		return nil, err
// 	}

// 	if err := auth.ValidateTenantAccesswithTenantCode(claims.Role, claims.TenantID, req.TenantID); err != nil {
// 		return nil, err
// 	}

// 	// if err := auth.ValidateTenantAccess(claims.Role, claims.EmployeeID, claims.EmployeeID); err != nil {
// 	// 	return nil, err
// 	// }

// 	return ser.store.CreateTenantShift(ctx, req)

// }

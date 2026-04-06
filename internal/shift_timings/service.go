package shifttiming

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rajesh_bond/production/internal/auth"
)

type Service struct {
	Store          *Store
	TenantProvider TenantProvider
}

func NewService(store *Store, tenantProvider TenantProvider) *Service {
	return &Service{
		Store: store,

		TenantProvider: tenantProvider,
	}
}

func (ser *Service) BulkCreateShift(ctx context.Context, req BulkCreateShiftRequest, claims *auth.UserClaims) (*BulkResult, error) {

	if claims.Role != "tenantadmin" {
		return nil, ErrOnlyTenantAllowed
	}

	tx, err := ser.Store.sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	insertCount := 0
	skipCount := 0

	for _, shift := range req {
		tenantID, err := ser.TenantProvider.GetTenantIDByCode(ctx, shift.TenantCode)
		if err != nil {
			return nil, err
		}

		tenantShiftID, err := ser.Store.CreateTenantShift(ctx, tx, tenantID, shift.ShiftName, claims.UserID)
		if err != nil {
			return nil, err
		}
		existingMap, err := ser.Store.GetExisttingTimings(ctx, tx, tenantShiftID)
		if err != nil {
			return nil, err
		}

		dayTotal := make(map[int]int)

		for _, t := range shift.Timings {

			start, end, err := convertToMinutes(t.ShiftStart, t.ShiftEnd)
			if err != nil {
				return nil, err
			}

			duration := end - start

			if duration <= 0 {
				duration += 1440
			}

			if dayTotal[t.Weekday]+duration > 1440 {
				return nil, fmt.Errorf("exeeds 24 hrs weekday %d ", t.Weekday)
			}

			dayTotal[t.Weekday] += duration

			key := t.ShiftStart + "-" + t.ShiftEnd + strconv.Itoa(t.Weekday)

			if existingMap[key] {
				skipCount++
				continue
			}
			rows, err := ser.Store.InsertShifttiming(ctx, tx, t, tenantShiftID, claims.UserID)
			if err != nil {
				return nil, err
			}

			if rows == 0 {
				skipCount++
				continue
			}

			insertCount++

		}

	}

	if insertCount == 0 {
		return nil, fmt.Errorf("all shifts already exist")
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &BulkResult{
		Inserted: insertCount,
		Skipped:  skipCount,
	}, nil

}

// func (ser *Service) CreateShiftTiming(ctx context.Context, req CreateShiftTimingRequest, claims *auth.UserClaims) ([]ShiftTimingResponse, error) {

// 	// Role validation
// 	if claims.Role != "tenantadmin" {
// 		return nil, ErrOnlyTenantAllowed
// 	}
// 	tx, err := ser.Store.sqlDB.BeginTx(ctx, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer tx.Rollback()

// 	fmt.Println("Claims Tenant Id ", claims.TenantID)

// 	// Tenant validation using tenant_shift
// 	tenantCode, err := ser.TenantProvider.GetTenantCodeByID(ctx, claims.TenantID)

// 	fmt.Println("Tenant Code", tenantCode)

// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := ser.Store.ValidateTenantShiftAccess(ctx, tx, req.TenantShiftID, tenantCode); err != nil {
// 		return nil, err
// 	}

// 	// Fetch existing timings

// 	existingMap, err := ser.Store.GetExisttingTimings(ctx, tx, req.TenantShiftID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	dayTotal := make(map[int]int)
// 	var result []ShiftTimingResponse

// 	for _, t := range req.Timings {
// 		if t.Weekday < 0 || t.Weekday > 6 {
// 			return nil, fmt.Errorf("invalid weekday")
// 		}

// 		startMin, endMin, err := convertToMinutes(t.ShiftStart, t.ShiftEnd)

// 		if err != nil {
// 			return nil, err
// 		}

// 		duration := endMin - startMin

// 		if duration <= 0 {
// 			duration += 1440 // overnight
// 		}

// 		// 24 - hours validation

// 		if dayTotal[t.Weekday]+duration > 1440 {
// 			return nil, fmt.Errorf("total shoft exceeeds 24 hrs from weekday %d", t.Weekday)
// 		}

// 		dayTotal[t.Weekday] += duration

// 		key := t.ShiftStart + "-" + t.ShiftEnd + "-" + strconv.Itoa(t.Weekday)

// 		// Skip duplicate in DB

// 		if existingMap[key] {
// 			continue
// 		}

// 		data, err := ser.Store.InsertShiftTimingTx(ctx, tx, t, req.TenantShiftID, claims.UserID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if data == nil {
// 			continue
// 		}

// 		result = append(result, ShiftTimingResponse{
// 			ID:            data.ID,
// 			TenantShiftID: data.TenantShiftID,
// 			ShiftStart:    data.ShiftStart,
// 			ShiftEnd:      data.ShiftEnd,
// 			Weekday:       data.Weekday,
// 		})
// 	}
// 	if err := tx.Commit(); err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

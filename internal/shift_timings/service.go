package shifttiming

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rajesh_bond/production/internal/auth"
)

type Service struct {
	Store        *Store
	RoleProvider RoleProvider
}

func NewService(store *Store, roleProvider RoleProvider) *Service {
	return &Service{
		Store:        store,
		RoleProvider: roleProvider,
	}
}

func (ser *Service) CreateShiftTiming(ctx context.Context, req CreateShiftTimingRequest, claims *auth.UserClaims) ([]ShiftTimingResponse, error) {

	// Role validation
	if claims.Role != "tenantadmin" {
		return nil, ErrOnlyTenantAllowed
	}
	tx, err := ser.Store.sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	// Tenant validation using tenant_shift
	tenantCode, err := ser.RoleProvider.GetRoleNameByID(ctx, claims.TenantID)

	if err != nil {
		return nil, err
	}

	if err := ser.Store.ValidateTenantShiftAccess(ctx, tx, req.TenantShiftID, tenantCode); err != nil {
		return nil, err
	}

	// Fetch existing timings

	existingMap, err := ser.Store.GetExisttingTimings(ctx, tx, req.TenantShiftID)
	if err != nil {
		return nil, err
	}

	dayTotal := make(map[int]int)
	var result []ShiftTimingResponse

	for _, t := range req.Timings {
		if t.Weekday < 0 || t.Weekday > 6 {
			return nil, fmt.Errorf("invalid weekday")
		}

		startMin, endMin, err := convertToMinutes(t.ShiftStart, t.ShiftEnd)

		if err != nil {
			return nil, err
		}

		duration := endMin - startMin

		if duration <= 0 {
			duration += 1440 // overnight
		}

		// 24 - hours validation

		if dayTotal[t.Weekday]+duration > 1440 {
			return nil, fmt.Errorf("total shoft exceeeds 24 hrs from weekday %d", t.Weekday)
		}

		dayTotal[t.Weekday] += duration

		key := t.ShiftStart + "-" + t.ShiftEnd + "-" + strconv.Itoa(t.Weekday)

		// Skip duplicate in DB

		if existingMap[key] {
			continue
		}

		data, err := ser.Store.InsertShiftTimingTx(ctx, tx, t, req.TenantShiftID, claims.UserID)
		if err != nil {
			return nil, err
		}
		if data == nil {
			continue
		}

		result = append(result, ShiftTimingResponse{
			ID:            data.ID,
			TenantShiftID: data.TenantShiftID,
			ShiftStart:    data.ShiftStart,
			ShiftEnd:      data.ShiftEnd,
			Weekday:       data.Weekday,
		})
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

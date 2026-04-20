package shift

import (
	"context"
	"fmt"
)

type Service struct {
	Store *Store
}

func NewService(store *Store) *Service {
	return &Service{Store: store}
}

// ==========================================
// SINGLE CREATE (optional)
// ==========================================
func (s *Service) Create(ctx context.Context, userID int64, req CreateShiftTimingRequest) error {

	tx, err := s.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := s.Store.Create(ctx, tx, userID, req); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// ==========================================
// BULK CREATE (MAIN API)
// ==========================================
func (s *Service) CreateBulk(
	ctx context.Context,
	userID int64,
	req BulkCreateShiftRequest,
) error {

	if len(req) == 0 {
		return fmt.Errorf("empty request")
	}

	tx, err := s.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// in-memory cache
	tenantMap := make(map[string]int64)
	dayMap := make(map[string]map[int][][2]int)

	for _, shift := range req {

		// ❌ validate empty timings
		if len(shift.Timings) == 0 {
			tx.Rollback()
			return fmt.Errorf("no timings for shift %s", shift.ShiftName)
		}

		// 🔹 get tenant_id (cached)
		tenantID, ok := tenantMap[shift.TenantCode]
		if !ok {
			id, err := s.Store.GetTenantIDByCode(ctx, tx, shift.TenantCode)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("tenant not found: %s", shift.TenantCode)
			}
			tenantMap[shift.TenantCode] = id
			tenantID = id
		}

		// init weekday map
		if _, ok := dayMap[shift.TenantCode]; !ok {
			dayMap[shift.TenantCode] = make(map[int][][2]int)
		}

		// 🔹 upsert shift
		shiftID, err := s.Store.UpsertTenantShift(
			ctx, tx, tenantID, shift.ShiftName, userID,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to upsert shift %s: %w", shift.ShiftName, err)
		}

		// 🔹 timings loop
		for _, t := range shift.Timings {

			start, err := toMinutes(t.ShiftStart)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("invalid start time %s", t.ShiftStart)
			}

			end, err := toMinutes(t.ShiftEnd)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("invalid end time %s", t.ShiftEnd)
			}

			// overnight handling
			if end <= start {
				end += 1440
			}

			existing := dayMap[shift.TenantCode][t.Weekday]

			// ❌ overlap check
			for _, ex := range existing {
				if start < ex[1] && end > ex[0] {
					tx.Rollback()
					return fmt.Errorf(
						"overlap detected for tenant %s weekday %d",
						shift.TenantCode, t.Weekday,
					)
				}
			}

			// add interval
			dayMap[shift.TenantCode][t.Weekday] =
				append(existing, [2]int{start, end})

			// ❌ 24-hour validation
			total := 0
			for _, ex := range dayMap[shift.TenantCode][t.Weekday] {
				total += ex[1] - ex[0]
			}

			if total > 1440 {
				tx.Rollback()
				return fmt.Errorf(
					"total shift exceeds 24h for tenant %s weekday %d",
					shift.TenantCode, t.Weekday,
				)
			}

			// insert timing
			if err := s.Store.InsertShiftTiming(ctx, tx, shiftID, t, userID); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to insert shift timing: %w", err)
			}
		}
	}

	return tx.Commit()
}

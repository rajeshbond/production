package operationdowntimemap

import (
	"context"
	"strings"

	"github.com/rajesh_bond/production/internal/auth"
)

type Service struct {
	store             *Store
	DowntimeProvide   DowntimeProvider
	OperationProvider OperationProvider
}

func NewService(store *Store, downtimeProvide DowntimeProvider, operationProvider OperationProvider) *Service {
	return &Service{
		store:             store,
		DowntimeProvide:   downtimeProvide,
		OperationProvider: operationProvider,
	}
}

func (ser *Service) CreateOperationWithDowntime(ctx context.Context, req OperationDowntimeCreateRequest, claims *auth.UserClaims) (OperationDowntimeCreateResponse, error) {
	// Start transaction
	tx, err := ser.store.db.BeginTx(ctx, nil)
	if err != nil {
		return OperationDowntimeCreateResponse{}, err
	}
	defer tx.Rollback()

	// Normalize operation name
	opname := strings.TrimSpace(req.OperationName)

	// 1️⃣ Get or create operation
	operationID, err := ser.OperationProvider.GetOperationIDByName(ctx, tx, claims.TenantID, opname)

	if err != nil {
		return OperationDowntimeCreateResponse{}, err

	}
	if operationID == 0 {
		operationID, err = ser.OperationProvider.CreateOperation(ctx, tx, claims.TenantID, claims.UserID, opname)

	}
	if err != nil {
		return OperationDowntimeCreateResponse{}, err
	}

	// 2. Deduplicate downtime name

	uniqueMap := make(map[string]string) // normalized -> original
	cleanedDefects := []string{}
	duplicateInputs := []string{}

	for _, d := range req.DownTimeName {
		trimmed := strings.ToLower(strings.TrimSpace(d))
		normalized := strings.ToLower(d)

		if _, exists := uniqueMap[normalized]; exists {
			duplicateInputs = append(duplicateInputs, trimmed)
			continue
		}

		uniqueMap[normalized] = trimmed
		cleanedDefects = append(cleanedDefects, trimmed)

	}

	inserted := []string{}
	skipped := []string{}

	// Process each defect

	for _, downtimeName := range cleanedDefects {

		// Get defect ID
		downTimeID, err := ser.DowntimeProvide.GetDowntimeIDByName(ctx, tx, claims.TenantID, downtimeName)

		if err != nil {
			skipped = append(skipped, downtimeName)
			continue
		}
		// Create defect if it doesn't exist
		if downTimeID == 0 {

			downTimeID, err = ser.DowntimeProvide.CreateDowntime(ctx, tx, claims.TenantID, claims.UserID, downtimeName)
			if err != nil || downTimeID == 0 {
				skipped = append(skipped, downtimeName)
				continue
			}
		}

		existID, err := ser.store.GetOperationDowntimeMap(ctx, tx, claims.TenantID, operationID, downTimeID)

		if err != nil {
			skipped = append(skipped, downtimeName)
			continue
		}
		if existID > 0 {
			skipped = append(skipped, downtimeName)
			continue
		}

		// Map operation <---> Downtime

		if existID == 0 {
			id, err := ser.store.InsertOperationDowntimeMap(ctx, tx, claims.TenantID, operationID, downTimeID, claims.UserID)
			if err != nil {
				skipped = append(skipped, downtimeName)
				continue
			}
			if id > 0 {
				inserted = append(inserted, downtimeName)
				continue
			}
			if id == 0 {
				skipped = append(skipped, downtimeName)
				continue
			}

			// Only append to inserted if downtime and mapping succeeded

		}
	}

	// Add duplicate to skipped

	skipped = append(skipped, duplicateInputs...)

	// Commit transation

	if err := tx.Commit(); err != nil {
		return OperationDowntimeCreateResponse{}, nil
	}

	return OperationDowntimeCreateResponse{
		OperationID: operationID,
		Inserted:    inserted,
		Skipped:     skipped,
	}, nil

}

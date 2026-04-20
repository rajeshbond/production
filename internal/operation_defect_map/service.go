package operationdefectmap

import (
	"context"
	"strings"

	"github.com/rajesh_bond/production/internal/auth"
)

type Service struct {
	store             *Store
	defectProvide     DefectProvider
	OperationProvider OperationProvider
}

func NewService(store *Store, defectProvide DefectProvider, operationProvider OperationProvider) *Service {
	return &Service{
		store:             store,
		defectProvide:     defectProvide,
		OperationProvider: operationProvider,
	}
}

func (s *Service) CreateOperationWithDefect(
	ctx context.Context,
	req OperationDefectCreateRequest,
	claims auth.UserClaims,
) (OperationDefectCreateResponse, error) {

	// Start transaction
	tx, err := s.store.db.BeginTx(ctx, nil)
	if err != nil {
		return OperationDefectCreateResponse{}, err
	}
	defer tx.Rollback()

	// Normalize operation name
	opName := strings.TrimSpace(req.OperationName)

	// 1️⃣ Get or create operation
	operationID, err := s.OperationProvider.GetOperationIDByName(ctx, tx, claims.TenantID, opName)
	if err != nil {
		return OperationDefectCreateResponse{}, err
	}

	if operationID == 0 {
		operationID, err = s.OperationProvider.CreateOperation(ctx, tx, claims.TenantID, claims.UserID, opName)
		if err != nil {
			return OperationDefectCreateResponse{}, err
		}
	}

	// 2️⃣ Deduplicate defect names
	uniqueMap := make(map[string]string) // normalized -> original
	cleanedDefects := []string{}
	duplicateInputs := []string{}

	for _, d := range req.DefectNames {
		trimmed := strings.TrimSpace(d)
		normalized := strings.ToLower(trimmed)

		if _, exists := uniqueMap[normalized]; exists {
			duplicateInputs = append(duplicateInputs, trimmed)
			continue
		}

		uniqueMap[normalized] = trimmed
		cleanedDefects = append(cleanedDefects, trimmed)
	}

	inserted := []string{}
	skipped := []string{}

	// 3️⃣ Process each defect
	for _, defectName := range cleanedDefects {
		// Get defect ID
		defectID, err := s.defectProvide.GetDefectIDByName(ctx, tx, claims.TenantID, defectName)
		if err != nil {
			skipped = append(skipped, defectName)
			continue
		}

		// Create defect if it doesn't exist
		if defectID == 0 {
			defectID, err = s.defectProvide.CreateDefect(ctx, tx, claims.TenantID, claims.UserID, defectName)
			if err != nil || defectID == 0 {
				skipped = append(skipped, defectName)
				continue
			}
		}

		exitID, err := s.store.GetOperationDefectMap(ctx, tx, claims.TenantID, operationID, defectID)

		if err != nil {
			skipped = append(skipped, defectName)
			continue
		}

		if exitID > 0 {
			skipped = append(skipped, defectName)
			continue
		}

		// Map operation ↔ defect
		if exitID == 0 {
			id, err := s.store.InsertOperationDefectMap(ctx, tx, claims.TenantID, operationID, defectID)
			if err != nil {
				skipped = append(skipped, defectName)
				continue
			}
			if id > 0 {
				inserted = append(inserted, defectName)
			}
			if id == 0 {
				skipped = append(skipped, defectName)
				continue
			}
		}

		// Only append to inserted if defect exists and mapping succeeded

	}

	// 4️⃣ Add duplicates to skipped
	skipped = append(skipped, duplicateInputs...)

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return OperationDefectCreateResponse{}, err
	}

	return OperationDefectCreateResponse{
		OperationID: operationID,
		Inserted:    inserted,
		Skipped:     skipped,
	}, nil
}

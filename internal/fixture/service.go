package fixture

import (
	"context"
	"encoding/json"

	"github.com/rajesh_bond/production/internal/auth"
)

type Service struct {
	Store *Store
}

func NewService(store *Store) *Service {
	return &Service{Store: store}
}

func (s *Service) Create(ctx context.Context, req *CreateFixtureRequest, claims *auth.UserClaims) (int64, error) {
	tx, err := s.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// var notesBytes []byte

	// // ✅ Handle JSONB
	// if req.SpecialNotes != nil {
	// 	notesBytes, err = json.Marshal(req.SpecialNotes)
	// 	if err != nil {
	// 		return 0, err
	// 	}
	// }

	var notesByte []byte

	// Handle JSONB
	if req.SpecialNotes != nil {
		notesByte, err = json.Marshal(req.SpecialNotes)
	}

	// Hansle JSONB
	f := Fixture{
		TenantID:     claims.TenantID,
		Type:         req.Type,
		FixtureNo:    req.FixtureNo,
		FixtureName:  req.FixtureName,
		Description:  req.Description,
		Cavities:     req.Cavities,
		LifeShots:    req.LifeShots,
		FixtureType:  req.FixtureType,
		Material:     req.Material,
		SpecialNotes: notesByte,
		CreatedBy:    &claims.UserID,
		UpdatedBy:    &claims.UserID,
	}

	id, err := s.Store.Create(ctx, tx, &f)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, nil
	}

	return id, nil
}

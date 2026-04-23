package fixture

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/rajesh_bond/production/internal/auth"
)

type Service struct {
	Store            *Store
	ResourceProvider ResourceProvider
}

func NewService(store *Store, resourceProvide ResourceProvider) *Service {
	return &Service{Store: store, ResourceProvider: resourceProvide}
}

func (s *Service) Create(ctx context.Context, req CreateFixtureRequest, claims *auth.UserClaims) (int64, error) {
	tx, err := s.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var notesByte []byte

	// Handle JSONB
	if req.SpecialNotes != nil {
		notesByte, err = json.Marshal(req.SpecialNotes)
	}
	typeString := strings.ToLower(req.Type)
	fixtureName := ""
	if req.FixtureName != nil {
		fixtureName = *req.FixtureName
	}

	description := ""
	if req.Description != nil {
		description = *req.Description
	}
	// Hansle JSONB
	f := Fixture{
		TenantID:     claims.TenantID,
		Type:         typeString,
		FixtureNo:    req.FixtureNo,
		FixtureName:  &fixtureName,
		Description:  &description,
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

	_, err = s.ResourceProvider.CreateResource(ctx, tx, id, claims.TenantID, claims.UserID, req.FixtureNo, fixtureName, typeString, description)
	if err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, nil
	}

	return id, nil
}

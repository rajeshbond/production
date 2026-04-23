package fixture

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Create(ctx context.Context, tx *sql.Tx, f *Fixture) (int64, error) {

	query := `
		INSERT INTO fixture(tenant_id,type,fixture_no,fixture_name,description,cavities,life_shots,fixture_type,material,special_notes,created_by, updated_by)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		RETURNING id
	`
	var id int64
	typeName := strings.ToLower(f.Type)
	err := tx.QueryRowContext(ctx, query,
		f.TenantID,
		typeName,
		f.FixtureNo,
		f.FixtureName,
		f.Description,
		f.Cavities,
		f.LifeShots,
		f.FixtureType,
		f.Material,
		f.SpecialNotes,
		f.CreatedBy,
		f.UpdatedBy,
	).Scan(&id)

	if err != nil {
		fmt.Println(err)
		return 0, nil
	}

	return id, nil

}

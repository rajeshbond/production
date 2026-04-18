package mold

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) BulkCreate(ctx context.Context, tx *sql.Tx, tenantID, userID int64, req BulkCreateMoldRequest) (inserted, skipped int, err error) {
	if len(req.Molds) == 0 {
		return 0, 0, nil
	}

	// -------------------------------------------------
	// 1. Collect mold_nos from request
	// -------------------------------------------------

	moldNos := make([]string, 0, len(req.Molds))
	for _, m := range req.Molds {
		moldNos = append(moldNos, m.MoldNo)
	}

	// -------------------------------------------------
	// 2. Fetch already existing (NOT deleted)
	// -------------------------------------------------

	checkQuery := `
		SELECT mold_no 
		FROM mold
		WHERE tenant_id = $1
			AND is_deleted = FALSE
			AND mold_no = ANY($2)
	`

	rows, err := tx.QueryContext(ctx, checkQuery, tenantID, pq.Array(moldNos))
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	existingMap := make(map[string]struct{})
	for rows.Next() {
		var moldNo string
		if err := rows.Scan(&moldNo); err != nil {
			return 0, 0, err
		}
		existingMap[moldNo] = struct{}{}
	}
	// -------------------------------------------------
	// 3. Insert only non-existing
	// -------------------------------------------------

	insertQuery := `
		INSERT INTO mold (
			tenant_id, mold_no, description, cavities, target_shots, special_notes, created_by, updated_by
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$7)
	`

	for _, m := range req.Molds {

		// Skip if already exists (in DB)

		if _, found := existingMap[m.MoldNo]; found {
			skipped++
			continue
		}

		_, err := tx.ExecContext(ctx, insertQuery,
			tenantID,
			m.MoldNo,
			m.Description,
			m.TargetShots,
			m.SpecialNotes,
			userID,
		)

		if err != nil {
			// Handle race condition diplicate
			if isUniqueViolation(err) {
				skipped++
				continue
			}
			return inserted, skipped, err
		}

		inserted++

	}

	return inserted, skipped, err

}

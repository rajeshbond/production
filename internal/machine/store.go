package machine

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Create(ctx context.Context, tx *sql.Tx, m *Machine) (int64, error) {

	query := `
	INSERT INTO machine(
		tenant_id,
		machine_code,
		machine_name,
		description,
		capacity,
		special_notes,
		created_by,
		updated_by
	)
	VALUES($1,$2,$3,$4,$5,$6,$7,$7)
	RETURNING id;
	`

	var id int64

	err := tx.QueryRowContext(ctx, query,
		m.TenantID,
		m.MachineCode,
		m.MachineName,
		m.Description,
		m.Capacity,
		m.SpecialNotes,
		m.CreatedBy,
	).Scan(&id)

	return id, err
}

// 1.Create
// func (s *Store) Create(ctx context.Context, tx *sql.Tx, m *Machine) error {
// 	query := `
// 	INSERT INTO machine
// 	(tenant_id, machine_code, machine_name, description, capacity, created_by, updated_by)
// 	VALUES ($1,$2,$3,$4,$5,$6,$6)

// 	ON CONFLICT (tenant_id, machine_code)
// 	WHERE is_deleted = FALSE
// 	DO UPDATE SET
// 	    machine_name = EXCLUDED.machine_name,
// 	    description  = EXCLUDED.description,
// 	    capacity     = EXCLUDED.capacity,
// 	    updated_by   = EXCLUDED.updated_by,
// 	    updated_at   = NOW()

// 	RETURNING id, created_at, updated_at
// 	`
// 	return tx.QueryRowContext(ctx, query,
// 		m.TenantID,
// 		m.MachineCode,
// 		m.MachineName,
// 		m.Description,
// 		m.Capacity,
// 		m.CreatedBy,
// 	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)

// }

// 2.Update

func (s *Store) Update(ctx context.Context, tx *sql.Tx, m *Machine) error {
	query := `
	UPDATE machine
	SET machine_code = $1,
		machine_name = $2,
		description = $3,
		capacity = $4,
		updated_by = $5,
		WHERE id = $6 AND tenant_id = $7 AND is_deleted = FALSE
`
	_, err := tx.ExecContext(ctx, query,
		m.MachineCode,
		m.MachineName,
		m.Description,
		m.Capacity,
		m.UpdatedBy,
		m.ID,
		m.TenantID,
	)

	return err

}

// 3.Get the Machine by ID

func (s *Store) GetMachineByID(ctx context.Context, tenant_id, id int64) (*Machine, error) {
	query := `
		SELECT id,tenant_id,machine_code,machine_name,description,capacity,created_at,updated_at
		FROM machine
		WHERE id = $1 AND tenant_id = $2 AND is_deleted = FALSE
	`
	var m Machine

	err := s.db.QueryRowContext(ctx, query,
		id, tenant_id).Scan(
		&m.ID, &m.TenantID, &m.MachineCode, &m.MachineName, &m.Capacity, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No machine found with Provided ID")
		}
		return nil, err
	}

	return &m, nil
}

// Fetch Machine by tenant ID

func (s *Store) GetAllMachineByTenantID(ctx context.Context, tenantID int64) ([]Machine, error) {
	query := `
		SELECT id, tenant_id, machine_code, machine_name, description, capacity,
	       created_at, updated_at
	FROM machine
	WHERE tenant_id = $1 AND is_deleted = FALSE
	ORDER BY id DESC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var list []Machine

	for rows.Next() {
		var m Machine
		err := rows.Scan(
			&m.ID,
			&m.TenantID,
			&m.MachineCode,
			&m.MachineName,
			&m.Description,
			&m.Capacity,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}

	return list, nil

}

// Delete machine

func (s *Store) DeleteMachine(ctx context.Context, tx *sql.Tx, tenantID, userID, machineID int64) error {
	query := `
		UPDATE machine 
		SET id_deleted = TRUE
				deleted_by = $1,
	    deleted_at = NOW()
		WHERE id = $2 AND tenant_id = $3 AND is_deleted = FALSE
	`
	_, err := tx.ExecContext(ctx, query, userID, machineID, tenantID)

	return err
}

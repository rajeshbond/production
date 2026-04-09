package operationdefectmap

import (
	"context"
	"database/sql"
)

type DefectProvider interface {
	GetDefectIDByName(ctx context.Context, tx *sql.Tx, tenantID int64, defectName string) (int64, error)
	CreateDefect(ctx context.Context, tx *sql.Tx, tenantID int64, userID int64, defectName string) (int64, error)
}

package operationdowntimemap

import (
	"context"
	"database/sql"
)

type OperationProvider interface {
	GetOperationIDByName(ctx context.Context, tx *sql.Tx, tenantID int64, operationName string) (int64, error)
	CreateOperation(ctx context.Context, tx *sql.Tx, tenantID int64, userID int64, operationName string) (int64, error)
}

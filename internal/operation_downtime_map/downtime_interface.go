package operationdowntimemap

import (
	"context"
	"database/sql"
)

type DowntimeProvider interface {
	CreateDowntime(ctx context.Context, tx *sql.Tx, tenantID int64, userID int64, downtimeName string) (int64, error)

	GetDowntimeIDByName(ctx context.Context, tx *sql.Tx, tenantID int64, downtimeName string) (int64, error)
}

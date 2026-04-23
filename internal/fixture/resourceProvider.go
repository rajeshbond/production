package fixture

import (
	"context"
	"database/sql"
)

type ResourceProvider interface {
	// Resource Provide store fucntion
	CreateResource(ctx context.Context, tx *sql.Tx, rid, tenantID, userID int64, resourceCode, resourceName, resourceType, description string) (int64, error)
}

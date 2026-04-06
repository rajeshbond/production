package defect

import "golang.org/x/net/context"

type TenantProvider interface {
	GetTenantCodeByID(ctx context.Context, tenantID int64) (string, error)
	GetTenantIDByCode(ctx context.Context, tenantName string) (int64, error)
}

package users

import "context"

type TenantProvider interface {
	GetTenantIDByCode(ctx context.Context, tenantCode string) (int64, error)
}

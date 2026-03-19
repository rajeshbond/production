package tenant

import "errors"

var (
	ErrTenantCodeExists = errors.New("tenant Code already exists")

	ErrTenantNameRequired = errors.New("tenant name is required")

	ErrTenantCodeRequired = errors.New("tenant code is required")

	ErrTenantCodeNotFount = errors.New("tenant not found")

	ErrTenantAddressRequired = errors.New("tenant address is required")

	ErrTenantAlreadyVerified = errors.New("tenant already verified")
	ErrTenantNotDeleted      = errors.New("tenant not deleted")
	ErrSameTenantCode        = errors.New("Same Tenant code exists")
)

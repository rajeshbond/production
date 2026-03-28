package tenantshifts

import "errors"

var (
	ErrTenantShiftAlreadyExists = errors.New("tenant shift already exists")
	ErrInvalidRequest           = errors.New("invalid request")
)

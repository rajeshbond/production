package shifttiming

import "errors"

var (
	ErrTenantShiftAlreadyExists = errors.New("tenant shift already exists")
	ErrInvalidRequest           = errors.New("invalid request")
	ErrUnAuthorized             = errors.New("UnAuthorised ")
	ErrOnlyTenantAllowed        = errors.New("Only Tenant Admin Allowed to create")
)

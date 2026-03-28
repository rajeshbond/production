package auth

import "errors"

var (
	ErrUnauthorized   = errors.New("you are not authorized")
	ErrTenantMismatch = errors.New("tenant mismatch: not allowed for other tenant")
	ErrForbidden      = errors.New("forbidden")
)

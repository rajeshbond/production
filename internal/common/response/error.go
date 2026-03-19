package response

import "errors"

var (
	ErrTeantNameRequired   = errors.New("tenanr name Required")
	ErrTenantCondeRequired = errors.New("tenant code Required")
	ErrInvalidRequest      = errors.New("invalid request")
)

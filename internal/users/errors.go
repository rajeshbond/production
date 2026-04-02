package users

import "errors"

var (
	// 🔹 Validation Errors
	ErrEmployeeIDRequired            = errors.New("employee id is required")
	ErrUserNameRequired              = errors.New("user name is required")
	ErrPasswordRequired              = errors.New("password is required")
	ErrTenantIDRequired              = errors.New("tenant id is required")
	ErrRoleIDRequired                = errors.New("role id is required")
	ErrUserDeleted                   = errors.New("User Already Deleted")
	ErrUserAlreadyExistForThisTenant = errors.New("User Already Exist for this Tenant")

	// 🔹 Business Logic Errors
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")

	// 🔹 Auth / Security
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized access")

	// 🔹 Generic
	ErrInvalidRequest       = errors.New("invalid request")
	ErrInternalServer       = errors.New("internal server error")
	ErrEmployeeIDReqyured   = errors.New("employee_id is required")
	ErrAlreadyTenantPresent = errors.New("tenant not Present")

	// Tenant

	ErrTenantDeleted         = errors.New("Tenant Already Deleted")
	ErrTenantInActive        = errors.New("Tenant Not active Please cal Admin")
	ErrTenantVerified        = errors.New("Tenant Not Verified, Please contact Admin")
	ErrRoleNotFound          = errors.New("Role not found")
	ErrOnlyTenantAdminCreate = errors.New("Only Can Create Tenant Admin")
	ErrTenantIDMismatched    = errors.New("Tenant ID Mismatched with Employee ID")
)

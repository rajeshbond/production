package auth

import (
	"errors"
	"fmt"
	"strings"
)

func IsSuper(role string) bool {
	return role == "admin" || role == "superadmin"
}

func ValidateTenantAccess(role, claimsEmpID, reqEmpID string) error {

	role = strings.ToLower(strings.TrimSpace(role))

	// ✅ Superadmin & Admin → full access
	if role == "superadmin" || role == "admin" {
		return nil
	}

	// ✅ Tenant Admin → restricted
	if role == "tenantadmin" {

		claimsTcode, err := Tcode(claimsEmpID)
		if err != nil {
			return errors.New("invalid claims employee id")
		}

		reqTcode, err := Tcode(reqEmpID)
		if err != nil {
			return errors.New("invalid request employee id")
		}

		if claimsTcode != reqTcode {
			return errors.New("tenant mismatch: not allowed for other Tenant")
		}

		return nil
	}

	// ❌ Other roles
	return errors.New("insufficient permissions")
}

func Tcode(employee_id string) (string, error) {
	parts := strings.SplitN(employee_id, "@", 2)

	if len(parts) < 2 || parts[1] == "" {
		return "", errors.New("invalid employee id format")
	}

	return strings.ToLower(strings.TrimSpace(parts[1])), nil
}

// Function Validate Tenant ID with Tenant code

func ValidateTenantAccesswithTenantCode(role string, claimsTenantID, reqTenantID int64) error {

	switch role {

	case RoleSuperAdmin, RoleAdmin:
		// Full Access
		return nil

	case RoleTenantAdmin:
		// Restricted to own tenant
		if claimsTenantID != reqTenantID {
			fmt.Println("Rajesh failed ")
			return ErrTenantMismatch
		}

		return nil // ✅ IMPORTANT FIX

	default:
		return ErrUnauthorized
	}
}

func TenantRoleCheck(role string) error {
	switch role {
	case "superadmin", "admin", "tenantadmin":
		return fmt.Errorf("not allowed to create admin role")
	default:
		return nil
	}
}

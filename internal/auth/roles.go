package auth

func IsSuper(role string) bool {
	return role == "admin" || role == "superadmin"
}

// func GetSuperSuperAdmin(ctx context.Context, roles ...Role) (*UserClaims, error) {

// 	claims, ok := GetUserClaimsFromContext(ctx)

// 	if !ok {
// 		return nil, ErrUnauthorized
// 	}

// 	for _, r := range roles {
// 		if strings.EqualFold(string(claims.Role), string(r)) {
// 			return claims, nil
// 		}
// 	}
// 	return nil, ErrForbidden

// }

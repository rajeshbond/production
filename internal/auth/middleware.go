package auth

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/rajesh_bond/production/internal/contextkey"
)

// Verifier middleware
func Verifier(tokenAuth *jwtauth.JWTAuth) func(next http.Handler) http.Handler {
	return jwtauth.Verifier(tokenAuth)
}

// Authenticator middleware
func Authenticator(tokenAuth *jwtauth.JWTAuth) func(next http.Handler) http.Handler {
	return jwtauth.Authenticator(tokenAuth)
}

// UserContextInjector extracts JWT claims and stores them in context
func UserContextInjector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			log.Printf("Error reading JWT claims: %v", err)
			http.Error(w, "Authentication error: JWT claims missing.", http.StatusUnauthorized)
			return
		}

		userClaims := UserClaims{}

		if v, ok := claims["employee_id"].(string); ok {
			userClaims.EmployeeID = v
		}

		if v, ok := claims["username"].(string); ok {
			userClaims.Username = v
		}

		if v, ok := claims["user_id"].(float64); ok {
			userClaims.UserID = int64(v)
		}

		if v, ok := claims["tenant_id"].(float64); ok {
			userClaims.TenantID = int64(v)
		}

		if v, ok := claims["role_id"].(float64); ok {
			userClaims.RoleID = int64(v)
		}

		if v, ok := claims["iat"].(float64); ok {
			userClaims.Iat = int64(v)
		}

		ctx := context.WithValue(r.Context(), contextkey.KeyUser, userClaims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserClaimsFromContext returns full JWT claims struct
func GetUserClaimsFromContext(ctx context.Context) (*UserClaims, bool) {

	claims, ok := ctx.Value(contextkey.KeyUser).(UserClaims)
	if !ok {
		return nil, false
	}

	return &claims, true
}

// UserContextInjector reads claims and injects userID
// func UserContextInjector(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		_, claims, err := jwtauth.FromContext(r.Context())
// 		if err != nil {
// 			log.Printf("Error: Failed to read claims from context: %v", err)
// 			http.Error(w, "Authentication error: Context claims missing.", http.StatusUnauthorized)
// 			return
// 		}

// 		var userID string

// 		switch v := claims["user_id"].(type) {

// 		case float64:
// 			userID = strconv.FormatInt(int64(v), 10)

// 		case string:
// 			userID = v

// 		case json.Number:
// 			userID = v.String()

// 		default:
// 			log.Printf("SECURITY ALERT: Invalid user_id type %T value %v", v, v)
// 			http.Error(w, "Authentication error: User ID claims missing or invalid.", http.StatusUnauthorized)
// 			return
// 		}

// 		if userID == "" {
// 			http.Error(w, "Authentication error: Empty user ID.", http.StatusUnauthorized)
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), contextkey.KeyUser, userID)

// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

// GetUserIDFromContext extracts userID from context
// func GetUserIDFromContext(ctx context.Context) string {

// 	userID, ok := ctx.Value(contextkey.KeyUser).(string)
// 	if !ok {
// 		return ""
// 	}

// 	return userID
// }

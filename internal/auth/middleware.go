package auth

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/rajesh_bond/production/internal/contextkey"
)

// Verifier middleware (v5)
func Verifier(tokenAuth *jwtauth.JWTAuth) func(next http.Handler) http.Handler {
	return jwtauth.Verifier(tokenAuth)
}

// Authenticator middleware (v5)
func Authenticator(tokenAuth *jwtauth.JWTAuth) func(next http.Handler) http.Handler {
	return jwtauth.Authenticator(tokenAuth) // ✅ must pass tokenAuth
}

// UserContextInjector reads claims from context and injects userID
func UserContextInjector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			log.Printf("Error: Failed to read claims from context: %v", err)
			http.Error(w, "Authentication error: Context claims missing.", http.StatusInternalServerError)
			return
		}

		// fmt.Println(claims)

		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			log.Printf("SECURITY ALERT: Token claims missing/invalid 'user_id'. Type %T, Value: %v", claims["user_id"], claims["user_id"])
			http.Error(w, "Authentication error: User ID claims missing or invalid.", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextkey.KeyUser, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext extracts userID from context
func GetUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value(contextkey.KeyUser).(string)
	if !ok {
		return ""
	}
	return userID
}

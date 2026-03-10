package service

import (
	"log"
	"os"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
	"github.com/rajesh_bond/production/internal/users"
)

var tokenAuth *jwtauth.JWTAuth

// var jwtSecret string

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found . Faling back to system enviormrnt variable.")
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	tokenAuth = jwtauth.New("HS256", []byte(secret), nil) // ✅ v5
}

func GetTokenAuth() *jwtauth.JWTAuth {
	return tokenAuth
}

func GenerateToken(payload users.TokenPayload, employeeId string) (string, error) {

	claims := map[string]interface{}{
		"tenant_id":   payload.TenantID,
		"user_id":     payload.UserID,
		"username":    payload.Username,
		"role_id":     payload.RoleID,
		"employee_id": employeeId,``
	}

	jwtauth.SetIssuedAt(claims, time.Now())

	_, tokenString, err := tokenAuth.Encode(claims)
	return tokenString, err

}

package utils

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func GenertaeToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// Generate tokenwith expiry retuns token and exporation time

func GenerateTokenWithTime(duration time.Duration) (string, time.Time) {
	token := GenertaeToken()
	ExpiresAt := time.Now().Add(duration)

	return token, ExpiresAt
}
